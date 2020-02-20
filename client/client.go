package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/ntwk"
)

//
func MakeRandomTx(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {
	//make a random transaction by requesting random account from node
	//get random account

	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_RANDOM_ACCOUNT, "emptydata")

	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

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

	response = protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	log.Print("response msg ", response)

	return nil
}

func CreateTx() {
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

func PushTx(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {

	dat, _ := ioutil.ReadFile("tx.json")
	var tx block.Tx

	if err := json.Unmarshal(dat, &tx); err != nil {
		panic(err)
	}

	//send Tx
	txJson, _ := json.Marshal(tx)
	log.Println("txJson ", string(txJson))

	req_msg := protocol.EncodeMessageTx(txJson)
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	log.Print("reply msg ", response)

	return nil
}

func Getbalance(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	txJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_BALANCE, string(txJson))
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	log.Println(response)
	var balance int
	if err := json.Unmarshal(response.Data, &balance); err != nil {
		panic(err)
	}
	log.Println("balance of account ", balance)

	return nil
}

func Getblockheight(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_BLOCKHEIGHT, "")
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

	var blockheight int
	if err := json.Unmarshal(response.Data, &blockheight); err != nil {
		panic(err)
	}
	log.Println("blockheight ", blockheight)

	return nil
}

func Gettxpool(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_GETTXPOOL, "")
	log.Println("> ", req_msg)
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

	log.Println("rcvmsg ", response)

	return nil
}

func GetFaucet(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	accountJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_FAUCET, string(accountJson))
	protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	//log.Println("resp ", resp)

	return nil
}

func MakePing(msg_in_chan chan protocol.Message, msg_out_chan chan protocol.Message) {
	emptydata := ""
	req_msg := protocol.EncodeMsgString(protocol.REQ, protocol.CMD_PING, emptydata)
	resp := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	log.Println("resp ", resp)
}

func client(ip string) (*bufio.ReadWriter, error) {

	// Open a connection to the server.
	rw, err := protocol.Open(ip + protocol.Port)
	//log.Println(rw)
	if err != nil {
		return nil, errors.Wrap(err, "Client: Failed to open connection to "+ip+protocol.Port)
	}

	return rw, err
}

type Configuration struct {
	PeerAddresses []string
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

//start client and connect to the host
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

	//if exists
	// dat, _ := ioutil.ReadFile("keys.txt")
	// //check(err)
	// fmt.Print("keys ", string(dat))

	optionPtr := flag.String("option", "createkeypair", "the command to be performed")
	flag.Parse()
	fmt.Println("option:", *optionPtr)

	mainPeer := configuration.PeerAddresses[0]
	rw, err := client(mainPeer)
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	req_chan := make(chan protocol.Message)
	rep_chan := make(chan protocol.Message)
	go protocol.RequestLoop(rw, req_chan, rep_chan)

	if *optionPtr == "createkeypair" {
		//TODO read secret from cmd
		kp := crypto.PairFromSecret("test")
		log.Println("keypair ", kp)

	} else if *optionPtr == "ping" {
		log.Println("ping")
		//protocol.Server_address
		MakePing(req_chan, rep_chan)

	} else if *optionPtr == "pingall" {
		//protocol.Server_address
		for _, peer := range configuration.PeerAddresses {
			log.Println("ping ", peer)

		}
		MakePing(req_chan, rep_chan)

	} else if *optionPtr == "pingconnect" {
		log.Println("ping continously")
		//protocol.Server_address
		for {
			MakePing(req_chan, rep_chan)
			time.Sleep(10 * time.Second)
		}

	} else if *optionPtr == "getbalance" {
		log.Println("getbalance")
		//protocol.Server_address

		Getbalance(req_chan, rep_chan)

	} else if *optionPtr == "blockheight" {
		log.Println("blockheight")

		Getblockheight(req_chan, rep_chan)

	} else if *optionPtr == "faucet" {
		log.Println("faucet")
		//get coins
		//GetFaucet(rw)
		GetFaucet(req_chan, rep_chan)

	} else if *optionPtr == "txpool" {
		err = Gettxpool(req_chan, rep_chan)
		return

	} else if *optionPtr == "pushtx" {
		err = PushTx(req_chan, rep_chan)
		return

	} else if *optionPtr == "randomtx" {
		err = MakeRandomTx(req_chan, rep_chan)
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
