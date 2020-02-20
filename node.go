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

type Configuration struct {
	PeerAddresses []string
	node_port     int
	web_port      int
}

//INBOUND
func addpeer(addr string) protocol.Peer {
	p := protocol.Peer{Address: addr, Req_chan: make(chan protocol.Message), Rep_chan: make(chan protocol.Message), Out_req_chan: make(chan protocol.Message), Out_rep_chan: make(chan protocol.Message)}
	Peers = append(Peers, p)
	nlog.Println("peers ", Peers)
	return p
}

func setupPeer(addr string, conn net.Conn) {
	peer := addpeer(addr)

	nlog.Println("setup channels for incoming requests")
	//TODO peers chan
	go channelNetwork(conn, peer)
}

// start listening on tcp and handle connection through channels
func ListenAll(node_port int) error {
	nlog.Println("listen all")
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", strconv.Itoa(node_port))
	if err != nil {
		nlog.Println(err)
		return errors.Wrapf(err, "Unable to listen on port %d\n", node_port) //protocol.Port
	}

	addr := listener.Addr().String()
	nlog.Println("Listen on", addr)

	//TODO check if peers are alive see
	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	for {
		nlog.Println("Accept a connection request")

		//TODO peer handshake
		//TODO client handshake

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

func HandlePing(msg_out_chan chan protocol.Message) {
	//reply := protocol.EncodeMsgString(protocol.REP, protocol.CMD_PONG, protocol.EMPTY_DATA)
	reply := protocol.EncodeMsg(protocol.REP, protocol.CMD_PONG, protocol.EMPTY_DATA)
	msg_out_chan <- reply
}

func HandleReqMsg(msg protocol.Message, rep_chan chan protocol.Message) {
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
		//s := strconv.Itoa(balance)
		data, _ := json.Marshal(balance)
		reply := protocol.EncodeMsgBytes(protocol.REP, protocol.CMD_BALANCE, data)
		log.Println(">> ", reply)

		rep_chan <- reply

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
		reply_string := chain.HandleTx(tx)
		nlog.Println("resp > ", reply_string)

		reply := protocol.EncodeMsg(protocol.REP, protocol.CMD_FAUCET, reply_string)

		rep_chan <- reply

	case protocol.CMD_BLOCKHEIGHT:

		data, _ := json.Marshal(len(chain.Blocks))
		reply := protocol.EncodeMsgBytes(protocol.REP, protocol.CMD_BLOCKHEIGHT, data)
		log.Println("CMD_BLOCKHEIGHT >> ", reply)

		rep_chan <- reply

	case protocol.CMD_TX:
		nlog.Println("Handle tx")

		dataBytes := msg.Data

		var tx block.Tx

		if err := json.Unmarshal(dataBytes, &tx); err != nil {
			panic(err)
		}
		nlog.Println(">> ", tx)

		resp := chain.HandleTx(tx)
		msg := protocol.EncodeMsg(protocol.REP, protocol.CMD_TX, resp)
		rep_chan <- msg

	case protocol.CMD_GETTXPOOL:
		nlog.Println("get tx pool")

		//TODO
		data, _ := json.Marshal(chain.Tx_pool)
		msg := protocol.EncodeMsg(protocol.REP, protocol.CMD_GETTXPOOL, string(data))
		rep_chan <- msg

		//var Tx_pool []block.Tx

	// case protocol.CMD_RANDOM_ACCOUNT:
	// 	nlog.Println("Handle random account")

	// 	txJson, _ := json.Marshal(chain.RandomAccount())

	default:
		nlog.Println("unknown cmd ", msg.Command)
		resp := "ERROR UNKONWN CMD"
		//TODO error message
		msg := protocol.EncodeMsg(protocol.REP, protocol.CMD_TX, resp)
		rep_chan <- msg
	}
}

//handle messages
func HandleMsg(req_chan chan protocol.Message, rep_chan chan protocol.Message) {
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

func ReplyLoop(rw *bufio.ReadWriter, req_chan chan protocol.Message, rep_chan chan protocol.Message) {

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

func ReqLoop(rw *bufio.ReadWriter, out_req_chan chan protocol.Message, out_rep_chan chan protocol.Message) {

	//
	for {
		request := <-out_req_chan
		log.Println("request ", request)
	}
}

//setup the network of channels
func channelNetwork(conn net.Conn, peer protocol.Peer) {

	//TODO use msg types
	//req_chan := make(chan protocol.Message)
	//rep_chan := make(chan protocol.Message)

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//TODO
	//when close?
	//defer conn.Close()

	//REQUEST<>REPLY protocol only so far

	go ReplyLoop(rw, peer.Req_chan, peer.Rep_chan)

	go ReqLoop(rw, peer.Out_req_chan, peer.Out_rep_chan)

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

	content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(Peers))
	for i := 0; i < len(Peers); i++ {
		content += fmt.Sprintf("peer ip address: %s", Peers[i].Address)
	}

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

func runweb(webport int) {
	//webserver to access node state through browser
	// HTTP
	nlog.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := loadContent()
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	nlog.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))

}

func connect_peers(node_port int, PeerAddresses []string) {

	//TODO

	for _, peer := range PeerAddresses {
		conn := protocol.OpenConn(peer + strconv.Itoa(node_port))
		log.Println(conn)

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		log.Println(rw)
		out_req := make(chan protocol.Message)
		out_rep := make(chan protocol.Message)
		ReqLoop(rw, out_req, out_rep)
		//log.Println("ping ", peer)
		//MakePing(req_chan, rep_chan)

	}
}

func run_node(node_port int) {
	nlog.Println("run node")

	//TODO signatures of genesis
	chain.InitAccounts()

	genBlock := chain.MakeGenesisBlock()
	chain.ApplyBlock(genBlock)
	chain.AppendBlock(genBlock)

	// create block every 10sec
	blockTime := 10000 * time.Millisecond
	go doEvery(blockTime, chain.MakeBlock)

	//connect_peers(configuration.PeerAddresses)

	go ListenAll(node_port)

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

	file, _ := os.Open("nodeconf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("PeerAddresses: ", configuration.PeerAddresses)

	run_node(configuration.node_port)

	nlog.Println("node running")

	runweb(configuration.web_port)

}
