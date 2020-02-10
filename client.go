package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/polygonledger/node/block"
	cryptoutil "github.com/polygonledger/node/crypto"
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

func client(ip string) error {

	// Open a connection to the server.
	rw, err := Open(ip + protocol.Port)
	log.Println(rw)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+protocol.Port)
	}

	//get random account

	//DELIM_HEAD := "#" //TODO from protocol
	//msg := protocol.REQ + string(protocol.DELIM_HEAD) + protocol.CMD_RANDOM_ACCOUNT + string(DELIM_HEAD) + "emptydata" + string(protocol.DELIM)
	msg := ""
	n, err := rw.WriteString(msg)
	if err != nil {
		return errors.Wrap(err, "Could not write JSON data ("+strconv.Itoa(n)+" bytes written)")
	}

	err = rw.Flush()

	msg2, err2 := rw.ReadString(protocol.DELIM)
	msg2 = strings.Trim(msg2, string(protocol.DELIM))
	var a block.Account
	dataBytes := []byte(msg2)
	if err := json.Unmarshal(dataBytes, &a); err != nil {
		panic(err)
	}
	log.Print(" account key ", a.AccountKey)
	if err != nil {
		log.Println("Failed ", err2)
		//log.Println(err.)
		//break
	}

	//use this random account to send coins from

	//send Tx
	testTx := protocol.RandomTx(a)
	txJson, _ := json.Marshal(testTx)
	log.Println("txJson ", txJson)

	//msg3 := protocol.REQ + protocol.DELIM_HEAD + "TX" + protocol.DELIM_HEAD + string(txJson) + string(protocol.DELIM)
	//msg3 := protocol.REQ + string(protocol.DELIM_HEAD) + "TX" + string(protocol.DELIM_HEAD) + string(txJSon)
	var msg3 string
	msg3 = string(protocol.DELIM)
	n3, err3 := rw.WriteString(msg3)
	if err3 != nil {
		return errors.Wrap(err, "Could not write JSON data ("+strconv.Itoa(n3)+" bytes written)")
	}

	err = rw.Flush()

	msg4, err2 := rw.ReadString(protocol.DELIM)
	log.Print("reply msg ", msg4)
	if err != nil {
		log.Println("Failed ", err2)
		//log.Println(err.)
		//break
	}

	return nil
}

/*
start client and connect to the host
*/
func main() {

	// err := client(protocol.Server_address)
	// if err != nil {
	// 	log.Println("Error:", errors.WithStack(err))
	// }
	// log.Println("Client done")
	// return
	kp := cryptoutil.PairFromSecret("test")
	log.Println("keypair ", kp)

}
