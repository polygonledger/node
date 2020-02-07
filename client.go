package main

import (
	"bufio"
	"log"
	"net"

	"github.com/pkg/errors"

	protocol "github.com/polygonledger/net"
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
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+protocol.Port)
	}

	protocol.SendTx(rw)

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
	log.Println("Client done.")
	return

}
