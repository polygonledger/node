package main

//telnet like client shows how to interact with a node
//can send an

import (
	"log"
	"net"

	"github.com/polygonledger/node/ntcl"
)

func main() {

	//fmt.Print("Enter command: ")
	//reader := bufio.NewReader(os.Stdin)
	//cmd, _ := reader.ReadString('\n')
	//cmd = strings.Trim(cmd, string('\n'))

	mainPeerAddress := "localhost:8888"
	mainPeer := ntcl.CreatePeer(mainPeerAddress, 8888)
	log.Println("client with mainPeer ", mainPeer)

	log.Println("Dial " + mainPeerAddress)
	conn, err := net.Dial("tcp", mainPeerAddress)
	if err != nil {
		log.Println("error ", err)
	}

	content := "telnet"
	num, err := ntcl.NetWrite(conn, content)

	log.Println(num)

	s, _ := ntcl.NetRead(conn, "}")
	log.Println(s)
}
