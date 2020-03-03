package main

//needs refactor

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
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/ntwk"
)

var Peers []ntwk.Peer

//banned IPs
// var nlog *log.Logger
// var logfile_name = "node.log"

var blockTime = 10000 * time.Millisecond

var tchan chan string

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
}

//inbound
func addpeer(addr string, nodeport int) ntwk.Peer {
	//ignored

	p := ntwk.CreatePeer(addr, nodeport)
	Peers = append(Peers, p)
	nlog.Println("peers ", Peers)
	return p
}

//setup the network of channels
//the main junction for managing message flow between types of messages
func ChannelPeerNetwork(conn net.Conn, peer ntwk.Peer) {

	log.Println("init channelPeerNetwork")

	ntchan := ntwk.ConnNtchan(conn, peer.Address)

	//main reader and writer setup
	go ntwk.ReaderWriterConnector(ntchan)

	//TODO need to formalize this
	go Reqprocessor(ntchan)

	go Reqoutprocessor(ntchan)

	//publishers
	//go pubhearbeat(ntchan)

	//TODO! handle disconnects
	//when hearbeat fails top all workers related to the peer

	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	//when close?
	//defer conn.Close()

	//go ReplyLoop(ntchan, peer.Req_chan, peer.Rep_chan)

	//----------------
	//TODO pubsub
	//go publishLoop(msg_in_chan, msg_out_chan)

	// go ntwk.Subtime(tchan, "peer1")

	// go ntwk.Subout(tchan, "peer1", peer.Pub_chan)

}

//simplified network for testing
func SimpleNetwork(conn net.Conn, peer ntwk.Peer) {

	log.Println("init channelPeerNetwork")

	ntchan := ntwk.ConnNtchan(conn, peer.Address)

	read_loop_time := 900 * time.Millisecond
	go func() {
		for {
			//read from network and put in channel
			msg := ntwk.NetworkReadMessage(ntchan)
			log.Println("ntwk read => " + msg)
			//ntchan.Reader_queue <- msg
			time.Sleep(read_loop_time)
			//fix: need ntchan to be a pointer
			//msg_reader_total++
		}
	}()

}

//inbound
func setupPeer(addr string, nodeport int, conn net.Conn) {
	peer := addpeer(addr, nodeport)

	nlog.Println("setup channels for incoming requests")
	//TODO peers chan
	//TODO handshake
	go ChannelPeerNetwork(conn, peer)
	//go SimpleNetwork(conn, peer)

}

// start listening on tcp and handle connection through channels
func ListenAll(node_port int) error {
	nlog.Println("listen all")
	var err error
	var listener net.Listener
	p := ":" + strconv.Itoa(node_port)
	nlog.Println("listen on ", p)
	listener, err = net.Listen("tcp", p)
	if err != nil {
		nlog.Println(err)
		return errors.Wrapf(err, "Unable to listen on port %d\n", node_port) //ntwk.Port
	}

	addr := listener.Addr().String()
	nlog.Println("Listen on", addr)

	//TODO check if peers are alive see
	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	for {
		//nlog.Println("Accept a connection request ")

		conn, err := listener.Accept()
		strRemoteAddr := conn.RemoteAddr().String()

		nlog.Println("accepted conn ", strRemoteAddr)
		if err != nil {
			nlog.Println("Failed accepting a connection request:", err)
			continue
		}

		setupPeer(strRemoteAddr, node_port, conn)

	}
}

func putMsg(msg_in_chan chan ntwk.Message, msg ntwk.Message) {
	msg_in_chan <- msg
}

func pubhearbeat(ntchan ntwk.Ntchan) {
	heartbeat_time := 1000 * time.Millisecond

	for {
		msg := ntwk.EncodeHeartbeat("peer1")
		ntchan.Writer_queue <- msg
		time.Sleep(heartbeat_time)
	}

}

//--- request handler ---

func HandlePing() ntwk.Message {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)
	return reply
	//msg_out_chan <- reply
}

func HandleHandshake(msg_out_chan chan ntwk.Message) {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_HANDSHAKE_STABLE, ntwk.EMPTY_DATA)
	msg_out_chan <- reply
}

func HandleReqMsg(msg ntwk.Message) ntwk.Message {
	log.Println("Handle ", msg.Command)

	//can use handler instead i.e. map[string] => func
	switch msg.Command {

	case ntwk.CMD_PING:
		log.Println("PING PONG")
		return HandlePing()

		// case ntwk.CMD_HANDSHAKE_HELLO:
		// 	nlog.Println("handshake")
		// 	HandleHandshake(rep_chan)

		// case ntwk.CMD_BALANCE:
		// 	nlog.Println("Handle balance")

		// 	dataBytes := msg.Data
		// 	nlog.Println("data ", dataBytes)
		// 	var account block.Account

		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("get balance for account ", account)

		// 	balance := chain.Accounts[account]
		// 	//s := strconv.Itoa(balance)
		// 	data, _ := json.Marshal(balance)
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BALANCE, data)
		// 	log.Println(">> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_FAUCET:
		// 	//send money to specified address

		// 	dataBytes := msg.Data
		// 	var account block.Account
		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("faucet for ... ", account)

		// 	randNonce := 0
		// 	amount := 10

		// 	keypair := chain.GenesisKeys()
		// 	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
		// 	Genesis_Account := block.AccountFromString(addr)

		// 	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}

		// 	tx = crypto.SignTxAdd(tx, keypair)
		// 	reply_string := chain.HandleTx(tx)
		// 	nlog.Println("resp > ", reply_string)

		// 	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_FAUCET, reply_string)

		// 	rep_chan <- reply

		// case ntwk.CMD_BLOCKHEIGHT:

		// 	data, _ := json.Marshal(len(chain.Blocks))
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BLOCKHEIGHT, data)
		// 	log.Println("CMD_BLOCKHEIGHT >> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_TX:
		// 	nlog.Println("Handle tx")

		// 	dataBytes := msg.Data

		// 	var tx block.Tx

		// 	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println(">> ", tx)

		// 	resp := chain.HandleTx(tx)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		// 	rep_chan <- msg

		// case ntwk.CMD_GETTXPOOL:
		// 	nlog.Println("get tx pool")

		// 	//TODO
		// 	data, _ := json.Marshal(chain.Tx_pool)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_GETTXPOOL, string(data))
		// 	rep_chan <- msg

		//var Tx_pool []block.Tx

		// case ntwk.CMD_RANDOM_ACCOUNT:
		// 	nlog.Println("Handle random account")

		// 	txJson, _ := json.Marshal(chain.RandomAccount())

	default:
		// 	nlog.Println("unknown cmd ", msg.Command)
		resp := "ERROR UNKONWN CMD"

		msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		return msg

	}
}

func HandleReqMsgString(msg_string string) string {
	msg := ntwk.ParseMessage(msg_string)
	reply := HandleReqMsg(msg)
	reply_string := ntwk.MsgString(reply)
	return reply_string
}

//process requests
func Reqprocessor(ntchan ntwk.Ntchan) {
	for {
		log.Println("init Reqprocessor")
		msg_string := <-ntchan.REQ_in

		//reply_string := "reply"
		//reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)
		reply_string := HandleReqMsgString(msg_string)

		log.Println("forward request to writer")
		//ntchan.Writer_queue <- reply_string
		ntchan.REP_out <- reply_string
	}
}

func Reqoutprocessor(ntchan ntwk.Ntchan) {
	log.Println("init Reqoutprocessor")
	for {
		log.Println("... Reqoutprocessor")
		msg_string := <-ntchan.REQ_out

		//reply_string := "reply"
		//reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)

		log.Println("forward request to writer")

		ntchan.Writer_queue <- msg_string
	}
}

// func Repoutprocessor(ntchan ntwk.Ntchan) {
// 	for {
// 		log.Println("handler ")
// 		msg := <-ntchan.REP_out
// 		log.Println(">>> REP_out ", msg)
// 	}
// }

//channel network optimised for client
// func channelPeerNetworkClient(conn net.Conn, peer ntwk.Peer) {

// 	ntchan := ntwk.ConnNtchan(conn, peer.Address)

// 	ntwk.ReaderWriterConnector(ntchan)

// }

func connect_peers(node_port int, PeerAddresses []string) {

	for _, peer := range PeerAddresses {
		log.Println(peer)
		//TODO old!

		//addr := peer + strconv.Itoa(node_port)

		//ntchan := ntwk.OpenNtchan(addr)

		//log.Println("ping ", peer)
		//MakePingOld(req_chan, rep_chan)

	}
}

// 	nlog.Printf("block height %d", len(chain.Blocks))

//HTTP
// func LoadContent(peers []ntwk.Peer) string {
// 	content := ""

// 	content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(peers))
// 	for i := 0; i < len(peers); i++ {
// 		content += fmt.Sprintf("peer ip address: %s<br>", peers[i].Address)
// 	}

// 	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(chain.Tx_pool))

// 	for i := 0; i < len(chain.Tx_pool); i++ {
// 		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
// 		ctx := chain.Tx_pool[i]
// 		content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
// 	}

// 	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(chain.Accounts))

// 	for k, v := range chain.Accounts {
// 		content += fmt.Sprintf("%s %d<br>", k, v)
// 	}

// 	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(chain.Blocks))

// 	for i := 0; i < len(chain.Blocks); i++ {
// 		t := chain.Blocks[i].Timestamp
// 		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
// 			t.Year(), t.Month(), t.Day(),
// 			t.Hour(), t.Minute(), t.Second())

// 		//summary
// 		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", chain.Blocks[i].Height, tsf, chain.Blocks[i].Hash, chain.Blocks[i].Prev_Block_Hash)

// 		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(chain.Blocks[i].Txs))
// 		for j := 0; j < len(chain.Blocks[i].Txs); j++ {
// 			ctx := chain.Blocks[i].Txs[j]
// 			content += fmt.Sprintf("%d from %s to %s %x<br>", ctx.Amount, ctx.Sender, ctx.Receiver, ctx.Id)
// 		}
// 	}

// 	return content
// }

// func Runweb(webport int) {
// 	//webserver to access node state through browser
// 	// HTTP
// 	nlog.Println("start webserver")

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		p := LoadContent(Peers)
// 		//nlog.Print(p)
// 		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
// 	})

// 	nlog.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))

// }

//TODO
func rungin(webport int) {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	//router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Setup route group for the API
	example := "test"
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, example)
		})

		api.GET("/peers", func(c *gin.Context) {
			c.JSON(http.StatusOK, Peers)
		})

		api.GET("/txpool", func(c *gin.Context) {
			c.JSON(http.StatusOK, chain.Tx_pool)
		})

		api.GET("/blockheight", func(c *gin.Context) {
			c.JSON(http.StatusOK, len(chain.Blocks))
		})

		//json: unsupported type: map[block.Account]int
		api.GET("/accounts", func(c *gin.Context) {
			c.JSON(http.StatusOK, chain.Accounts)
		})

		//blocks

	}

	// Start and run the server
	router.Run(":" + strconv.Itoa(webport))
}

// func main() {

// 	setupLogfile()

// 	config := LoadConfiguration("nodeconf.json")

// 	nlog.Println("run node with config ", config)

// 	run_node(config)

// 	log.Println("run web on ", config.WebPort)

// 	//rungin(config.WebPort)

// 	Runweb(config.WebPort)

// }
