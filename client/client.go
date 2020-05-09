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
	"github.com/polygonledger/node/parser"
)

var Peers []ntcl.Peer

const node_port = 8888

const keysfile = "keys.wfe"

// --- utils ---

func ReadKeys(keysfile string) crypto.Keypair {

	log.Println("read keys from ", keysfile)
	dat, _ := ioutil.ReadFile(keysfile)
	s := string(dat)
	vs, _ := parser.ReadMapP(s)

	privHex := parser.StringUnWrap(vs[0])
	pubkeyHex := parser.StringUnWrap(vs[1])
	log.Println("pub ", pubkeyHex)

	return crypto.Keypair{PubKey: crypto.PubKeyFromHex(pubkeyHex), PrivKey: crypto.PrivKeyFromHex(privHex)}
}

func CreateKeypairFormat(privkey string, pubkey_string string, address string) string {
	mp := map[string]string{"privkey": parser.StringWrap(privkey), "pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
	m := parser.MakeMap(mp)
	return m
}

func CreatePubKeypairFormat(pubkey_string string, address string) string {
	mp := map[string]string{"pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
	m := parser.MakeMap(mp)
	return m
}

//write keys to file with format
func WriteKeys(kp crypto.Keypair, keysfile string) {
	pubkeyHex := crypto.PubKeyToHex(kp.PubKey)
	privHex := crypto.PrivKeyToHex(kp.PrivKey)
	address := crypto.Address(pubkeyHex)
	s := CreateKeypairFormat(privHex, pubkeyHex, address)
	ioutil.WriteFile(keysfile, []byte(s), 0644)
}

func WritePubKeys(kp crypto.Keypair, keysfile string) {
	pubkeyHex := crypto.PubKeyToHex(kp.PubKey)
	address := crypto.Address(pubkeyHex)
	s := CreatePubKeypairFormat(pubkeyHex, address)
	ioutil.WriteFile(keysfile, []byte(s), 0644)
}

func CreateKeys() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.Trim(pw, string('\n'))

	//check if exists
	//dat, _ := ioutil.ReadFile(keysfile)
	//check(err)

	kp := crypto.PairFromSecret(pw)

	WriteKeys(kp, keysfile)

	WritePubKeys(kp, "pubkeys.wfe")

}

func Createtx() {
	// kp := ReadKeys(keysfile)

	// pubk := crypto.PubKeyToHex(kp.PubKey)
	// addr := crypto.Address(pubk)
	// s := block.AccountFromString(addr)
	// log.Println("using account ", s)

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter amount: ")
	// amount, _ := reader.ReadString('\n')
	// amount = strings.Trim(amount, string('\n'))
	// amount_int, _ := strconv.Atoi(amount)

	// reader = bufio.NewReader(os.Stdin)
	// fmt.Print("Enter recipient: ")
	// recv, _ := reader.ReadString('\n')
	// recv = strings.Trim(recv, string('\n'))

	// tx := block.Tx{Nonce: 1, Amount: amount_int, Sender: block.Account{AccountKey: addr}, Receiver: block.Account{AccountKey: recv}}
	// log.Println("tx ", tx)

	// signature := crypto.SignTx(tx, kp)
	// sighex := hex.EncodeToString(signature.Serialize())

	// tx.Signature = sighex
	// tx.SenderPubkey = crypto.PubKeyToHex(kp.PubKey)
	// log.Println("tx ", tx)

	// txJson, _ := json.Marshal(tx)
	// // //write to file
	// // log.Println(txJson)
	// f := "tx.json"
	// ioutil.WriteFile(f, []byte(txJson), 0644)

	// log.Println("wrote to " + f)
}

//
func MakeRandomTx(peer ntcl.Peer) {
	//make a random transaction by requesting random account from node
	//get random account

	// req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_RANDOM_ACCOUNT, "emptydata")

	// response := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	// var a block.Account
	// dataBytes := []byte(response.Data)
	// if err := json.Unmarshal(dataBytes, &a); err != nil {
	// 	panic(err)
	// }
	// log.Print(" account key ", a.AccountKey)

	// //use this random account to send coins from

	// //send Tx
	// testTx := ntcl.RandomTx(a)
	// txJson, _ := json.Marshal(testTx)
	// log.Println("txJson ", txJson)

	// req_msg = ntcl.EncodeMessageTx(txJson)

	// response = ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	// log.Print("response msg ", response)

	// return nil
}

//handlers TODO this is higher level and should be somewhere else
func RandomTx(account_s block.Account) {

	// rand.Seed(time.Now().UnixNano())
	// randNonce := rand.Intn(100)

	// kp := crypto.PairFromSecret("test111")
	// log.Println("PUBKEY ", kp.PubKey)

	// r := crypto.RandomPublicKey()
	// address_r := crypto.Address(r)
	// account_r := block.AccountFromString(address_r)

	// //TODO make sure the amount is covered by sender
	// rand.Seed(time.Now().UnixNano())
	// randomAmount := rand.Intn(20)

	// // log.Printf("randomAmount %s", randomAmount)
	// // log.Printf("randNonce %d", randNonce)
	// testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	// sig := crypto.SignTx(testTx, kp)
	// sighex := hex.EncodeToString(sig.Serialize())
	// testTx.Signature = sighex
	// log.Println(">> ran tx", testTx.Signature)
	// return testTx
}

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
	verbose       bool
}

func addPeerOut(p ntcl.Peer) {
	Peers = append(Peers, p)
	log.Println("peers now", Peers)
}

func initClient(mainPeerAddress string, verbose bool) ntcl.Ntchan {

	addr := mainPeerAddress + ":" + strconv.Itoa(node_port)
	log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("cant run")
		//return
	}

	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr, verbose)

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

func PushTx(ntchan ntcl.Ntchan) error {
	log.Println("PushTx")

	dat, _ := ioutil.ReadFile("tx.json")
	var tx block.Tx

	if err := json.Unmarshal(dat, &tx); err != nil {
		panic(err)
	}

	//send Tx
	txJson, _ := json.Marshal(tx)
	log.Println("txJson ", string(txJson))

	req_msg := ntcl.EncodeMessageTx(txJson)
	log.Print(" ", req_msg)

	ntchan.REQ_out <- req_msg

	rep := <-ntchan.REP_in
	log.Println("reply ", rep)

	return nil
}

func fetchbalance(ntchan ntcl.Ntchan, addr string) {
	req_msg := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_BALANCE, string(addr))
	log.Println(req_msg)

	ntchan.REQ_out <- req_msg

	// rep := <-ntchan.REP_in
	// //log.Println("reply ", rep)
	// rep = strings.Trim(rep, string(ntcl.DELIM))
	// s := strings.Split(rep, string(ntcl.DELIM_HEAD))
	// balance_int, _ := strconv.Atoi(string(s[2]))
	// log.Println(fmt.Sprintf("balance of %s %d", addr, balance_int))
}

func Mybalance(ntchan ntcl.Ntchan) error {

	kp := ReadKeys(keysfile)
	pubk := crypto.PubKeyToHex(kp.PubKey)
	myaddr := crypto.Address(pubk)

	log.Println("request balance for my address ", myaddr)
	//json, _ := json.Marshal(block.Account{AccountKey: myaddr})
	fetchbalance(ntchan, myaddr)

	//log.Println("reply ", strconv.Atoi(int(msg.Data)))

	return nil
}

func Getbalance(ntchan ntcl.Ntchan) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	fetchbalance(ntchan, addr)
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
	req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_BLOCKHEIGHT)
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
	req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_GETTXPOOL)
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

//run client against multiple nodes
func runPeermode(cmd string, config Configuration) {
	log.Println("runPeermode")

	log.Println("setup peers")

	for _, peerAddress := range config.PeerAddresses {

		//p := ntcl.CreatePeer(peerAddress, config.NodePort)
		//log.Println("add peer ", p)

		ntchan := initClient(peerAddress, config.verbose)

		//
		p := ntcl.CreatePeer(peerAddress, peerAddress, config.NodePort, ntchan)

		//log.Println("init  ", ntchan)
		ping(p)

		//err := setupPeerClient(p)
		// if err != nil {
		// 	//remove peer
		// 	log.Println("dont add peer to list")
		// } else {
		// 	addPeerOut(p)
		// }
	}

	// switch cmd {

	// case "pingall":

	// 	defer track(runningtime("execute ping"))
	// 	successCount := 0

	// 	log.Println("pinged peers ", len(config.PeerAddresses), " successCount:", successCount)

	// 	// case "blockheight":

	// 	// 	for _, peerAddress := range config.PeerAddresses {
	// 	// 		log.Println("setup  peer ", peerAddress)
	// 	// 	}

	// }

}

//TODO common
func ping(peer ntcl.Peer) {

	// req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, "")

	//subscribe example
	//reqs := "REQ#PING#|"
	req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_PING)
	peer.NTchan.REQ_out <- req_msg
	//ntcl.NetWrite(ntchan, reqs)

	time.Sleep(1000 * time.Millisecond)

	reply := <-peer.NTchan.REP_in
	success := reply == "{:REP PONG}"
	log.Println("success ", success)

}

func request_reply(peer ntcl.Peer, req_msg string) {
	peer.NTchan.REQ_out <- req_msg
	time.Sleep(1000 * time.Millisecond)
	reply := <-peer.NTchan.REP_in
	log.Println("reply ", reply)
}

func status(peer ntcl.Peer) {
	req_msg := "{:REQ STATUS}"
	request_reply(peer, req_msg)
}

func MakeFaucet(ntchan ntcl.Ntchan) {
	log.Println("read keys")
	kp := ReadKeys(keysfile)
	pubk := crypto.PubKeyToHex(kp.PubKey)

	addr := crypto.Address(pubk)
	log.Println("request faucet to ", addr)
	req_msg := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_FAUCET, addr)

	ntchan.REQ_out <- req_msg

	rep := <-ntchan.REP_in
	log.Println("reply ", rep)
	log.Println("wait for block....")
	time.Sleep(5000 * time.Millisecond)

	req_msg2 := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_BALANCE, addr)
	ntchan.REQ_out <- req_msg2

	rep2 := <-ntchan.REP_in

	log.Println("reply ", rep2)

}

//run client against single node, just use first IP address in peers i.e. mainpeer
func runSingleMode(cmd string, config Configuration) {

	mainPeerAddress := config.PeerAddresses[0]
	ntchan := initClient(mainPeerAddress, config.verbose)
	p := ntcl.CreatePeer(mainPeerAddress, mainPeerAddress, config.NodePort, ntchan)
	log.Println("init ", ntchan)

	switch cmd {

	case "ping":
		log.Println("ping")
		//ntchan := initClient(config)
		ping(p)

		time.Sleep(100 * time.Millisecond)

	case "status":
		status(p)

	case "faucet":
		MakeFaucet(ntchan)

	case "faucetloop":
		for {
			MakeFaucet(ntchan)
			time.Sleep(10 * time.Second)
		}

	case "pushtx":
		PushTx(ntchan)

	case "mybalance":
		Mybalance(ntchan)

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

	case "getbalance":
		log.Println("getbalance")
		Getbalance(ntchan)

	case "dnslook":
		dnslook()

		// case "blockheight":
		// 	log.Println("blockheight")

		// 	Getblockheight(mainPeer)

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
func runListenMode(cmd string, config Configuration) {
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

//TODO move
func createsigmap(pubkey_string string, txsighex string) string {
	return `{:SenderPubkey "` + pubkey_string + `" :Signature "` + txsighex + `"}`
}

func txVector(simpletx string, sigmap string) string {
	return `[:STX ` + simpletx + ` ` + sigmap + ` ]`
}

//run client without client or server
func runOffline(cmd string, config Configuration) {

	switch cmd {
	case "createkeys":
		CreateKeys()

	case "readkeys":
		kp := ReadKeys(keysfile)
		log.Println(kp)

	case "signtx":
		fmt.Println("signtx")

		kp := ReadKeys(keysfile)
		dat, _ := ioutil.ReadFile("example.txp")
		msg := string(dat)
		fmt.Println("sign >> ", msg)
		signature := crypto.SignMsgHash(kp, msg)
		log.Println("signature ", signature)

		pubkey_string := crypto.PubKeyToHex(kp.PubKey)
		sigmap := createsigmap(pubkey_string, msg)
		fmt.Println("sigmap ", sigmap)

		v := txVector(msg, sigmap)
		ioutil.WriteFile("example.txo", []byte(v), 0644)

		// sighex := hex.EncodeToString(signature.Serialize())
		// log.Println("sighex ", sighex)

	case "sign":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to sign: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		//fmt.Println(msg)
		kp := ReadKeys(keysfile)
		signature := crypto.SignMsgHash(kp, msg)
		//log.Println("signature ", signature)

		sighex := hex.EncodeToString(signature.Serialize())
		log.Println("signature as hex ", sighex)

		//sigmap := CreateSigmap(pubk, sig)

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

	case "verifytx":

		dat, _ := ioutil.ReadFile("example.txo")
		msg := string(dat)
		fmt.Println("to verify ", msg)

	}
}

//dns functions for later, as we can use txt records to get pubkey
func dnslook() {
	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter Domain to check ")

	domain, _ := reader.ReadString('\n')
	domain = strings.Trim(domain, string('\n'))
	log.Println(domain)
	//domain = "polygonnode.com"

	txtrecords, _ := net.LookupTXT(domain)
	log.Println("txtrecords ", txtrecords)

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

func readFlags() string {
	cPtr := flag.String("cmd", "", "the command to be performed")
	flag.Parse()
	cmd := *cPtr
	log.Println("run client with cmd:", cmd)
	return cmd
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
	ntchan := ntcl.ConnNtchan(conn, "client", addr, true)

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

//run client based on cmds
func main() {

	//kp := ReadKeys(keysfile)
	//fmt.Println("=> ", kp)

	config := getConfig()

	cmd := readFlags()

	//dnslook()

	switch cmd {

	case "test", "ping", "status", "heartbeat", "getbalance", "faucet", "faucetloop", "txpool", "pushtx", "randomtx", "mybalance", "dnslook":
		runSingleMode(cmd, config)

	case "createkeys", "sign", "signtx", "createtx", "verify", "verifytx":
		runOffline(cmd, config)

	case "pingall", "blockheight":
		runPeermode(cmd, config)

	case "listen":
		runListenMode(cmd, config)

	default:
		log.Println("unknown cmd")
	}

}
