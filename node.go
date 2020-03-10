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
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
)

//simple node that runs standalone without peers

//var srv Server

const node_port = 8888

var blockTime = 10000 * time.Millisecond

var nlog *log.Logger
var logfile_name = "node.log"

type Configuration struct {
	PeerAddresses  []string
	NodePort       int
	WebPort        int
	DelgateEnabled bool
}

//TODO! rename TCP Node
type TCPNode struct {
	Name          string
	addr          string
	server        net.Listener
	accepting     bool
	ConnectedChan chan net.Conn //channel of newly connected clients/peers
	Peers         []ntcl.Peer
	Mgr           *chain.ChainManager
}

func (t *TCPNode) GetPeers() []ntcl.Peer {
	if &t.Peers == nil {
		return nil
	}
	return t.Peers
}

// start listening on tcp and handle connection through channels
func (t *TCPNode) Run() (err error) {

	log.Println("node listen on ", t.addr)
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		//return errors.Wrapf(err, "Unable to listen on port %s\n", t.addr)
	}
	//run forever and don't close
	//defer t.Close()

	for {
		t.accepting = true
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			break
		}
		if conn == nil {
			err = errors.New("could not create connection")
			break
		}

		log.Println("new conn accepted ", conn)
		//we put the new connection on the chan and handle there
		t.ConnectedChan <- conn

		// 	//TODO check if peers are alive see
		// 	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
		// 	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	}
	log.Println("end run")
	return
}

func (t *TCPNode) HandleDisconnect() {

}

//handle new connection
func (t *TCPNode) HandleConnect() {

	//TODO! hearbeart, check if peers are alive
	//TODO! handshake

	for {
		newpeerConn := <-t.ConnectedChan
		strRemoteAddr := newpeerConn.RemoteAddr().String()
		log.Println("accepted conn ", strRemoteAddr, t.accepting)
		log.Println("new peer ", newpeerConn)
		// log.Println("> ", t.Peers)
		// log.Println("# peers ", len(t.Peers))

		ntchan := ntcl.ConnNtchan(newpeerConn, "server", strRemoteAddr)

		p := ntcl.Peer{Address: strRemoteAddr, NodePort: node_port, NTchan: ntchan}
		t.Peers = append(t.Peers, p)

		go t.handleConnection(t.Mgr, ntchan)
		//go ChannelPeerNetwork(conn, peer)
		//setupPeer(strRemoteAddr, node_port, conn)

		//conn.Close()

	}
}

//--- request handlers ---

func HandleEcho(ins string) string {
	resp := "Echo:" + ins
	return resp
}

func HandlePing(msg ntcl.Message) string {
	reply_msg := ntcl.EncodeMsgString(ntcl.REP, "PONG", "")
	return reply_msg
}

func HandleBlockheight(mgr *chain.ChainManager, msg ntcl.Message) string {
	bh := len(mgr.Blocks)
	data := strconv.Itoa(bh)
	reply_msg := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_BLOCKHEIGHT, data)
	//log.Println("BLOCKHEIGHT ", reply_msg)
	return reply_msg
}

func HandleBalance(mgr *chain.ChainManager, msg ntcl.Message) string {
	dataBytes := msg.Data
	log.Println("HandleBalance data ", string(msg.Data), dataBytes)

	a := block.Account{AccountKey: string(msg.Data)}

	// var account block.Account

	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }
	// nlog.Println("get balance for account ", account)

	balance := mgr.Accounts[a]

	//s := strconv.Itoa(balance)
	// data, _ := json.Marshal(balance)
	data := strconv.Itoa(balance)
	reply_msg := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_BALANCE, data)
	return reply_msg
}

func HandleFaucet(mgr *chain.ChainManager, msg ntcl.Message) string {
	// dataBytes := msg.Data
	// var account block.Account
	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	account := block.Account{AccountKey: string(msg.Data)}
	log.Println("faucet for ... ", account.AccountKey)

	randNonce := 0
	amount := 10

	keypair := chain.GenesisKeys()
	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
	Genesis_Account := block.AccountFromString(addr)

	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}

	tx = crypto.SignTxAdd(tx, keypair)
	reply_string := chain.HandleTx(mgr, tx)
	log.Println("resp > ", reply_string)

	reply := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_FAUCET, reply_string)
	return reply
}

//Standard Tx handler
//InteractiveTx also possible
//client requests tranaction <=> server response with challenge <=> client proves
func HandleTx(mgr *chain.ChainManager, msg ntcl.Message) string {
	dataBytes := msg.Data

	var tx block.Tx

	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		panic(err)
	}
	log.Println("tx >> ", tx)

	resp := chain.HandleTx(mgr, tx)
	reply := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_TX, resp)
	//reply_msg := ntcl.EncodeMsgString(ntcl.REP, "PONG", "")
	return reply
}

//handle requests in telnet style i.e. string encoding
func RequestHandlerTel(mgr *chain.ChainManager, ntchan ntcl.Ntchan) {
	for {
		msg_string := <-ntchan.REQ_in
		log.Println("handle ", msg_string)
		msg := ntcl.ParseMessage(msg_string)

		var reply_msg string

		log.Println("Handle ", msg.Command)

		switch msg.Command {

		case ntcl.CMD_PING:
			reply_msg = HandlePing(msg)

		case ntcl.CMD_BALANCE:
			reply_msg = HandleBalance(mgr, msg)

		case ntcl.CMD_FAUCET:
			//send money to specified address
			reply_msg = HandleFaucet(mgr, msg)

		case ntcl.CMD_BLOCKHEIGHT:
			reply_msg = HandleBlockheight(mgr, msg)

			//Login would be challenge response protocol
			// case ntcl.CMD_LOGIN:
			// 	log.Println("> ", msg.Data)

		case ntcl.CMD_TX:
			log.Println("Handle tx")
			reply_msg = HandleTx(mgr, msg)

		case ntcl.CMD_RANDOM_ACCOUNT:
			log.Println("Handle random account")

			txJson, _ := json.Marshal(mgr.RandomAccount())
			reply_msg = ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_RANDOM_ACCOUNT, string(txJson))

		//PUBSUB
		case ntcl.CMD_SUB:
			log.Println("subscribe to topic ", msg.Data)

			//quitpub := make(chan int)
			go ntcl.PublishTime(ntchan)
			go ntcl.PubWriterLoop(ntchan)
			//TODO reply sub ok

		case ntcl.CMD_SUBUN:
			log.Println("unsubscribe from topic ", msg.Data)

			go func() {
				//time.Sleep(5000 * time.Millisecond)
				close(ntchan.PUB_time_quit)
			}()

			//TODO reply unsub ok

			// case ntcl.CMD_GETTXPOOL:
			// 	nlog.Println("get tx pool")

			// 	//TODO
			// 	data, _ := json.Marshal(chain.Tx_pool)
			// 	msg := ntcl.EncodeMsg(ntcl.REP, ntcl.CMD_GETTXPOOL, string(data))
			// 	rep_chan <- msg

			//var Tx_pool []block.Tx

		}

		//ntchan.Writer_queue <- reply_msg
		log.Println("reply_msg ", reply_msg)
		ntchan.REP_out <- reply_msg
	}
}

func (t *TCPNode) handleConnection(mgr *chain.ChainManager, ntchan ntcl.Ntchan) {
	//tr := 100 * time.Millisecond
	//defer ntchan.Conn.Close()
	log.Println("handleConnection")

	ntcl.NetConnectorSetup(ntchan)

	go RequestHandlerTel(mgr, ntchan)

	//go ntcl.WriteLoop(ntchan, 100*time.Millisecond)

}

//HTTP
func LoadContent(mgr *chain.ChainManager) string {
	content := ""

	// content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(peers))
	// for i := 0; i < len(peers); i++ {
	// 	content += fmt.Sprintf("peer ip address: %s<br>", peers[i].Address)
	// }

	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(mgr.Tx_pool))

	for i := 0; i < len(mgr.Tx_pool); i++ {
		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
		ctx := mgr.Tx_pool[i]
		content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
	}

	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(mgr.Accounts))

	for k, v := range mgr.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(mgr.Blocks))

	for i := 0; i < len(mgr.Blocks); i++ {
		current_block := mgr.Blocks[i]

		t := current_block.Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		//summary
		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", current_block.Height, tsf, current_block.Hash, current_block.Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(current_block.Txs))
		for j := 0; j < len(current_block.Txs); j++ {
			ctx := current_block.Txs[j]
			content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
		}
	}

	return content
}

func Runweb(mgr *chain.ChainManager, webport int) {
	//webserver to access node state through browser
	// HTTP
	log.Printf("start webserver %d", webport)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := LoadContent(mgr)
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	nlog.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))

}

// create a new node
func NewNode(addr string) (*TCPNode, error) {
	return &TCPNode{
		addr:          addr,
		accepting:     false,
		ConnectedChan: make(chan net.Conn),
		//Peers: make([]ntcl.Peer)
	}, nil

}

// Close shuts down the TCP Server
func (t *TCPNode) Close() (err error) {
	return t.server.Close()
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

	logger := log.New(logFile, "node ", log.LstdFlags)
	logger.SetOutput(mw)

	//log.SetOutput(file)

	nlog = logger
	return logger

}

func run_node(mgr *chain.ChainManager, config Configuration) {

	setupLogfile()

	nlog.Println("run node ", config.NodePort)

	// 	//if file exists read the chain

	// create block every 10sec

	if config.DelgateEnabled {
		//go utils.DoEvery(blockTime, chain.MakeBlock(mgr, blockTime))
		go chain.MakeBlockLoop(mgr, 10000*time.Millisecond)
	}

	node, err := NewNode(":" + strconv.Itoa(node_port))
	node.Mgr = mgr

	if err != nil {
		log.Println("error creating TCP server")
		return
	}

	// if err2 != nil {
	// 	log.Println("error starting TCP server ", err2)
	// 	return
	// }

	go node.HandleConnect()

	node.Run()
}

func LoadConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func pubexample() {

	// ntchan := ntcl.ConnNtchanStub("test")

	// go ntcl.PublishTime(ntchan)
	// go PubWriterLoop(ntchan)
	// go ntcl.WriteLoop(ntchan, 100*time.Millisecond)

	// time.Sleep(2000 * time.Millisecond)
}

func main() {

	mgr := chain.CreateManager()
	// 	//TODO signatures of genesis
	mgr.InitAccounts()

	//chain.ReadChain(&mgr)
	//log.Println(mgr.BlockHeight())

	// 	nlog.Println("PeerAddresses: ", config.PeerAddresses)

	success := mgr.ReadChain()
	log.Println("read chain success ", success)
	log.Printf("block height %d", len(mgr.Blocks))
	//chain.WriteGenBlock(chain.Blocks[0])

	// 	//create new genesis block (demo)
	createDemo := true //!success
	if createDemo {
		genBlock := chain.MakeGenesisBlock()
		mgr.ApplyBlock(genBlock)
		//TODO!
		mgr.AppendBlock(genBlock)
	}

	config := LoadConfiguration("nodeconf.json")

	go run_node(&mgr, config)

	Runweb(&mgr, config.WebPort)

}
