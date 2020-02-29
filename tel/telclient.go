package main

//telnet like client shows how to interact with a node

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
	ntwk "github.com/polygonledger/node/ntwk"
)

func main() {

	fmt.Print("Enter command: ")
	reader := bufio.NewReader(os.Stdin)

	cmd, _ := reader.ReadString('\n')
	cmd = strings.Trim(cmd, string('\n'))

	mainPeerAddress := "localhost:8888" // "polygonnode.com"
	mainPeer := ntwk.CreatePeer(mainPeerAddress, 8888)
	log.Println("client with mainPeer ", mainPeer)

	log.Println("Dial " + mainPeerAddress)
	conn, err := net.Dial("tcp", mainPeerAddress)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
		log.Println("error ", err)
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	log.Println(rw)

	n, err := rw.WriteString("test|")
	log.Println("n ", n)
}
