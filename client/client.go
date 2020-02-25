package main

//client based application to request information from peers

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
	cryptoutil "github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/ntwk"
)

var Peers []protocol.Peer

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
}

func addPeerOut(p protocol.Peer) {
	Peers = append(Peers, p)
	log.Println("peers now", Peers)
}

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End measure time:", s, "took", endTime.Sub(startTime))
}

//
func MakeRandomTx(peer protocol.Peer) error {
	//make a random transaction by requesting random account from node
	//get random account

	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_RANDOM_ACCOUNT, "emptydata")

	response := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	var a block.Account
	dataBytes := []byte(response.Data)
	if err := json.Unmarshal(dataBytes, &a); err != nil {
		panic(err)
	}
	log.Print(" account key ", a.AccountKey)

	//use this random account to send coins from

	//send Tx
	testTx := protocol.RandomTx(a)
	txJson, _ := json.Marshal(testTx)
	log.Println("txJson ", txJson)

	req_msg = protocol.EncodeMessageTx(txJson)

	response = protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Print("response msg ", response)

	return nil
}

func CreateTx(peer protocol.Peer) {
	// keypair := crypto.PairFromSecret("test")
	// var tx block.Tx
	// s := block.AccountFromString("Pa033f6528cc1")
	// r := s //TODO
	// tx = block.Tx{Nonce: 0, Amount: 0, Sender: s, Receiver: r}
	// signature := crypto.SignTx(tx, keypair)
	// sighex := hex.EncodeToString(signature.Serialize())
	// tx.Signature = sighex
	// tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

}

func PushTx(peer protocol.Peer) error {

	dat, _ := ioutil.ReadFile("tx.json")
	var tx block.Tx

	if err := json.Unmarshal(dat, &tx); err != nil {
		panic(err)
	}

	//send Tx
	txJson, _ := json.Marshal(tx)
	log.Println("txJson ", string(txJson))

	req_msg := protocol.EncodeMessageTx(txJson)
	response := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Print("reply msg ", response)

	return nil
}

func Getbalance(peer protocol.Peer) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	txJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_BALANCE, string(txJson))
	response := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Println("response ", response)
	var balance int
	if err := json.Unmarshal(response.Data, &balance); err != nil {
		panic(err)
	}
	log.Println("balance of account ", balance)

	return nil
}

func Getblockheight(peer protocol.Peer) error {
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_BLOCKHEIGHT, "")
	response := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	var blockheight int
	if err := json.Unmarshal(response.Data, &blockheight); err != nil {
		panic(err)
	}
	log.Println("blockheight ", blockheight)

	return nil
}

func Gettxpool(peer protocol.Peer) error {
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_GETTXPOOL, "")
	log.Println("> ", req_msg)
	resp := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	log.Println("rcvmsg ", resp)
	log.Println("data ", resp.Data)

	var txp []block.Tx
	if err := json.Unmarshal(resp.Data, &txp); err != nil {
		panic(err)
	}
	log.Println("txp ", txp)

	return nil
}

func GetFaucet(peer protocol.Peer) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	accountJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_FAUCET, string(accountJson))
	resp := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Println("resp ", resp)

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

//setup connection to a peer from client side for requests
func setupPeerClient(peer protocol.Peer) error {

	rw, err := protocol.OpenOut(peer.Address, peer.NodePort)
	if err != nil {
		log.Println("error open ", err)
		return err
	}

	go protocol.RequestLoop(rw, peer.Req_chan, peer.Rep_chan)
	return nil
}

func readKeys(keysfile string) cryptoutil.Keypair {

	dat, _ := ioutil.ReadFile(keysfile)
	s := strings.Split(string(dat), string("\n"))

	pubkeyHex := s[0]
	log.Println("pub ", pubkeyHex)

	privHex := s[1]
	log.Println("privHex ", privHex)

	return cryptoutil.Keypair{PubKey: cryptoutil.PubKeyFromHex(pubkeyHex), PrivKey: cryptoutil.PrivKeyFromHex(privHex)}
}

func writeKeys(kp cryptoutil.Keypair, keysfile string) {

	pubkeyHex := cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println("pub ", pubkeyHex)

	privHex := cryptoutil.PrivKeyToHex(kp.PrivKey)
	log.Println("privHex ", privHex)

	address := cryptoutil.Address(pubkeyHex)

	t := pubkeyHex + "\n" + privHex + "\n" + address
	//log.Println("address ", address)
	ioutil.WriteFile(keysfile, []byte(t), 0644)
}

func createKeys() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.Trim(pw, string('\n'))
	fmt.Println(pw)

	//check if exists
	//dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)

	kp := cryptoutil.PairFromSecret(pw)
	log.Println("keypair ", kp)

	writeKeys(kp, "keys.txt")

}

func createtx() {
	kp := readKeys("keys.txt")

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter password: ")
	// pw, _ := reader.ReadString('\n')
	// pw = strings.Trim(pw, string('\n'))

	// keypair := crypto.PairFromSecret(pw)

	pubk := cryptoutil.PubKeyToHex(kp.PubKey)
	addr := cryptoutil.Address(pubk)
	s := block.AccountFromString(addr)
	log.Println("using account ", s)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter amount: ")
	amount, _ := reader.ReadString('\n')
	amount = strings.Trim(amount, string('\n'))
	amount_int, _ := strconv.Atoi(amount)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter recipient: ")
	recv, _ := reader.ReadString('\n')
	recv = strings.Trim(recv, string('\n'))

	tx := block.Tx{Nonce: 1, Amount: amount_int, Sender: block.Account{AccountKey: addr}, Receiver: block.Account{AccountKey: recv}}
	log.Println("tx ", tx)

	signature := cryptoutil.SignTx(tx, kp)
	sighex := hex.EncodeToString(signature.Serialize())

	tx.Signature = sighex
	tx.SenderPubkey = cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println("tx ", tx)

	txJson, _ := json.Marshal(tx)
	// //write to file
	// log.Println(txJson)

	ioutil.WriteFile("tx.json", []byte(txJson), 0644)
}

func setupAllPeers(config Configuration) {

	for _, peerAddress := range config.PeerAddresses {
		log.Println("setup  peer ", peerAddress)
		p := protocol.CreatePeer(peerAddress, config.NodePort)

		err := setupPeerClient(p)
		if err != nil {
			log.Println("connect failed")
			continue
		} else {
			protocol.MakePing(p)
		}
	}

}

//run client against multiple nodes
func runPeermode(option string, config Configuration) {
	log.Println("runPeermode")

	log.Println("setup peers")

	for _, peerAddress := range config.PeerAddresses {

		p := protocol.CreatePeer(peerAddress, config.NodePort)
		log.Println("add peer ", p)

		err := setupPeerClient(p)
		if err != nil {
			//remove peer
			log.Println("dont add peer to list")
		} else {
			addPeerOut(p)
		}
	}

	switch option {

	case "pingall":

		defer track(runningtime("execute ping"))
		successCount := 0
		for _, peerAddress := range config.PeerAddresses {
			log.Println("setup  peer ", peerAddress, config.NodePort)
			p := protocol.CreatePeer(peerAddress, config.NodePort)

			err := setupPeerClient(p)
			if err != nil {
				log.Println("connect failed")
				continue
			} else {
				success := protocol.MakePing(p)
				if success {
					successCount++
				}
			}
		}

		log.Println("pinged peers ", len(config.PeerAddresses), " successCount:", successCount)

	case "blockheight":

		for _, peerAddress := range config.PeerAddresses {
			log.Println("setup  peer ", peerAddress)
			p := protocol.CreatePeer(peerAddress, config.NodePort)

			err := setupPeerClient(p)
			if err == nil {
				log.Println("block height ", p)
				Getblockheight(p)
			}
		}

	}

}

func heartbeat() {

}

//run client against single node, just use first IP address in peers i.e. mainpeer
func runSingleMode(option string, config Configuration) {

	mainPeerAddress := config.PeerAddresses[0]
	log.Println("setup main peer ", mainPeerAddress, config.NodePort)
	mainPeer := protocol.CreatePeer(mainPeerAddress, config.NodePort)
	log.Println("client with mainPeer ", mainPeer)

	setupPeerClient(mainPeer)

	switch option {

	case "ping":
		log.Println("ping")
		protocol.MakePing(mainPeer)

	case "heartbeat":
		log.Println("heartbeat")
		success := protocol.MakeHandshake(mainPeer)
		if success {
			log.Println("start heartbeat")
			protocol.Hearbeat(mainPeer)
		}

	case "getbalance":
		log.Println("getbalance")

		Getbalance(mainPeer)

	case "blockheight":
		log.Println("blockheight")

		Getblockheight(mainPeer)

	case "faucet":
		log.Println("faucet")
		//get coins
		//GetFaucet(rw)
		GetFaucet(mainPeer)

	case "txpool":
		_ = Gettxpool(mainPeer)
		return

	case "pushtx":
		_ = PushTx(mainPeer)
		return

	case "randomtx":
		_ = MakeRandomTx(mainPeer)
		return
	}
}

//client that only listens to events
func runListenMode(option string, config Configuration) {
	log.Println("listen")

	mainPeerAddress := config.PeerAddresses[0]
	log.Println("setup main peer ", mainPeerAddress)
	mainPeer := protocol.CreatePeer(mainPeerAddress, config.NodePort)
	success := protocol.MakeHandshake(mainPeer)
	log.Println(success)
	// log.Println("start heartbeat")
	// if success {
	// 	hTime := 2000 * time.Millisecond

	// 	for _ = range time.Tick(hTime) {
	// 		//log.Println(x)
	// 		protocol.Hearbeat(mainPeer)
	// 	}

	// }
}

//run client without client or server
func runOffline(option string, config Configuration) {

	switch option {
	case "createkeys":
		createKeys()

	case "readkeys":
		kp := readKeys("keys.txt")
		log.Println(kp)

	case "sign":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to sign: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		fmt.Println(msg)
		kp := readKeys("keys.txt")
		signature := cryptoutil.SignMsgHash(kp, msg)
		log.Println("signature ", signature)

		sighex := hex.EncodeToString(signature.Serialize())
		log.Println("sighex ", sighex)

	case "createtx":
		createtx()

	case "verify":
		reader := bufio.NewReader(os.Stdin)
		log.Print("Enter message to verify: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		log.Println(msg)

		fmt.Print("Enter signature to verify: ")
		msgsig, _ := reader.ReadString('\n')
		msgsig = strings.Trim(msgsig, string('\n'))

		sign := cryptoutil.SignatureFromHex(msgsig)

		fmt.Print("Enter pubkey to verify: ")
		msgpub, _ := reader.ReadString('\n')
		log.Println(msgpub)
		msgpub = strings.Trim(msgpub, string('\n'))

		pubkey := cryptoutil.PubKeyFromHex(msgpub)

		verified := cryptoutil.VerifyMessageSignPub(sign, pubkey, msg)
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

//run client based on options
func main() {

	//config := getConfig()

	//option := readOption()

	dnslook()

	// switch option {

	// case "ping", "heartbeat", "getbalance", "faucet", "txpool", "pushtx", "randomtx":
	// 	runSingleMode(option, config)

	// case "createkeys", "sign", "createtx", "verify":
	// 	runOffline(option, config)

	// case "pingall", "blockheight":
	// 	runPeermode(option, config)

	// case "listen":
	// 	runListenMode(option, config)

	// default:
	// 	log.Println("unknown option")
	// }

}
