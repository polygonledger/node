package main

//telnet like client shows how to interact with a node
//can send an

import (
	"log"
	"net"

	ntwk "github.com/polygonledger/node/ntwk"
)

const DELIM = '|'

func main() {

	//fmt.Print("Enter command: ")
	//reader := bufio.NewReader(os.Stdin)
	//cmd, _ := reader.ReadString('\n')
	//cmd = strings.Trim(cmd, string('\n'))

	mainPeerAddress := "localhost:8888"
	mainPeer := ntwk.CreatePeer(mainPeerAddress, 8888)
	log.Println("client with mainPeer ", mainPeer)

	log.Println("Dial " + mainPeerAddress)
	conn, err := net.Dial("tcp", mainPeerAddress)
	if err != nil {
		log.Println("error ", err)
	}

	content := "telnet"
	num, err := ntwk.NtwkWrite(conn, content)

	log.Println(num)

	s, _ := ntwk.NtwkRead(conn, DELIM)
	log.Println(s)
}
