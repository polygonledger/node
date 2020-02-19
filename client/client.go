package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/ntwk"
)

//
func MakeRandomTx(msg_in_chan chan string, msg_out_chan chan string) error {
	//make a random transaction by requesting random account from node
	//get random account

	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_RANDOM_ACCOUNT, "emptydata")

	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

	var a block.Account
	dataBytes := []byte(response)
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

func PushTx(msg_in_chan chan string, msg_out_chan chan string) error {

	// keypair := crypto.PairFromSecret("test")
	// var tx block.Tx
	// s := block.AccountFromString("Pa033f6528cc1")
	// r := s //TODO
	// tx = block.Tx{Nonce: 0, Amount: 0, Sender: s, Receiver: r}
	// signature := crypto.SignTx(tx, keypair)
	// sighex := hex.EncodeToString(signature.Serialize())
	// tx.Signature = sighex
	// tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

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

func Getbalance(msg_in_chan chan string, msg_out_chan chan string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	txJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_BALANCE, string(txJson))
	//fmt.Println("msg ", msg)
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

	x, _ := strconv.Atoi(response)
	log.Println("balance of account ", x)

	return nil
}

func Gettxpool(msg_in_chan chan string, msg_out_chan chan string) error {
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_GETTXPOOL, "")
	log.Println("> ", req_msg)
	response := protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)

	log.Println("rcvmsg ", response)

	return nil
}

func GetFaucet(msg_in_chan chan string, msg_out_chan chan string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	accountJson, _ := json.Marshal(block.Account{AccountKey: addr})
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_FAUCET, string(accountJson))
	protocol.RequestReplyChan(req_msg, msg_in_chan, msg_out_chan)
	//log.Println("resp ", resp)

	return nil
}

func MakePing(msg_in_chan chan string, msg_out_chan chan string) {
	emptydata := ""
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_PING, emptydata)
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
	ServerAddress string
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

	//prepare to run client

	//read config

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("ServerAddress: ", configuration.ServerAddress)

	//if exists
	// dat, _ := ioutil.ReadFile("keys.txt")
	// //check(err)
	// fmt.Print("keys ", string(dat))

	optionPtr := flag.String("option", "createkeypair", "the command to be performed")
	flag.Parse()
	fmt.Println("option:", *optionPtr)

	rw, err := client(configuration.ServerAddress)
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	// //sub loop
	// for {
	// 	log.Println("..")
	// 	msg, _ := rw.ReadString(protocol.DELIM)
	// 	log.Println(msg)
	// 	msg = strings.Trim(msg, string(protocol.DELIM))
	// 	time.Sleep(2 * time.Second)
	// }

	msg_in_chan := make(chan string)
	msg_out_chan := make(chan string)
	go protocol.RequestLoop(rw, msg_in_chan, msg_out_chan)

	if *optionPtr == "createkeypair" {
		//TODO read secret from cmd
		kp := crypto.PairFromSecret("test")
		log.Println("keypair ", kp)

	} else if *optionPtr == "ping" {
		log.Println("ping")
		//protocol.Server_address
		MakePing(msg_in_chan, msg_out_chan)
	} else if *optionPtr == "pingconnect" {
		log.Println("ping continously")
		//protocol.Server_address
		for {
			MakePing(msg_in_chan, msg_out_chan)
			time.Sleep(10 * time.Second)
		}

	} else if *optionPtr == "getbalance" {
		log.Println("getbalance")
		//protocol.Server_address

		Getbalance(msg_in_chan, msg_out_chan)

	} else if *optionPtr == "faucet" {
		log.Println("faucet")
		//get coins
		//GetFaucet(rw)
		GetFaucet(msg_in_chan, msg_out_chan)

	} else if *optionPtr == "txpool" {
		err = Gettxpool(msg_in_chan, msg_out_chan)
		return
	} else if *optionPtr == "pushtx" {
		err = PushTx(msg_in_chan, msg_out_chan)
		return
	} else if *optionPtr == "randomtx" {
		err = MakeRandomTx(msg_in_chan, msg_out_chan)
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
