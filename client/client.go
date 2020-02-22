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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/polygonledger/node/block"
	//"github.com/polygonledger/node/crypto"
	cryptoutil "github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/ntwk"
)

var Peers []protocol.Peer

type Configuration struct {
	PeerAddresses []string
	NodePort      int
	WebPort       int
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
	log.Println(response)
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

func MakePing(peer protocol.Peer) {
	emptydata := ""
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_PING, emptydata)
	resp := protocol.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	log.Println("resp ", resp)
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

func client(ip string, NodePort int) *bufio.ReadWriter {

	// Open a connection to the server.
	log.Println(">>> ", ip, NodePort)
	rw := protocol.OpenOut(ip, NodePort)

	return rw
}

//setup connection to a peer from client side for requests
func setupPeerClient(peer protocol.Peer) {

	rw := client(peer.Address, peer.NodePort)

	go protocol.RequestLoop(rw, peer.Req_chan, peer.Rep_chan)
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
	log.Println(tx)

	signature := cryptoutil.SignTx(tx, kp)
	sighex := hex.EncodeToString(signature.Serialize())

	tx.Signature = sighex
	tx.SenderPubkey = cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println(tx)

	txJson, _ := json.Marshal(tx)
	// //write to file
	// log.Println(txJson)

	ioutil.WriteFile("tx.json", []byte(txJson), 0644)
}

func createPeer(ipAddress string, NodePort int) protocol.Peer {
	//addr := ip
	p := protocol.Peer{Address: ipAddress, NodePort: NodePort, Req_chan: make(chan protocol.Message), Rep_chan: make(chan protocol.Message), Out_req_chan: make(chan protocol.Message), Out_rep_chan: make(chan protocol.Message)}
	return p
}

//run client based on options
func main() {

	//read config

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("PeerAddresses: ", configuration.PeerAddresses)

	optionPtr := flag.String("option", "createkeypair", "the command to be performed")
	flag.Parse()
	//fmt.Println("option:", *optionPtr)

	mainPeerAddress := configuration.PeerAddresses[0]
	log.Println("setup peer ", mainPeerAddress)
	mainPeer := createPeer(mainPeerAddress, configuration.NodePort)
	log.Println("client to mainPeer ", mainPeer)

	setupPeerClient(mainPeer)

	if *optionPtr == "createkeys" {
		createKeys()

	} else if *optionPtr == "readkeys" {
		kp := readKeys("keys.txt")
		log.Println(kp)

	} else if *optionPtr == "sign" {
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

	} else if *optionPtr == "createtx" {

		createtx()

	} else if *optionPtr == "verify" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to verify: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		fmt.Println(msg)

		fmt.Print("Enter signature to verify: ")
		msgsig, _ := reader.ReadString('\n')
		msgsig = strings.Trim(msgsig, string('\n'))
		//examplesig := "3045022100dd2781cc37edb84c5ed21b3d8fc03d49ebddf5647d23a9132eeea8bd2b951bd1022041519c47b77803d528d1b428ccb4d84a90ce3b67a22662d5feaa84c4521e5759"

		sign := cryptoutil.SignatureFromHex(msgsig)

		fmt.Print("Enter pubkey to verify: ")
		msgpub, _ := reader.ReadString('\n')
		fmt.Println(msgpub)
		msgpub = strings.Trim(msgpub, string('\n'))

		//msgpub = "039f6095ba1afa34c437a88fceb444bf177326eb9222d4938336387ecb4cbe7234"

		pubkey := cryptoutil.PubKeyFromHex(msgpub)

		verified := cryptoutil.VerifyMessageSignPub(sign, pubkey, msg)
		log.Println("verified ", verified)

	} else if *optionPtr == "createkeypair" {
		//TODO read secret from cmd
		kp := cryptoutil.PairFromSecret("test")
		log.Println("keypair ", kp)

	} else if *optionPtr == "ping" {
		log.Println("ping")
		//protocol.Server_address
		MakePing(mainPeer)

	} else if *optionPtr == "pingall" {
		//protocol.Server_address
		for _, peer := range configuration.PeerAddresses {
			log.Println("ping ", peer)

		}
		MakePing(mainPeer)

	} else if *optionPtr == "pingconnect" {
		log.Println("ping continously")
		//protocol.Server_address
		for {
			MakePing(mainPeer)
			time.Sleep(10 * time.Second)
		}

	} else if *optionPtr == "getbalance" {
		log.Println("getbalance")
		//protocol.Server_address

		Getbalance(mainPeer)

	} else if *optionPtr == "blockheight" {
		log.Println("blockheight")

		Getblockheight(mainPeer)

	} else if *optionPtr == "faucet" {
		log.Println("faucet")
		//get coins
		//GetFaucet(rw)
		GetFaucet(mainPeer)

	} else if *optionPtr == "txpool" {
		err = Gettxpool(mainPeer)
		return

	} else if *optionPtr == "pushtx" {
		err = PushTx(mainPeer)
		return

	} else if *optionPtr == "randomtx" {
		err = MakeRandomTx(mainPeer)
		return
	}
	// else if *optionPtr == "pushtx" {
	// 	//read locally created tx file and push it to server
	// 	data, _ := ioutil.ReadFile("tx.json")
	// 	log.Println(string(data))
	// 	var tx block.Tx
	// 	if err := json.Unmarshal(data, &tx); err != nil {
	// 		panic(err)
	// 	}
	// 	log.Println(">> ", tx)

	// 	return
	// }

}
