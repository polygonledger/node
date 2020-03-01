package main

//telnet like client shows how to interact with a node
//can send an

import (
	"bufio"
	"log"
	"net"

	ntwk "github.com/polygonledger/node/ntwk"
)

//const DELIM = "|"
const DELIM byte = '|'

func main() {

	//fmt.Print("Enter command: ")
	//reader := bufio.NewReader(os.Stdin)
	//cmd, _ := reader.ReadString('\n')
	//cmd = strings.Trim(cmd, string('\n'))

	mainPeerAddress := "localhost:8888" // "polygonnode.com"
	mainPeer := ntwk.CreatePeer(mainPeerAddress, 8888)
	log.Println("client with mainPeer ", mainPeer)

	log.Println("Dial " + mainPeerAddress)
	conn, err := net.Dial("tcp", mainPeerAddress)
	if err != nil {
		log.Println("error ", err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	//n, err := rw.WriteString("REQ#PING#" + string(DELIM))
	//s1 := "REQ#PING#" + string(DELIM)
	s1 := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, "")

	n, err := rw.WriteString(s1)
	if err != nil {
		log.Println("err ", err)
	}
	rw.WriteString(string('|'))
	log.Println("bytes written: ", n)
	log.Println(s1)

	s, _ := rw.ReadString(DELIM)
	log.Println(s)
}
