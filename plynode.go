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
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/polygonledger/edn"

	"github.com/pkg/errors"
	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/config"
	"github.com/polygonledger/node/ntcl"
)

var blocktime = 10000 * time.Millisecond
var logfile_name = "node.log"

const LOGLEVEL_OFF = 0

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
	Config        config.Configuration
}

func (t *TCPNode) GetPeers() []ntcl.Peer {
	if &t.Peers == nil {
		return nil
	}
	return t.Peers
}

func (t *TCPNode) log(s string) {
	// if t.Loglevel != LOGLEVEL_OFF {
	// 	t.Logger.Println(s)
	// }
	fmt.Println(s)
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
func (t *TCPNode) HandleConnectTCP() {

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

	//log.Println("connected")NetMsgRead
	ntchan := ntcl.ConnNtchan(conn, "client", addr, verbose)

	go ntcl.ReadLoop(ntchan)
	go ntcl.ReadProcessor(ntchan)
	go ntcl.WriteProcessor(ntchan)
	go ntcl.WriteLoop(ntchan, 300*time.Millisecond)
	return ntchan

}

func ping(peer ntcl.Peer) bool {
	req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_PING)
	peer.NTchan.REQ_out <- req_msg
	time.Sleep(1000 * time.Millisecond)
	reply := <-peer.NTchan.REP_in
	success := reply == "{:REP PONG}"
	//log.Println("success ", success)
	return success
}

func FetchBlocksPeer(config config.Configuration, peer ntcl.Peer) []block.Block {

	//log.Println("FetchBlocksPeer ", peer)
	ping(peer)
	req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_GETBLOCKS)
	//log.Println(req_msg)

	peer.NTchan.REQ_out <- req_msg
	time.Sleep(1000 * time.Millisecond)
	reply := <-peer.NTchan.REP_in
	//log.Println("reply ", reply)
	reply_msg := ntcl.ParseMessageMap(reply)
	var blocks []block.Block
	if err := json.Unmarshal(reply_msg.Data, &blocks); err != nil {
		panic(err)
	}

	return blocks

}

func FetchAllBlocks(config config.Configuration, t *TCPNode) {

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

func (t *TCPNode) handleConnection(mgr *chain.ChainManager, ntchan ntcl.Ntchan) {
	//tr := 100 * time.Millisecond
	//defer ntchan.Conn.Close()
	t.log(fmt.Sprintf("handleConnection"))

	//ntcl.NetConnectorSetup(ntchan)
	ntcl.NetConnectorSetup(ntchan)

	go RequestHandlerTel(t, ntchan)

	//go ntcl.WriteLoop(ntchan, 100*time.Millisecond)

}

type Status struct {
	Blockheight   int       `edn:"Blockheight"`
	LastBlocktime time.Time `edn:"LastBlocktime"`
	Servertime    time.Time `edn:"Servertime"`
	Starttime     time.Time `edn:"Starttime"`
	Timebehind    int64     `edn:"Timebehind"`
	Uptime        int64     `edn:"Uptime"`
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
	log.Println(fmt.Sprintf("run node on port: %d", t.Config.NodePort))

	t.log(fmt.Sprintf("run node on port: %d", t.Config.NodePort))

	// 	//if file exists read the chain

	// create block every blocktime sec

	if t.Config.DelgateEnabled {
		//go utils.DoEvery(, chain.MakeBlock(mgr, blockTime))

		//TODO!
		//go chain.MakeBlockLoop(t.Mgr, blocktime)
	}

	go t.HandleConnectTCP()

	t.Run()
}

//init sync or load of blocks
//if we have only genesis then load from mainpeer
//TODO check if we are mainpeer
//Set genesis will only be run at true genesis, after this we assume there is a longer chain out there

//WIP currently in testnet there is a single initiator which is the delegate expected to create first block
//TODO! replace with quering for blockheight?
func (t *TCPNode) initSyncChain(config config.Configuration) {
	if config.CreateGenesis {
		fmt.Println("CreateGenesis")
		genBlock := chain.MakeGenesisBlock()
		t.Mgr.ApplyBlock(genBlock)
		//TODO!
		t.Mgr.AppendBlock(genBlock)
		fmt.Println("accounts\n ", t.Mgr.Accounts)

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

func runAll(config config.Configuration) {

	log.Println("runNodeAll with config ", config)

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
	// if !config.CreateGenesis {
	// 	go func() {
	// 		for {
	// 			log.Println("fetch blocks loop")
	// 			FetchAllBlocks(config, node)
	// 			time.Sleep(10000 * time.Millisecond)
	// 		}
	// 	}()
	// }

	go runNode(node)

	log.Println("run web")
	go runWeb(node)

}

func getConf(conffile string) config.Configuration {

	f, err := os.Open(conffile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	dec := edn.NewDecoder(f)

	var c config.Configuration

	err = dec.Decode(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println("Config (raw go):")
	//fmt.Printf("%v\n", c.NodePort, c.WebPort, c.PeerAddresses)
	return c
}

func runNodeWithConfig() {

	conffile := "conf.edn"
	log.Println("config file ", conffile)

	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		log.Println("config file does not exist. create a file named ", conffile)
		return
	}

	config := getConf(conffile)
	log.Println("config ", config)
	log.Println("DelegateName ", config.DelegateName)
	log.Println("CreateGenesis ", config.CreateGenesis)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go runAll(config)

	<-quit
	// log.Println("Got quit signal: shutdown node ...")
	// signal.Reset(os.Interrupt)

	log.Println("node exiting")

	//handle shutdown should never happen, need restart on OS level and error handling

}

func main() {
	GitCommit := os.Getenv("GIT_COMMIT")
	fmt.Printf("--- run polygon ---\ngit commit: %s ----\n", GitCommit)

	runNodeWithConfig()
}
