package net

import (
	"bufio"
	"log"
	"net"

	"github.com/pkg/errors"
)

type Peer struct {
	Address string `json:"Address"`
}

//peer functions
//onReceiveBlock
//validateBlockSlot
//generateBlock
//loadBlocksFromPeer
//loadBlocksOffset
//getCommonBlock //Performs chain comparison with remote peer

// Open connects to a TCP Address
// It returns a TCP connection with a timeout wrapped into a buffered ReadWriter.
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func OpenOut(ip string) *bufio.ReadWriter {
	rw, _ := Open(ip + Port)
	return rw
}
