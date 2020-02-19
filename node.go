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
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/ntwk"
)

var Peers []protocol.Peer
var nlog *log.Logger
var logfile_name = "node.log"

func addpeer(addr string) {
	p := protocol.Peer{Address: addr}
	Peers = append(Peers, p)
	nlog.Println("peers ", Peers)
}

func setupPeer(addr string, conn net.Conn) {
	addpeer(addr)

	nlog.Println("setup channels for incoming requests")
	go channelNetwork(conn)
}

// start listening on tcp and handle connection through channels
func ListenAll() error {
	nlog.Println("listen all")
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		nlog.Println(err)
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}

	addr := listener.Addr().String()
	nlog.Println("Listen on", addr)

	//TODO check if peers are alive see
	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	for {
		nlog.Println("Accept a connection request")

		conn, err := listener.Accept()
		strRemoteAddr := conn.RemoteAddr().String()

		nlog.Println("accepted conn ", strRemoteAddr)
		if err != nil {
			nlog.Println("Failed accepting a connection request:", err)
			continue
		}

		setupPeer(strRemoteAddr, conn)

	}
}

func putMsg(msg_in_chan chan protocol.Message, msg protocol.Message) {
	msg_in_chan <- msg
}

func HandlePing(msg_out_chan chan string) {
	reply := protocol.EncodeMsgString(protocol.REP, protocol.CMD_PONG, protocol.EMPTY_DATA)
	msg_out_chan <- reply
}

func HandleReqMsg(msg protocol.Message, rep_chan chan string) {
	nlog.Println("Handle ", msg.Command)

	switch msg.Command {

	case protocol.CMD_PING:
		nlog.Println("PING PONG")
		HandlePing(rep_chan)

	case protocol.CMD_BALANCE:
		nlog.Println("Handle balance")

		dataBytes := msg.Data
		nlog.Println("data ", dataBytes)
		var account block.Account

		if err := json.Unmarshal(dataBytes, &account); err != nil {
			panic(err)
		}
		nlog.Println("get balance for account ", account)

		balance := chain.Accounts[account]
		s := strconv.Itoa(balance)
		rep_chan <- s

	case protocol.CMD_FAUCET:
		//send money to specified address

		dataBytes := msg.Data
		var account block.Account
		if err := json.Unmarshal(dataBytes, &account); err != nil {
			panic(err)
		}
		nlog.Println("faucet for ... ", account)

		randNonce := 0
		amount := 10

		keypair := chain.GenesisKeys()
		addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
		Genesis_Account := block.AccountFromString(addr)

		tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}

		tx = crypto.SignTxAdd(tx, keypair)
		reply := chain.HandleTx(tx)
		nlog.Println("resp > ", reply)

		rep_chan <- reply

	case protocol.CMD_GETTXPOOL:
		nlog.Println("get tx pool")

		//TODO
		data, _ := json.Marshal(chain.Tx_pool)
		msg := protocol.EncodeMsgString(protocol.REP, protocol.CMD_GETTXPOOL, string(data))
		rep_chan <- msg

		//var Tx_pool []block.Tx

	//case CMD_GETBLOCKS:

	case protocol.CMD_TX:
		nlog.Println("Handle tx")

		dataBytes := msg.Data

		var tx block.Tx

		if err := json.Unmarshal(dataBytes, &tx); err != nil {
			panic(err)
		}
		nlog.Println(">> ", tx)

		resp := chain.HandleTx(tx)
		rep_chan <- resp

	// case protocol.CMD_RANDOM_ACCOUNT:
	// 	nlog.Println("Handle random account")

	// 	txJson, _ := json.Marshal(chain.RandomAccount())

	default:
		nlog.Println("unknown cmd ", msg.Command)
		resp := "ERROR UNKONWN CMD"
		rep_chan <- resp
	}
}

//handle messages
func HandleMsg(req_chan chan protocol.Message, rep_chan chan string) {
	req_msg := <-req_chan
	//msgString := <-req_chan
	//msgString :=
	//	fmt.Println("handle msg string ", msgString)

	fmt.Println("msg type ", req_msg.MessageType)

	if req_msg.MessageType == protocol.REQ {
		HandleReqMsg(req_msg, rep_chan)
	} else if req_msg.MessageType == protocol.REP {
		nlog.Println("handle reply")
	}
}

func RequestReplyLoop(rw *bufio.ReadWriter, req_chan chan protocol.Message, rep_chan chan string) {

	//continously read for requests and respond with reply
	for {

		// read from network
		msgString := protocol.NetworkReadMessage(rw)
		if msgString == protocol.EMPTY_MSG {
			//log.Println("empty message, ignore")
			time.Sleep(500 * time.Millisecond)
			continue
		}

		msg := protocol.ParseMessage(msgString)
		nlog.Print("Receive message over network ", msgString)

		//put in the channel
		go putMsg(req_chan, msg)

		//handle in channel and put reply in msg_out channel
		go HandleMsg(req_chan, rep_chan)

		//take from reply channel and send over network
		reply := <-rep_chan
		fmt.Println("msg out ", reply)
		protocol.ReplyNetwork(rw, reply)

	}
}

//setup the network of channels
func channelNetwork(conn net.Conn) {

	//TODO use msg types
	req_chan := make(chan protocol.Message)
	rep_chan := make(chan string)

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//TODO
	//when close?
	//defer conn.Close()

	//REQUEST<>REPLY protocol only so far

	go RequestReplyLoop(rw, req_chan, rep_chan)

	//go publishLoop(msg_in_chan, msg_out_chan)

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
		content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
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
			content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
		}
	}

	return content
}

func runweb() {
	//webserver to access node state through browser
	// HTTP
	nlog.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := loadContent()
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	nlog.Fatal(http.ListenAndServe(":8080", nil))

}

func run_node() {
	nlog.Println("run node")

	//TODO signatures of genesis
	chain.InitAccounts()

	genBlock := chain.MakeGenesisBlock()
	chain.ApplyBlock(genBlock)
	chain.AppendBlock(genBlock)

	// create block every 10sec
	blockTime := 10000 * time.Millisecond
	go doEvery(blockTime, chain.MakeBlock)

	go ListenAll()

}

func setupLogfile() *log.Logger {
	//setup log file

	logFile, err := os.OpenFile(logfile_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		nlog.Fatal(err)
	}

	//defer logfile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//logger := log.New(logFile, "", log.LstdFlags)
	logger := log.New(logFile, "node ", log.LstdFlags)
	logger.SetOutput(mw)

	//log.SetOutput(file)

	nlog = logger
	return logger

}

// start node listening for incoming requests
func main() {

	setupLogfile()

	run_node()

	nlog.Println("node running")

	runweb()

}
