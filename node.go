package main

//node.go is the main software which delegates run. it currently contains a webserver
//which should be a gateway later

//kill -9 $(lsof -t -i:8888)

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"olympos.io/encoding/edn"

	"github.com/pkg/errors"
	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
)

var blocktime = 10000 * time.Millisecond
var logfile_name = "node.log"

const LOGLEVEL_OFF = 0

type Configuration struct {
	DelegateName   string
	PeerAddresses  []string
	NodePort       int
	WebPort        int
	DelgateEnabled bool
	CreateGenesis  bool
	//TODO
	Verbose bool
}

type TCPNode struct {
	NodePort      int
	Name          string
	addr          string
	server        net.Listener
	accepting     bool
	ConnectedChan chan net.Conn //channel of newly connected clients/peers
	Peers         []ntcl.Peer
	Mgr           *chain.ChainManager
	Starttime     time.Time
	Logger        *log.Logger
	Loglevel      int
	Config        Configuration
}

func (t *TCPNode) GetPeers() []ntcl.Peer {
	if &t.Peers == nil {
		return nil
	}
	return t.Peers
}

func (t *TCPNode) log(s string) {
	if t.Loglevel != LOGLEVEL_OFF {
		t.Logger.Println(s)
	}
}

// start listening on tcp and handle connection through channels
func (t *TCPNode) Run() (err error) {
	t.Starttime = time.Now()

	t.log("node listen on " + t.addr)
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

		t.log(fmt.Sprintf("new conn accepted %v", conn))
		//we put the new connection on the chan and handle there
		t.ConnectedChan <- conn

		// 	//TODO check if peers are alive see
		// 	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
		// 	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	}
	t.log("end run")
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
		t.log(fmt.Sprintf("accepted conn %v %v", strRemoteAddr, t.accepting))
		t.log(fmt.Sprintf("new peer %v ", newpeerConn))
		// log.Println("> ", t.Peers)
		// log.Println("# peers ", len(t.Peers))
		Verbose := true
		ntchan := ntcl.ConnNtchan(newpeerConn, "server", strRemoteAddr, Verbose)

		p := ntcl.Peer{Address: strRemoteAddr, NodePort: t.NodePort, NTchan: ntchan}
		t.Peers = append(t.Peers, p)

		go t.handleConnection(t.Mgr, ntchan)

		//conn.Close()

	}
}

//--- request out ---

//init an output connection
//TODO check if connected inbound already
func initOutbound(mainPeerAddress string, node_port int, verbose bool) ntcl.Ntchan {

	addr := mainPeerAddress + ":" + strconv.Itoa(node_port)
	//log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//log.Println("cant run")
		//return
	}

	//log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr, verbose)

	go ntcl.ReadLoop(ntchan)
	go ntcl.ReadProcessor(ntchan)
	go ntcl.WriteProcessor(ntchan)
	go ntcl.WriteLoop(ntchan, 300*time.Millisecond)
	return ntchan

}

func ping(peer ntcl.Peer) bool {
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, "")
	peer.NTchan.REQ_out <- req_msg
	time.Sleep(1000 * time.Millisecond)
	reply := <-peer.NTchan.REP_in
	success := reply == "REP#PONG#|"
	//log.Println("success ", success)
	return success
}

func FetchBlocksPeer(config Configuration, peer ntcl.Peer) []block.Block {

	//log.Println("FetchBlocksPeer ", peer)
	ping(peer)
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_GETBLOCKS, "")
	log.Println(req_msg)

	peer.NTchan.REQ_out <- req_msg
	time.Sleep(1000 * time.Millisecond)
	reply := <-peer.NTchan.REP_in
	//log.Println("reply ", reply)
	reply_msg := ntcl.ParseMessage(reply)
	var blocks []block.Block
	if err := json.Unmarshal(reply_msg.Data, &blocks); err != nil {
		panic(err)
	}

	return blocks

}

func FetchAllBlocks(config Configuration, t *TCPNode) {

	mainPeerAddress := config.PeerAddresses[0]
	verbose := true
	ntchan := initOutbound(mainPeerAddress, config.NodePort, verbose)
	peer := ntcl.CreatePeer(mainPeerAddress, mainPeerAddress, config.NodePort, ntchan)
	blocks := FetchBlocksPeer(config, peer)
	//log.Println("got blocks ", len(blocks))
	t.Mgr.Blocks = blocks
	t.Mgr.ApplyBlocks(blocks)
	//log.Println("set blocks ", len(t.Mgr.Blocks))
	for _, block := range t.Mgr.Blocks {
		log.Println(block)
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

func HandleBlockheight(t *TCPNode, msg ntcl.Message) string {
	bh := len(t.Mgr.Blocks)
	data := strconv.Itoa(bh)
	reply_msg := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_BLOCKHEIGHT, data)
	//("BLOCKHEIGHT ", reply_msg)
	return reply_msg
}

func HandleBalance(t *TCPNode, msg ntcl.Message) string {
	dataBytes := msg.Data
	t.log(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))

	//a := block.Account{AccountKey: string(msg.Data)}

	// var account block.Account

	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	a := string(msg.Data)
	balance := t.Mgr.Accounts[a]

	//s := strconv.Itoa(balance)
	// data, _ := json.Marshal(balance)
	data := strconv.Itoa(balance)
	reply_msg := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_BALANCE, data)
	return reply_msg
}

func HandleFaucet(t *TCPNode, msg ntcl.Message) string {
	t.log(fmt.Sprintf("HandleFaucet"))
	// dataBytes := msg.Data
	// var account block.Account
	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	//account := block.Account{AccountKey: string(msg.Data)}
	//t.log(fmt.Sprintf("faucet for ... %v", account.AccountKey))

	randNonce := 0
	amount := rand.Intn(10)

	keypair := chain.GenesisKeys()
	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
	//Genesis_Account := block.AccountFromString(addr)

	//tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}
	a := string(msg.Data)
	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: addr, Receiver: a}

	tx = crypto.SignTxAdd(tx, keypair)
	reply_string := chain.HandleTx(t.Mgr, tx)
	t.log(fmt.Sprintf("resp > %s", reply_string))

	reply := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_FAUCET, reply_string)
	return reply
}

//Standard Tx handler
//InteractiveTx also possible
//client requests tranaction <=> server response with challenge <=> client proves
func HandleTx(t *TCPNode, msg ntcl.Message) string {
	dataBytes := msg.Data

	var tx block.Tx

	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		panic(err)
	}
	t.log(fmt.Sprintf("tx %v ", tx))

	resp := chain.HandleTx(t.Mgr, tx)
	reply := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_TX, resp)
	return reply
}

//handle requests in telnet style i.e. string encoding
func RequestHandlerTel(t *TCPNode, ntchan ntcl.Ntchan) {
	for {
		msg_string := <-ntchan.REQ_in
		t.log(fmt.Sprintf("handle %s ", msg_string))
		msg := ntcl.ParseMessage(msg_string)

		var reply_msg string

		t.log(fmt.Sprintf("Handle %v", msg.Command))

		switch msg.Command {

		case ntcl.CMD_PING:
			reply_msg = HandlePing(msg)

		case ntcl.CMD_BALANCE:
			reply_msg = HandleBalance(t, msg)

		case ntcl.CMD_FAUCET:
			//send money to specified address
			reply_msg = HandleFaucet(t, msg)

		case ntcl.CMD_BLOCKHEIGHT:
			reply_msg = HandleBlockheight(t, msg)

		case ntcl.CMD_GETTXPOOL:
			t.log("get tx pool")

			//TODO
			data, _ := json.Marshal(t.Mgr.Tx_pool)
			reply_msg = ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_GETTXPOOL, string(data))

		case ntcl.CMD_GETBLOCKS:
			t.log("get tx pool")

			//TODO
			data, _ := json.Marshal(t.Mgr.Blocks)
			reply_msg = ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_GETBLOCKS, string(data))

			//Login would be challenge response protocol
			// case ntcl.CMD_LOGIN:
			// 	log.Println("> ", msg.Data)

		case ntcl.CMD_TX:
			t.log("Handle tx")
			reply_msg = HandleTx(t, msg)

		case ntcl.CMD_RANDOM_ACCOUNT:
			t.log("Handle random account")

			txJson, _ := json.Marshal(t.Mgr.RandomAccount())
			reply_msg = ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_RANDOM_ACCOUNT, string(txJson))

		//PUBSUB
		case ntcl.CMD_SUB:
			t.log(fmt.Sprintf("subscribe to topic %v", msg.Data))

			//quitpub := make(chan int)
			go ntcl.PublishTime(ntchan)
			go ntcl.PubWriterLoop(ntchan)
			//TODO reply sub ok

		case ntcl.CMD_SUBUN:
			t.log(fmt.Sprintf("unsubscribe from topic %v", msg.Data))

			go func() {
				//time.Sleep(5000 * time.Millisecond)
				close(ntchan.PUB_time_quit)
			}()

			//TODO reply unsub ok

		}

		//ntchan.Writer_queue <- reply_msg
		t.log(fmt.Sprintf("reply_msg %s", reply_msg))
		ntchan.REP_out <- reply_msg
	}
}

func (t *TCPNode) handleConnection(mgr *chain.ChainManager, ntchan ntcl.Ntchan) {
	//tr := 100 * time.Millisecond
	//defer ntchan.Conn.Close()
	t.log(fmt.Sprintf("handleConnection"))

	ntcl.NetConnectorSetup(ntchan)

	go RequestHandlerTel(t, ntchan)

	//go ntcl.WriteLoop(ntchan, 100*time.Millisecond)

}

func BlockContent(mgr *chain.ChainManager) string {
	content := ""

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

func AccountContent(mgr *chain.ChainManager) string {

	content := ""
	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(mgr.Accounts))

	for k, v := range mgr.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}
	return content
}

func Txpoolcontent(mgr *chain.ChainManager) string {
	content := ""
	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(mgr.Tx_pool))

	for i := 0; i < len(mgr.Tx_pool); i++ {
		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
		ctx := mgr.Tx_pool[i]
		content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
	}
	return content
}

//HTTP
func LoadContent(mgr *chain.ChainManager) string {
	content := ""

	// content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(peers))
	// for i := 0; i < len(peers); i++ {
	// 	content += fmt.Sprintf("peer ip address: %s<br>", peers[i].Address)
	// }

	content += Txpoolcontent(mgr)
	content += "<br>"

	content += "<a href=\"/blocks\">blocks</a><br>"
	content += "<a href=\"/accounts\">accounts</a><br>"

	//content += BlockContent(mgr)

	return content
}

func StatusContent(mgr *chain.ChainManager, t *TCPNode) []byte {

	servertime := time.Now()
	uptimedur := time.Now().Sub(t.Starttime)
	uptime := int64(uptimedur / time.Second)
	lastblocktime := t.Mgr.LastBlock().Timestamp
	timebehind := int64(servertime.Sub(lastblocktime) / time.Second)
	status := Status{Blockheight: len(mgr.Blocks), Starttime: t.Starttime, Uptime: uptime, Servertime: servertime, LastBlocktime: lastblocktime, Timebehind: timebehind}
	jData, _ := json.Marshal(status)
	return jData
}

type Status struct {
	Blockheight   int       `json:"Blockheight"`
	LastBlocktime time.Time `json:"LastBlocktime"`
	Servertime    time.Time `json:"Servertime"`
	Starttime     time.Time `json:"Starttime"`
	Timebehind    int64     `json:"Timebehind"`
	Uptime        int64     `json:"Uptime"`
}

func runWeb(t *TCPNode) {
	//webserver to access node state through browser
	// HTTP
	log.Printf("start webserver %d", t.Config.WebPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := LoadContent(t.Mgr)
		statusdata := StatusContent(t.Mgr, t)
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1>Status%s<br><div>%s </div>", statusdata, p)
	})

	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		p := BlockContent(t.Mgr)
		//nlog.Print(p)
		fmt.Fprintf(w, "<div>%s</div>", p)
	})

	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		p := AccountContent(t.Mgr)
		//nlog.Print(p)
		fmt.Fprintf(w, "<div>%s</div>", p)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {

		statusdata := StatusContent(t.Mgr, t)
		w.Header().Set("Content-Type", "application/json")
		w.Write(statusdata)

		//w.WriteHeader(http.StatusCreated)
		//json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {

		dat, _ := ioutil.ReadFile("node.log")
		fmt.Fprintf(w, "%s", dat)
	})

	//log.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))
	http.ListenAndServe(":"+strconv.Itoa(t.Config.WebPort), nil)

}

// create a new node
func NewNode() (*TCPNode, error) {
	return &TCPNode{
		//addr:          addr,
		accepting:     false,
		ConnectedChan: make(chan net.Conn),
	}, nil
}

// Close shuts down the TCP Server
func (t *TCPNode) Close() (err error) {
	return t.server.Close()
}

//TODO! fix nlog
func (t *TCPNode) setupLogfile() {
	//setup log file

	logFile, err := os.OpenFile(logfile_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		//nlog.Fatal(err)
	}

	//defer logfile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	logger := log.New(logFile, "node ", log.LstdFlags)
	logger.SetOutput(mw)

	//log.SetOutput(file)

	//nlog = logger
	//return logger
	t.Logger = logger

}

func runNode(t *TCPNode) {

	//setupLogfile()

	t.log(fmt.Sprintf("run node %d", t.Config.NodePort))

	// 	//if file exists read the chain

	// create block every blocktime sec

	if t.Config.DelgateEnabled {
		//go utils.DoEvery(, chain.MakeBlock(mgr, blockTime))

		go chain.MakeBlockLoop(t.Mgr, blocktime)
	}

	// if err != nil {
	// 	log.Println("error creating TCP server")
	// 	return
	// }

	go t.HandleConnect()

	t.Run()
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

	// go ntcl.PublishTime(ntchan)
	// go PubWriterLoop(ntchan)
	// go ntcl.WriteLoop(ntchan, 100*time.Millisecond)

	// time.Sleep(2000 * time.Millisecond)
}

//init sync or load of blocks
//if we have only genesis then load from mainpeer
//TODO check if we are mainpeer
//Set genesis will only be run at true genesis, after this we assume there is a longer chain out there

//WIP currently in testnet there is a single initiator which is the delegate expected to create first block
//TODO! replace with quering for blockheight?
func (t *TCPNode) initSyncChain(config Configuration) {
	if config.CreateGenesis {

		genBlock := chain.MakeGenesisBlock()
		t.Mgr.ApplyBlock(genBlock)
		//TODO!
		t.Mgr.AppendBlock(genBlock)

	} else {

		//TODO! apply blocks
		success := t.Mgr.ReadChain()
		t.log(fmt.Sprintf("read chain success %v", success))
		loaded_height := len(t.Mgr.Blocks)
		t.log(fmt.Sprintf("block height %d", loaded_height))

		//TODO! age of latest block compared to local time
		are_behind := loaded_height < 2
		if are_behind {
			t.Mgr.ResetBlocks()
			log.Println("blocks after reset ", len(t.Mgr.Blocks))
			FetchAllBlocks(config, t)
		}

	}
}

func runAll(config Configuration) {

	node, err := NewNode()
	node.Config = config
	node.addr = ":" + strconv.Itoa(node.Config.NodePort)
	node.setupLogfile()

	node.log(fmt.Sprintf("PeerAddresses: %v", node.Config.PeerAddresses))

	mgr := chain.CreateManager()
	node.Mgr = &mgr

	//TODO signatures of genesis
	node.Mgr.InitAccounts()

	node.initSyncChain(config)

	if err != nil {
		node.log(fmt.Sprintf("error creating TCP server"))
		return
	}

	//TODO! this will be intrement sync, not get full chain after the init sync
	if !config.CreateGenesis {
		go func() {
			for {
				log.Println("fetch blocks loop")
				FetchAllBlocks(config, node)
				time.Sleep(10000 * time.Millisecond)
			}
		}()
	}

	go runNode(node)

	go runWeb(node)

}

func getConf() Configuration {
	conffile := "conf.edn"
	f, err := os.Open(conffile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	dec := edn.NewDecoder(f)

	var c Configuration

	err = dec.Decode(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println("Config (raw go):")
	//fmt.Printf("%v\n", c.NodePort, c.WebPort, c.PeerAddresses)
	return c
}

func main() {

	conffile := "nodeconf.json"

	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		log.Println("config file does not exist. create a file named ", conffile)
		return
	}

	//config := LoadConfiguration(conffile)
	config := getConf()
	log.Println("DelegateName ", config.DelegateName)
	log.Println("CreateGenesis ", config.CreateGenesis)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go runAll(config)

	<-quit
	log.Println("Got quit signal: shutdown node ...")
	signal.Reset(os.Interrupt)

	log.Println("node exiting")

	//handle shutdown should never happen, need restart on OS level and error handling

}
