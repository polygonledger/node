package main

//basic server based protocol
//server receives tx messages
//adds tx messages to a pool
//block gets created every 10 secs

//TODO
//Tx, sender receiver
//var hash = sha256("secret")
//var keypair = MakeKeypair(hash)

//Delegates
//rounds
//slotTime = getSlotNumber(currentBlockData.time))
//if slotTime generate block

//newWallet

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"encoding/gob"
	"encoding/json"
	"net/http"

	block "github.com/polygon/block"

	chain "github.com/polygon/chain"
	cryptoutil "github.com/polygon/crypto"
	protocol "github.com/polygon/net"
)

var tx_pool []block.Tx
var blocks []block.Block
var latest_block block.Block
var accounts map[chain.Account]int

//#### blockchain functions

//empty the pool
func emptyPool() {
	tx_pool = []block.Tx{}
}

func blockHash(block block.Block) block.Block {
	//FIX hash of proper data, merkle and such
	new_hash := []byte(string(block.Timestamp.Format("2020-02-02 16:06:06"))[:])
	blockheight_string := []byte(strconv.FormatInt(int64(block.Height), 10))
	//new_hash = append(new_hash, blockheight_string)
	log.Printf("newhash %x", new_hash)
	block.Hash = sha256.Sum256(blockheight_string)
	return block
}

func moveCash(SenderAccount chain.Account, ReceiverAccount chain.Account, amount int) {
	//if accounts[SenderAccount] >= amount { //sufficient balance
	accounts[SenderAccount] -= amount
	accounts[ReceiverAccount] += amount
	//}
}

func applyTx(tx block.Tx) {
	//TODO check transaction type, not implemented yet
	moveCash(tx.Sender, tx.Receiver, tx.Amount)
}

func setAccount(account chain.Account, balance int) {
	accounts[account] = balance
}

func showAccount(account chain.Account) {
	log.Printf("%s %d", account, accounts[account])
}

func makeGenesisBlock() block.Block {
	emptyhash := [32]byte{}
	timestamp := time.Now() //.Unix()
	b := []byte("banks on brink again")[:]
	genHash := sha256.Sum256(b)

	//add 10 genesis tx
	genesisTx := []block.Tx{}
	for i := 0; i < 10; i++ {
		someTx := protocol.GenesisTx()
		genesisTx = append(genesisTx, someTx)
	}

	genesis_block := block.Block{Height: 0, Txs: genesisTx, Prev_Block_Hash: emptyhash, Hash: genHash, Timestamp: timestamp}
	return genesis_block
}

func appendBlock(new_block block.Block) {
	latest_block = new_block
	blocks = append(blocks, new_block)
}

func applyBlock(block block.Block) {
	//apply
	for j := 0; j < len(block.Txs); j++ {
		applyTx(block.Txs[j])
	}

}

func makeBlock(t time.Time) {

	log.Printf("make block?")
	start := time.Now()
	//elapsed := time.Since(start)
	log.Printf("%s", start)

	//create new block if there is tx in the pool
	if len(tx_pool) > 0 {

		timestamp := time.Now() //.Unix()
		new_block := block.Block{Height: len(blocks), Txs: tx_pool, Prev_Block_Hash: latest_block.Hash, Timestamp: timestamp}
		new_block = blockHash(new_block)
		applyBlock(new_block)
		appendBlock(new_block)

		log.Printf("new block %v", new_block)
		emptyPool()

		latest_block = new_block
	} else {
		log.Printf("no block to make")
		//handle special case of no tx
		//now we don't add blocks, which means there are empty periods and blocks are not evenly spaced in time
	}

}

/*
Outgoing connections
 A `net.Conn` satisfies the io.Reader and `io.Writer` interfaces
*/

// connects to a TCP Address and returns a connection with a timeout and wrapped into a buffered ReadWriter
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

/*
Incoming connections
receive the name of a command terminated by `\n`, followed by data.
create an `Endpoint` object with the following properties:
* It allows to register one or more handler functions, where each can handle a
  particular command.
* It dispatches incoming commands to the associated handler based on the commands
  name.
*/

// handles an incoming command. It receives the open connection wrapped in a `ReadWriter` interface
type HandleFunc func(*bufio.ReadWriter)

// provides an endpoint to other processess that they can send data to
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	// Maps are not threadsafe, so we need a mutex to control access
	m sync.RWMutex
}

// create new endpoint
func NewEndpoint() *Endpoint {
	// empty list of handler funcs
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc adds a new function for handling incoming data
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// starts listening on the endpoint port on all interfaces
func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

// handleMessages reads the connection up to the first newline
// Based on this string, it calls the appropriate HandleFunc
func (e *Endpoint) handleMessages(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		// Trim the request string - ReadString does not strip any newlines.
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd)

		// Fetch the appropriate handler function from the 'handler' map and call it.
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

// handles "GOB" request
func handleTx(rw *bufio.ReadWriter) {
	//log.Print("Receive GOB data")
	var tx block.Tx
	// decodes into a struct
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&tx)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	tx_json, _ := json.Marshal(tx)
	log.Printf("receive data: %s\n", tx_json)

	//TODO its own function txhash
	//hash of timestamp is same, check lenght of bytes used??
	timestamp := time.Now().Unix()
	//b := []byte(append(string(timestamp)[:], string(tx.Nonce)[:]))

	b := []byte(string(tx.Nonce)[:])
	hash := sha256.Sum256(b)
	log.Println("hash %x time %s", hash, timestamp)

	//hash := sha256.Sum256(xdata)

	tx.Id = hash

	//TODO its own function
	tx_pool = append(tx_pool, tx)
	//blockheight += 1

	//log.Printf("tx_pool_size: \n%d", tx_pool_size)

	log.Printf("tx_pool size: %d\n", len(tx_pool))

}

// server listens for incoming requests and dispatches them to
// registered handler functions
func server() error {
	endpoint := NewEndpoint()

	// Add the handle funcs
	endpoint.AddHandleFunc(protocol.CMD_TX, handleTx)

	// Start listening.
	return endpoint.Listen()
}

//basic threading helper
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

//HTTP
func loadContent() string {
	content := ""

	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(tx_pool))

	for i := 0; i < len(tx_pool); i++ {
		content += fmt.Sprintf("Nonce %d, Id %x<br>", tx_pool[i].Nonce, tx_pool[i].Id[:])
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(blocks))

	for i := 0; i < len(blocks); i++ {
		t := blocks[i].Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", blocks[i].Height, tsf, blocks[i].Hash, blocks[i].Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(blocks[i].Txs))
		for j := 0; j < len(blocks[i].Txs); j++ {
			ctx := blocks[i].Txs[j]
			content += fmt.Sprintf("%x, %d from %s to %s<br>", ctx.Id, ctx.Amount, ctx.Sender, ctx.Receiver)
		}
	}

	return content
}

/*
start server listening for incoming requests
*/
func main() {

	accounts = make(map[chain.Account]int)
	setAccount(chain.AccountFromString("test"), 22)
	showAccount(chain.AccountFromString("test"))

	//cryptoutil.KeyExample()

	//btcec.PublicKey
	s := cryptoutil.RandomPublicKey()
	log.Printf("%s", s)

	appendBlock(makeGenesisBlock())

	//create block every 10sec
	blockTime := 10000 * time.Millisecond
	go doEvery(blockTime, makeBlock)

	//node server
	go server()
	// if err != nil {
	// 	log.Println("Error:", errors.WithStack(err))
	// }

	//webserver to access node state through browser
	// HTTP
	log.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := loadContent()
		//log.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

	log.Println("Server running")
}
