package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strconv"

	"github.com/pkg/errors"

	block "github.com/polygon/block"
	protocol "github.com/polygon/net"
)

var server_address string = "127.0.0.1"

/*
Outgoing connections. A `net.Conn` satisfies the io.Reader and `io.Writer` interfaces
we can treat a TCP connection just like any other `Reader` or `Writer`.
*/

// Open connects to a TCP Address.
// It returns a TCP connection armed with a timeout and wrapped into a
// buffered ReadWriter.
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

/*
The client function connects to the server and sends GOB requests.
*/
func sendTx(rw *bufio.ReadWriter) error {

	// Send a GOB request.
	// Create an encoder that directly transmits to `rw`.
	// Send the request name.
	// Send the GOB.

	testTx := block.Tx{Nonce: 42}

	log.Println("Send a struct as GOB:")
	log.Printf("testTx: \n%#v\n", testTx)

	enc := gob.NewEncoder(rw)
	//Command

	n, err := rw.WriteString(protocol.CMD_GOB)
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(testTx)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", testTx)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

func client(ip string) error {

	// Open a connection to the server.
	rw, err := Open(ip + protocol.Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+protocol.Port)
	}

	sendTx(rw)

	return nil
}

/*
start as a client and connects to the host
*/
func main() {

	err := client(server_address)
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	log.Println("Client done.")
	return

}
