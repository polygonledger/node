package main

//client based application to interact with node

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
)

var Peers []ntcl.Peer

const node_port = 8888

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
}

func addPeerOut(p ntcl.Peer) {
	Peers = append(Peers, p)
	log.Println("peers now", Peers)
}

func initClient(config Configuration) ntcl.Ntchan {
	mainPeerAddress := config.PeerAddresses[0]
	addr := mainPeerAddress + ":" + strconv.Itoa(node_port)
	log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("cant run")
		//return
	}

	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr)

	go ntcl.ReadLoop(ntchan)
	go ntcl.ReadProcessor(ntchan)
	go ntcl.WriteProcessor(ntchan)
	go ntcl.WriteLoop(ntchan, 300*time.Millisecond)
	return ntchan

}

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End measure time:", s, "took", endTime.Sub(startTime))
}

//request account address
// func RequestAccount(rw *bufio.ReadWriter) error {
// 	msg := ConstructMessage(CMD_RANDOM_ACCOUNT)

// func ReceiveAccount(rw *bufio.ReadWriter) error {
// 	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)

func PushTx(peer ntcl.Peer) error {

	dat, _ := ioutil.ReadFile("tx.json")
	var tx block.Tx

	if err := json.Unmarshal(dat, &tx); err != nil {
		panic(err)
	}

	//send Tx
	txJson, _ := json.Marshal(tx)
	log.Println("txJson ", string(txJson))

	req_msg := ntcl.EncodeMessageTx(txJson)
	// response := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Print(" ", req_msg)

	return nil
}

func Getbalance(peer ntcl.Peer) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	txJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_BALANCE, string(txJson))
	log.Println(req_msg)
	// response := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	// log.Println("response ", response)
	// var balance int
	// if err := json.Unmarshal(response.Data, &balance); err != nil {
	// 	panic(err)
	// }
	// log.Println("balance of account ", balance)

	return nil
}

func Getblockheight(peer ntcl.Peer) error {
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_BLOCKHEIGHT, "")
	log.Println(req_msg)
	// response := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	// var blockheight int
	// if err := json.Unmarshal(response.Data, &blockheight); err != nil {
	// 	panic(err)
	// }
	// log.Println("blockheight ", blockheight)

	return nil
}

func Gettxpool(peer ntcl.Peer) error {
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_GETTXPOOL, "")
	log.Println("> ", req_msg)
	// resp := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	// log.Println("rcvmsg ", resp)
	// log.Println("data ", resp.Data)

	// var txp []block.Tx
	// if err := json.Unmarshal(resp.Data, &txp); err != nil {
	// 	panic(err)
	// }
	// log.Println("txp ", txp)

	return nil
}

func GetFaucet(peer ntcl.Peer) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	accountJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_FAUCET, string(accountJson))
	log.Println(req_msg)
	// resp := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	// log.Println("resp ", resp)

	return nil
}

func readdns() {
	// domain := "example.com"
	// ips, err1 := net.LookupIP(domain)
	// if err1 != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err1)
	// 	os.Exit(1)
	// }
	// for _, ip := range ips {
	// 	fmt.Printf(domain+". IN A %s\n", ip.String())
	// }

}

func setupAllPeers(config Configuration) {

	for _, peerAddress := range config.PeerAddresses {
		log.Println("setup  peer ", peerAddress)
		//p := ntcl.CreatePeer(peerAddress, config.NodePort)

		//err := setupPeerClient(p)
		// if err != nil {
		// 	log.Println("connect failed")
		// 	continue
		// }
	}

}

//run client against multiple nodes
func runPeermode(option string, config Configuration) {
	log.Println("runPeermode")

	log.Println("setup peers")

	for _, peerAddress := range config.PeerAddresses {

		p := ntcl.CreatePeer(peerAddress, config.NodePort)
		log.Println("add peer ", p)

		//err := setupPeerClient(p)
		// if err != nil {
		// 	//remove peer
		// 	log.Println("dont add peer to list")
		// } else {
		// 	addPeerOut(p)
		// }
	}

	switch option {

	case "pingall":

		defer track(runningtime("execute ping"))
		successCount := 0
		for _, peerAddress := range config.PeerAddresses {
			log.Println("setup  peer ", peerAddress, config.NodePort)
			//p := ntcl.CreatePeer(peerAddress, config.NodePort)

			//err := setupPeerClient(p)
			// if err != nil {
			// 	log.Println("connect failed")
			// 	continue
			// } else {
			// 	// success := ntcl.MakePingOld(p)
			// 	// if success {
			// 	// 	successCount++
			// 	// }
			// }
		}

		log.Println("pinged peers ", len(config.PeerAddresses), " successCount:", successCount)

	case "blockheight":

		for _, peerAddress := range config.PeerAddresses {
			log.Println("setup  peer ", peerAddress)
			//p := ntcl.CreatePeer(peerAddress, config.NodePort)

			// err := setupPeerClient(p)
			// if err == nil {
			// 	log.Println("block height ", p)
			// 	Getblockheight(p)
			// }
		}

	}

}

func requestreply(ntchan ntcl.Ntchan, req_msg string) {

	//TODO! use readloop and REQ/REP chans

	log.Println("requestreply >> ", req_msg)
	//REQUEST
	ntchan.REQ_out <- req_msg
	//REPLY

	resp_string := <-ntchan.REP_in
	log.Println("REP_in >> ", resp_string)

	// msg := ntcl.ParseMessage(resp_string)
	// log.Println("response ", msg.MessageType)
	// if msg.MessageType == ntcl.REP {
	// 	//need to match to know this is the same request ID?
	// 	log.Println("REPLY ", msg)
	// }
}

//TODO! move
func ping(ntchan ntcl.Ntchan) {

	// req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, "")

	// requestreply(ntchan, req_msg)

	//subscribe example
	//reqs := "REQ#PING#|"
	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, "")
	ntchan.REQ_out <- req_msg
	//ntcl.NetWrite(ntchan, reqs)

	time.Sleep(1000 * time.Millisecond)

	reply := <-ntchan.REP_in
	success := reply == "REP#PONG#|"
	log.Println("success ", success)

}

// func ReplyInProcessor(ntchan ntcl.Ntchan) {
// 	for {
// 		log.Println("ReplyInProcessor ")
// 		msg := <-ntchan.REP_in
// 		log.Println("ReplyInProcessor >> ", msg)
// 	}
// }

//run client against single node, just use first IP address in peers i.e. mainpeer
func runSingleMode(option string, config Configuration) {

	//mainPeerAddress := config.PeerAddresses[0]
	//log.Println("setup main peer ", mainPeerAddress, config.NodePort)
	//mainPeer := ntcl.CreatePeer(mainPeerAddress, config.NodePort)
	//log.Println("client with mainPeer ", mainPeer)
	//setupPeerClient(mainPeer)
	//conn := ntcl.OpenConn(mainPeerAddress + ":" + strconv.Itoa(config.NodePort))
	//ntcl.ChannelPeerNetwork(conn, mainPeer)
	//ntchan := ntcl.ConnNtchan(conn, mainPeerAddress)
	ntchan := initClient(config)
	log.Println("init ", ntchan)

	switch option {

	case "ping":
		log.Println("ping")
		//ntchan := initClient(config)
		ping(ntchan)

		time.Sleep(100 * time.Millisecond)

		// case "heartbeat":
		// 	log.Println("heartbeat")

		// 	for {
		// 		go ping(ntchan)
		// 		time.Sleep(1 * time.Second)
		// 	}

		// 	// success := ntcl.MakeHandshake(mainPeer)
		// 	// if success {
		// 	// 	log.Println("start heartbeat")
		// 	// 	ntcl.Hearbeat(mainPeer)
		// 	// }

		// case "getbalance":
		// 	log.Println("getbalance")

		// 	Getbalance(mainPeer)

		// case "blockheight":
		// 	log.Println("blockheight")

		// 	Getblockheight(mainPeer)

		// case "faucet":
		// 	log.Println("faucet")
		// 	//get coins
		// 	//GetFaucet(rw)
		// 	GetFaucet(mainPeer)

		// case "txpool":
		// 	_ = Gettxpool(mainPeer)
		// 	return

		// case "pushtx":
		// 	_ = PushTx(mainPeer)
		// 	return

		// case "randomtx":
		// 	_ = MakeRandomTx(mainPeer)
		// 	return

	}
}

//client that only listens to events
func runListenMode(option string, config Configuration) {
	log.Println("listen")

	//ntchan := initClient()

	// mainPeerAddress := config.PeerAddresses[0]
	// log.Println("setup main peer ", mainPeerAddress)
	// mainPeer := ntcl.CreatePeer(mainPeerAddress, config.NodePort)
	// success := ntcl.MakeHandshake(mainPeer)
	// log.Println(success)
	// log.Println("start heartbeat")
	// if success {
	// 	hTime := 2000 * time.Millisecond

	// 	for _ = range time.Tick(hTime) {
	// 		//log.Println(x)
	// 		ntcl.Hearbeat(mainPeer)
	// 	}

	// }
}

//run client without client or server
func runOffline(option string, config Configuration) {

	switch option {
	case "createkeys":
		CreateKeys()

	case "readkeys":
		kp := ReadKeys("keys.txt")
		log.Println(kp)

	case "sign":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to sign: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		fmt.Println(msg)
		kp := ReadKeys("keys.txt")
		signature := crypto.SignMsgHash(kp, msg)
		log.Println("signature ", signature)

		sighex := hex.EncodeToString(signature.Serialize())
		log.Println("sighex ", sighex)

	case "createtx":
		Createtx()

	case "verify":
		reader := bufio.NewReader(os.Stdin)
		log.Print("Enter message to verify: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		log.Println(msg)

		fmt.Print("Enter signature to verify: ")
		msgsig, _ := reader.ReadString('\n')
		msgsig = strings.Trim(msgsig, string('\n'))

		sign := crypto.SignatureFromHex(msgsig)

		fmt.Print("Enter pubkey to verify: ")
		msgpub, _ := reader.ReadString('\n')
		log.Println(msgpub)
		msgpub = strings.Trim(msgpub, string('\n'))

		pubkey := crypto.PubKeyFromHex(msgpub)

		verified := crypto.VerifyMessageSignPub(sign, pubkey, msg)
		log.Println("verified ", verified)

	}
}

//dns functions for later, as we can use txt records to get pubkey
func dnslook() {
	//domain := "test.polygonnode.com"
	domain := "swix.io"

	txtrecords, _ := net.LookupTXT(domain)
	// log.Println(txtrecords)

	// for _, txt := range txtrecords {
	// 	fmt.Println(txt)
	// }

	frec := txtrecords[0]
	log.Println("pubkey ", frec)

	// nameserver, _ := net.LookupNS(domain)
	// for _, ns := range nameserver {
	// 	fmt.Println(ns)
	// }

	// cname, _ := net.LookupCNAME(domain)
	// fmt.Println(cname)

	iprecords, _ := net.LookupIP(domain)
	for _, ip := range iprecords {
		fmt.Println("ip ", ip)
	}
}

func getConfig() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	return config
}

func readOption() string {
	optionPtr := flag.String("option", "", "the command to be performed")
	flag.Parse()
	option := *optionPtr
	log.Println("run client with option:", option)
	return option
}

func testclient_subscribe() {
	time.Sleep(200 * time.Millisecond)
	addr := ":" + strconv.Itoa(node_port)
	log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("cant run")
		return
	}

	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr)

	go ntcl.ReadLoop(ntchan)

	//subscribe example
	//reqs := "REQ#PING#|"
	reqs := "REQ#SUBTO#TIME|"
	log.Println("subscribe")
	ntcl.NetWrite(ntchan, reqs)

	//log.Println(clientNt.SrcName)

	go func() {
		for {
			rmsg := <-ntchan.Reader_queue
			log.Println("> ", rmsg)
		}
	}()
	time.Sleep(2000 * time.Millisecond)

	reqs = "REQ#SUBUN#TIME|"
	ntcl.NetWrite(ntchan, reqs)

	log.Println("unsubscribe")

	time.Sleep(10000 * time.Millisecond)

	//defer conn.Close()
	return

}

//run client based on options
func main() {

	config := getConfig()

	option := readOption()

	//dnslook()

	switch option {

	case "test", "ping", "heartbeat", "getbalance", "faucet", "txpool", "pushtx", "randomtx":
		runSingleMode(option, config)

	case "createkeys", "sign", "createtx", "verify":
		runOffline(option, config)

	case "pingall", "blockheight":
		runPeermode(option, config)

	case "listen":
		runListenMode(option, config)

	default:
		log.Println("unknown option")
	}

}
