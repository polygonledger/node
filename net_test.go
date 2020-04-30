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

func initserver() *TCPNode {
	log.Println("initserver")
	// Start the new server

	testsrv, err := NewNode()
	testsrv.addr = ":" + strconv.Itoa(test_node_port)

	if err != nil {
		log.Println("error starting TCP server")
		return testsrv
	} else {
		log.Println("start ", testsrv)
	}

	// Run the server in Goroutine to stop tests from blocking
	// test execution
	//log.Println("initserver  ", testsrv)

	go testsrv.Run()
	//log.Println("waiting ", newpeerchan)
	go testsrv.HandleConnectTCP()

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
	ntchan := ntcl.ConnNtchan(conn, "client", addr, false)
	//defer conn.Close()
	return ntchan

}

func TestServer_Run(t *testing.T) {
	log.Println("TestServer_Run")

	testsrv := initserver()
	defer testsrv.Close()

	time.Sleep(800 * time.Millisecond)

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
	log.Println("TestServer_Write")

	testsrv := initserver()
	defer testsrv.Close()

	clientNt := testclient()
	go ntcl.ReadLoop(clientNt)
	//go ntcl.ReadProcessor(clientNt)

	time.Sleep(2000 * time.Millisecond)

	peers := testsrv.GetPeers()
	if len(peers) != 1 {
		t.Error("no peers ", testsrv.Peers, len(peers))
	}

	firstpeer := peers[0]

	if !isEmpty(firstpeer.NTchan.Writer_queue, 1*time.Second) {
		t.Error("fail")
	}

	reqs := "hello world"
	n, err := ntcl.NetWrite(firstpeer.NTchan, reqs)

	if err != nil {
		t.Error("could not write to server:", err)
	}

	delimsize := 1
	l := len([]byte(reqs)) + delimsize
	if n != l {
		t.Error("wrong bytes written ", l)
	}

	time.Sleep(100 * time.Millisecond)

	//log.Println(clientNt.SrcName)

	rmsg := <-clientNt.Reader_queue
	if rmsg != reqs {
		t.Error("different message on reader ", rmsg)
	}
	// if isEmpty(clientNt.Reader_queue, 1*time.Second) {
	// 	t.Error("fail")
	//}

}

// func TestServer_Request(t *testing.T) {

// 	addr := ":" + strconv.Itoa(test_node_port)
// 	conn, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		t.Error("could not connect to server: ", err)
// 	}
// 	defer conn.Close()

// 	reqs := "hello world"

// 	n, err := ntcl.NetWrite(conn, reqs)
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
