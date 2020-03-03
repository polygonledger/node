package main

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/polygonledger/node/ntcl"
)

const test_node_port = 8888

func initserver() *TCPServer {
	// Start the new server

	testsrv, err := NewServer(":" + strconv.Itoa(test_node_port))

	if err != nil {
		log.Println("error starting TCP server")
		return testsrv
	} else {
		log.Println("start ", testsrv)
	}

	// Run the server in Goroutine to stop tests from blocking
	// test execution
	log.Println("initserver >>> ", testsrv)

	go testsrv.Run()
	//log.Println("waiting ", newpeerchan)
	go testsrv.HandleConnect()

	return testsrv
}

func testclient() ntcl.Ntchan {
	time.Sleep(200 * time.Millisecond)
	addr := ":" + strconv.Itoa(test_node_port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//t.Error("could not connect to server: ", err)
	}
	//t.Error("...")
	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr)
	//defer conn.Close()
	return ntchan

}

func TestServer_Run(t *testing.T) {

	testsrv := initserver()

	time.Sleep(1100 * time.Millisecond)

	// Simply check that the server is up and can accept connections
	go testclient()
	// for ok := true; ok; testsrv.accepting = false {
	// 	log.Println(testsrv.accepting)
	// 	time.Sleep(100 * time.Millisecond)
	// }

	time.Sleep(900 * time.Millisecond)
	//log.Println("TestServer_Run > ", testsrv, testsrv.Peers)

	if !testsrv.accepting {
		t.Error("not accepting")
	}

	peers := testsrv.GetPeers()
	if len(peers) != 1 {
		t.Error("no peers ", testsrv.Peers, len(peers))
	}

}

func TestServer_Write(t *testing.T) {

	testsrv := initserver()

	clientNt := testclient()
	go ntcl.ReadLoop(clientNt)
	//go ntcl.ReadProcessor(clientNt)

	time.Sleep(1000 * time.Millisecond)

	peers := testsrv.GetPeers()
	if len(peers) != 1 {
		t.Error("no peers ", testsrv.Peers, len(peers))
	}

	firstpeer := peers[0]

	if !isEmpty(firstpeer.NTchan.Writer_queue, 1*time.Second) {
		t.Error("fail")
	}

	reqs := "hello world"
	n, err := ntcl.NtwkWrite(firstpeer.NTchan, reqs)

	if err != nil {
		t.Error("could not write to server:", err)
	}

	delimsize := 2
	l := len([]byte(reqs)) + delimsize
	if n != l {
		t.Error("wrong bytes written ", l)
	}

	time.Sleep(100 * time.Millisecond)

	log.Println(clientNt.SrcName)

	if isEmpty(clientNt.Reader_queue, 1*time.Second) {
		t.Error("fail")
	}

	// msg, _ := ntcl.MsgRead(clientNt)
	// log.Println("msg ", msg)

	//n, err := ntcl.MsgRead(firstpeer.NTchan, reqs)

	//TODO! need to start reader on firstpeer
	// if isEmpty(firstpeer.NTchan.Reader_queue, 1*time.Second) {
	// 	t.Error("Reader_queue empty")
	// }

}

// func TestServer_Write(t *testing.T) {
// 	testsrv := initserver()

// 	if len(testsrv.GetPeers()) > 0 {
// 		x := <-testsrv.Peers[0].NTchan.Writer_queue
// 		if x == "" {
// 			t.Error("peer nil")
// 		}
// 	}
// 	defer conn.Close()
// }

// func TestServer_Request(t *testing.T) {

// 	addr := ":" + strconv.Itoa(test_node_port)
// 	conn, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Error("could not connect to server: ", err)
// 	}
// 	defer conn.Close()

// 	reqs := "hello world"

// 	n, err := ntcl.NtwkWrite(conn, reqs)
// 	if err != nil {
// 		t.Error("could not write payload to server:", err)
// 	} else {
// 		log.Println("bytes written ", n)
// 	}

// 	read_msg, err := ntcl.MsgRead(conn)

// 	expected_result := ntcl.EncodeMsg("Echo:" + reqs)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	if read_msg != expected_result {
// 		t.Error("response did match expected output ", read_msg, expected_result)
// 	}

// }
