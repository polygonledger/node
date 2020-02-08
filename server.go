package main

//basic server based protocol
//server receives tx messages
//adds tx messages to a pool
//block gets created every 10 secs

//TODO
//package polygonledger/node
//basic signatures

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
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"encoding/gob"
	"encoding/json"

	block "github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	protocol "github.com/polygonledger/node/net"
)

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

// starts listening
func ListenAll() error {
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}
	log.Println("Listen on", listener.Addr().String())
	for {
		log.Println("Accept a connection request")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages")
		go handleMessagesChan(conn)
	}
}

func handleMessagesChan(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	defer conn.Close()
	for {
		// read
		log.Print("Receive message")

		msg, err := rw.ReadString(protocol.DELIM)
		log.Print("msg ", msg)
		if err != nil {
			log.Println("Failed ", err)
			//log.Println(err.)
			continue
		}

		msg = strings.Trim(msg, string(protocol.DELIM))
		//msg = strings.Trim(msg, string("\n"))
		s := strings.Split(msg, "#")
		mtype := s[0]
		cmd := s[1]
		fmt.Println(mtype, cmd)

		fmt.Println(protocol.IsValidMsgType(mtype))
	}
}

// starts listening on the endpoint port
func (e *Endpoint) ListenCmd() error {
	var err error
	e.listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages")
		go e.handleMessagesProt(conn)
	}
}

// handleMessages reads the connection up to the first newline
// based on the command, it calls the appropriate HandleFunc
func (e *Endpoint) handleMessagesProt(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	for {
		// read
		log.Print("Receive message")

		msg, err := rw.ReadString(protocol.DELIM)
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - no delim?")
			log.Println(msg)
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+msg+"'\n", err)
			return
		}
		// Trim the request string - ReadString does not strip any newlines.
		log.Println("msg: ", msg)
		msg = strings.Trim(msg, string(protocol.DELIM))
		log.Println("msg: ", msg)

		var msgStruct protocol.Message
		dec := gob.NewDecoder(rw)
		errDecode := dec.Decode(&msgStruct)

		if errDecode != nil {
			log.Println("Error decoding GOB data:", errDecode)
			return
		} else {
			log.Println("decoded ", msgStruct)
		}

		//cmd = strings.Trim(cmd, "\n")

	}
}

// handleMessages reads the connection up to the first newline
// based on the command, it calls the appropriate HandleFunc
func (e *Endpoint) handleMessagesCmd(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	for {
		// read
		log.Print("Receive command")
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
		log.Println("cmd: ", cmd)

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

// handle ranaccount request
func handleRandomAccountRequest(rw *bufio.ReadWriter) {
	protocol.SendAccount(rw)
}

// handle tx request
func handleTxRequest(rw *bufio.ReadWriter) {
	//log.Print("Receive GOB data")

	//GOB basics
	// decodes into a struct
	var tx block.Tx
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&tx)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	//json example, not used
	tx_json, _ := json.Marshal(tx)
	log.Printf("receive data: %s\n", tx_json)
	chain.HandleTx(tx)

}

// server listens for incoming requests and dispatches them to
// registered handler functions
func serverCmd() error {
	endpoint := NewEndpoint()

	// Add the handle funcs
	endpoint.AddHandleFunc(protocol.CMD_TX, handleTxRequest)
	endpoint.AddHandleFunc(protocol.CMD_RANDOM_ACCOUNT, handleRandomAccountRequest)

	// Start listening
	return endpoint.ListenCmd()
}

func serverNode() error {
	// Start listening
	return ListenAll()
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

	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(chain.Tx_pool))

	for i := 0; i < len(chain.Tx_pool); i++ {
		content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
	}

	content += fmt.Sprintf("<h2>Accounts</h2>%d<br>", len(chain.Accounts))

	for k, v := range chain.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(chain.Blocks))

	for i := 0; i < len(chain.Blocks); i++ {
		t := chain.Blocks[i].Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", chain.Blocks[i].Height, tsf, chain.Blocks[i].Hash, chain.Blocks[i].Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(chain.Blocks[i].Txs))
		for j := 0; j < len(chain.Blocks[i].Txs); j++ {
			ctx := chain.Blocks[i].Txs[j]
			content += fmt.Sprintf("%x, %d from %s to %s<br>", ctx.Id, ctx.Amount, ctx.Sender, ctx.Receiver)
		}
	}

	return content
}

func runweb() {
	//webserver to access node state through browser
	// HTTP
	log.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := loadContent()
		//log.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

/*
start server listening for incoming requests
*/
func main() {

	serverNode()

	//demo
	// kp := cryptoutil.SomeKeypair()
	// fmt.Println("some key ", kp)
	// cryptoutil.SignExample(kp)

	//cryptoutil.KeyExample()

	//btcec.PublicKey
	// s := cryptoutil.RandomPublicKey()
	// log.Printf("%s", s)

	//account := block.Account{AccountKey: "test"}
	//accountJson, _ := json.Marshal(account)
	//fmt.Println(string(accountJson))

	/////

	// chain.InitAccounts()
	//demo
	//chain.SetAccount(block.AccountFromString("test"), 22)
	//chain.ShowAccount(block.AccountFromString("test"))

	// rk := cryptoutil.RandomPublicKey()
	// ra := cryptoutil.Address(rk)
	// log.Printf("random address %s", ra)

	// genBlock := chain.MakeGenesisBlock()
	// chain.ApplyBlock(genBlock)
	// chain.AppendBlock(genBlock)

	// //create block every 10sec
	// blockTime := 10000 * time.Millisecond
	// go doEvery(blockTime, chain.MakeBlock)

	// //node server
	// go serverCmd()
	// // if err != nil {
	// // 	log.Println("Error:", errors.WithStack(err))
	// // }

	// log.Println("Server running")

	//go runweb()

}
