package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	protocol "github.com/polygonledger/node/net"
)

/*
Outgoing connections. A `net.Conn` satisfies the io.Reader and `io.Writer` interfaces
*/

// Open connects to a TCP Address
// It returns a TCP connection with a timeout wrapped into a buffered ReadWriter.
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

//
func MakeRandomTx(rw *bufio.ReadWriter) error {
	//make a random transaction by requesting random account from node
	//get random account

	msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_RANDOM_ACCOUNT, "emptydata")

	protocol.WritePipe(rw, msg)

	//response
	msg2 := protocol.ReadMsg(rw)

	var a block.Account
	dataBytes := []byte(msg2)
	if err := json.Unmarshal(dataBytes, &a); err != nil {
		panic(err)
	}
	log.Print(" account key ", a.AccountKey)

	//use this random account to send coins from

	//send Tx
	testTx := protocol.RandomTx(a)
	txJson, _ := json.Marshal(testTx)
	log.Println("txJson ", txJson)

	req_msg := protocol.EncodeMessageTx(txJson)

	//REQUEST
	//protocol.WritePipe(rw, req_msg)

	//REPLY
	//resp_msg := protocol.ReadMsg(rw)
	resp_msg := protocol.RequestReply(rw, req_msg)
	log.Print("reply msg ", resp_msg)

	return nil
}

func PushTx(rw *bufio.ReadWriter) error {

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
	resp_msg := protocol.RequestReply(rw, req_msg)
	log.Print("reply msg ", resp_msg)

	return nil
}

func Getbalance(rw *bufio.ReadWriter) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter address: ")
	addr, _ := reader.ReadString('\n')
	addr = strings.Trim(addr, string('\n'))

	txJson, _ := json.Marshal(block.Account{AccountKey: addr})
	msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_BALANCE, string(txJson))
	//fmt.Println("msg ", msg)
	protocol.WritePipe(rw, msg)

	rcvmsg := protocol.ReadMsg(rw)

	x, _ := strconv.Atoi(rcvmsg)
	log.Println("balance of account ", x)

	return nil
}

func Gettxpool(rw *bufio.ReadWriter) error {
	msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_GETTXPOOL, "")
	log.Println("> ", msg)
	protocol.WritePipe(rw, msg)

	rcvmsg := protocol.ReadMsg(rw)

	log.Println("rcvmsg ", rcvmsg)

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
	rw, err := Open(ip + protocol.Port)
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

func requestLoop(rw *bufio.ReadWriter, msg_in_chan chan string, msg_out_chan chan string) {
	for {

		//take from channel and request
		request := <-msg_in_chan
		fmt.Println("request ", request)

		resp_msg := protocol.RequestReply(rw, request)
		fmt.Println("resp_msg ", resp_msg)

		msg_out_chan <- resp_msg

	}
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
	go requestLoop(rw, msg_in_chan, msg_out_chan)

	if *optionPtr == "createkeypair" {
		//TODO read secret from cmd
		kp := crypto.PairFromSecret("test")
		log.Println("keypair ", kp)

	} else if *optionPtr == "ping" {
		log.Println("ping")
		//protocol.Server_address
		MakePing(msg_in_chan, msg_out_chan)

	} else if *optionPtr == "getbalance" {
		log.Println("getbalance")
		//protocol.Server_address

		Getbalance(rw)

	} else if *optionPtr == "faucet" {
		log.Println("faucet")
		//get coins
		//GetFaucet(rw)
		GetFaucet(msg_in_chan, msg_out_chan)

	} else if *optionPtr == "txpool" {
		err = Gettxpool(rw)
		return
	} else if *optionPtr == "pushtx" {
		err = PushTx(rw)
		return
	} else if *optionPtr == "randomtx" {
		err = MakeRandomTx(rw)
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
