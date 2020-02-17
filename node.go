package main

//kill -9 $(lsof -t -i:8888)
//node should run via DNS
//nodexample.com

//basic protocol
//node receives tx messages
//adds tx messages to a pool
//block gets created every 10 secs

//getBlocks
//registerPeer
//pickRandomAccount
//storeBalance

//newWallet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/net"
)

// start listening on tcp and put connection into go routine
func ListenAll() error {
	log.Println("listen all")
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		log.Println(err)
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
		go handleMessagesConn(conn)
	}
}

func Reply(rw *bufio.ReadWriter, resp string) {
	response := protocol.EncodeReply(resp)
	n, err := rw.WriteString(response)
	if err != nil {
		log.Println(err, n)
		//err:= errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	rw.Flush()
}

func ReadMessage(rw *bufio.ReadWriter) protocol.Message {
	var msg protocol.Message
	msgString := protocol.ReadStream(rw)
	if msgString == protocol.EMPTY_MSG {
		return protocol.EmptyMsg()
	}
	msg = protocol.ParseMessage(msgString)
	return msg
}

func putMsg(msg_in_chan chan string, msg string) {
	msg_in_chan <- msg
}

func handleMsg(msg_in_chan chan string, msg_out_chan chan string) {
	msgString := <-msg_in_chan
	fmt.Println("handle msg string ", msgString)

	if msgString == protocol.EMPTY_MSG {
		fmt.Println("empty msg")
		return
	}

	msg := protocol.ParseMessage(msgString)

	fmt.Println("msg type ", msg.MessageType)

	if msg.MessageType == protocol.REQ {

		log.Println("Handle ", msg.Command)

		//TODO is msg.command  valid?

		switch msg.Command {

		//TODO
		//CMD_VALIDTX

		case protocol.CMD_PING:
			log.Println("PING PONG")
			reply := "PONG"
			msg_out_chan <- reply

		case protocol.CMD_BALANCE:
			log.Println("Handle balance")

			dataBytes := msg.Data
			log.Println("data ", dataBytes)
			var account block.Account

			if err := json.Unmarshal(dataBytes, &account); err != nil {
				panic(err)
			}
			log.Println("get balance for account ", account)

			balance := chain.Accounts[account]
			s := strconv.Itoa(balance)
			msg_out_chan <- s

		case protocol.CMD_FAUCET:
			//send money to specified address

			dataBytes := msg.Data
			var account block.Account
			if err := json.Unmarshal(dataBytes, &account); err != nil {
				panic(err)
			}
			log.Println("faucet for ... ", account)

			randNonce := 0
			amount := 10

			keypair := chain.GenesisKeys()
			addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
			Genesis_Account := block.AccountFromString(addr)

			tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}
			//log.Println("tx >>> ", tx)

			tx = crypto.SignTxAdd(tx, keypair)
			reply := chain.HandleTx(tx)
			log.Println("resp > ", reply)

			msg_out_chan <- reply

		case protocol.CMD_GETTXPOOL:
			log.Println("get tx pool")

			//TODO
			data, _ := json.Marshal(chain.Tx_pool)
			msg := protocol.EncodeMsg(protocol.REP, protocol.CMD_GETTXPOOL, string(data))
			msg_out_chan <- msg

			//var Tx_pool []block.Tx

		//case CMD_GETBLOCKS:

		case protocol.CMD_TX:
			log.Println("Handle tx")

			dataBytes := msg.Data

			var tx block.Tx

			if err := json.Unmarshal(dataBytes, &tx); err != nil {
				panic(err)
			}
			log.Println(">> ", tx)

			resp := chain.HandleTx(tx)
			msg_out_chan <- resp

		// 	//log.Println("amount ", tx.Amount)
		// 	//n, err := rw.WriteString("response " + strconv.Itoa(tx.Amount) + string(protocol.DELIM))

		// case protocol.CMD_RANDOM_ACCOUNT:
		// 	log.Println("Handle random account")

		// 	txJson, _ := json.Marshal(chain.RandomAccount())
		// 	Reply(rw, string(txJson))

		// 	//log.Println("amount ", tx.Amount)
		// 	//n, err := rw.WriteString("response " + strconv.Itoa(tx.Amount) + string(protocol.DELIM))

		// 	//log.Println("amount ", tx.Amount)
		// 	//n, err := rw.WriteString("response " + strconv.Itoa(tx.Amount) + string(protocol.DELIM))

		default:
			log.Println("unknown cmd ", msg.Command)
		}
	}
}

func requestReplyLoop(rw *bufio.ReadWriter, msg_in_chan chan string, msg_out_chan chan string) {

	//continously read for requests and reply
	for {
		//REQUEST<>REPLY protocol only so far

		// read from network
		msgString := protocol.ReadStream(rw)
		log.Print("Receive message ", msgString)

		//put in the channel
		go putMsg(msg_in_chan, msgString)

		//handle in channel and put reply in msg_out channel
		go handleMsg(msg_in_chan, msg_out_chan)

		//take from channel and send over network
		reply := <-msg_out_chan
		fmt.Println("msg out ", reply)
		Reply(rw, reply)

	}
}

//handle connections
func handleMessagesConn(conn net.Conn) {

	//TODO use msg types
	msg_in_chan := make(chan string)
	msg_out_chan := make(chan string)

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//TODO
	//when close?
	//defer conn.Close()

	go requestReplyLoop(rw, msg_in_chan, msg_out_chan)

	//go publishLoop(rw, msg_in_chan, msg_out_chan)

}

func serverNode() {
	// Start listening
	ListenAll()
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
		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
		ctx := chain.Tx_pool[i]
		content += fmt.Sprintf("%x, %d from %s to %s<br>", ctx.Id, ctx.Amount, ctx.Sender, ctx.Receiver)
	}

	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(chain.Accounts))

	for k, v := range chain.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(chain.Blocks))

	for i := 0; i < len(chain.Blocks); i++ {
		t := chain.Blocks[i].Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		//summary
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

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// start node listening for incoming requests
func main() {

	log.Println("run node")

	//TODO signatures of genesis
	chain.InitAccounts()

	genBlock := chain.MakeGenesisBlock()
	chain.ApplyBlock(genBlock)
	chain.AppendBlock(genBlock)

	// create block every 10sec
	blockTime := 10000 * time.Millisecond
	go doEvery(blockTime, chain.MakeBlock)

	// //node server

	go ListenAll()

	runweb()
	//log.Println("Server running")

}
