package main

// import (
// 	"log"
// 	"net"
// 	"strconv"
// 	"testing"

// 	"github.com/polygonledger/node/ntcl"
// )

// const test_node_port = 8888

// func init() {
// 	// Start the new server

// 	srv, err := NewServer(":" + strconv.Itoa(test_node_port))
// 	if err != nil {
// 		log.Println("error starting TCP server")
// 		return
// 	}

// 	// Run the server in Goroutine to stop tests from blocking
// 	// test execution
// 	go func() {
// 		srv.Run()
// 	}()
// }

// func TestServer_Run(t *testing.T) {
// 	// Simply check that the server is up and can
// 	// accept connections
// 	conn, err := net.Dial("tcp", ":"+strconv.Itoa(test_node_port))
// 	if err != nil {
// 		t.Error("could not connect to server: ", err)
// 	} else {
// 		//t.Error("...")
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
