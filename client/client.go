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

// func getRandomAccount() block.Account {
// 	//TODO
// }

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

func Getbalance(rw *bufio.ReadWriter) error {
	//TODO use messageencoder
	data := "P0614579c42f2"
	//data := "P4e968b02dd42"
	txJson, _ := json.Marshal(block.Account{AccountKey: data})
	msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_BALANCE, string(txJson))
	protocol.WritePipe(rw, msg)

	rcvmsg := protocol.ReadMsg(rw)

	x, _ := strconv.Atoi(rcvmsg)
	log.Println("balance of account ", x)

	return nil
}

func MakePing(rw *bufio.ReadWriter) error {
	//TODO use messageencoder
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_PING, "")
	fmt.Println("send ", req_msg)

	resp_msg := protocol.RequestReply(rw, req_msg)

	log.Println("reply ", resp_msg)

	return nil
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

/*
start client and connect to the host
*/
func main() {

	//prepare to run client

	// domain := "google.com"
	// ips, err1 := net.LookupIP(domain)
	// if err1 != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err1)
	// 	os.Exit(1)
	// }
	// for _, ip := range ips {
	// 	fmt.Printf(domain+". IN A %s\n", ip.String())
	// }

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
	dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)
	fmt.Print("keys ", string(dat))

	optionPtr := flag.String("option", "createkeypair", "the command to be performed")
	flag.Parse()
	fmt.Println("option:", *optionPtr)

	if *optionPtr == "createkeypair" {
		//TODO read secret from cmd
		kp := crypto.PairFromSecret("test")
		log.Println("keypair ", kp)

	} else if *optionPtr == "ping" {
		log.Println("ping")
		//protocol.Server_address
		rw, err := client(configuration.ServerAddress)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		MakePing(rw)

	} else if *optionPtr == "getbalance" {
		log.Println("getbalance")
		//protocol.Server_address
		rw, err := client(configuration.ServerAddress)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		Getbalance(rw)

	} else if *optionPtr == "randomtx" {
		rw, err := client(protocol.Server_address)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		//log.Println("Client done")
		err = MakeRandomTx(rw)
		return
	}

}
