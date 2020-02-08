package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/pkg/errors"

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

func client(ip string) error {

	// Open a connection to the server.
	rw, err := Open(ip + protocol.Port)
	log.Println(rw)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+protocol.Port)
	}

	//get random account address
	//protocol.RequestAccount(rw)
	//protocol.RequestTest(rw)

	//protocol.SendTx(rw)
	testTx := protocol.RandomTx()
	txJson, _ := json.Marshal(testTx)
	log.Println("txJson ", txJson)

	msg := "REQ#TX#" + string(txJson) + string(protocol.DELIM)
	n, err := rw.WriteString(msg)
	if err != nil {
		return errors.Wrap(err, "Could not write JSON data ("+strconv.Itoa(n)+" bytes written)")
	}

	err = rw.Flush()

	msg2, err2 := rw.ReadString(protocol.DELIM)
	log.Print("reply msg ", msg2)
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

	err := client(protocol.Server_address)
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	log.Println("Client done")
	return

}
