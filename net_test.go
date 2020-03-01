package main

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"testing"
)

const test_node_port = 8888

func init() {
	// Start the new server

	srv, err := NewServer(":" + strconv.Itoa(test_node_port))
	if err != nil {
		log.Println("error starting TCP server")
		return
	}

	// Run the server in Goroutine to stop tests from blocking
	// test execution
	go func() {
		srv.Run()
	}()
}

func TestServer_Run(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(test_node_port))
	if err != nil {
		t.Error("could not connect to server: ", err)
	} else {
		//t.Error("...")
	}
	defer conn.Close()
}

func TestServer_Request(t *testing.T) {

	addr := ":" + strconv.Itoa(test_node_port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()

	req := []byte("hello world" + string(DELIM))

	if _, err := conn.Write(req); err != nil {
		t.Error("could not write payload to server:", err)
	}

	r := []byte("Request received: hello world")
	out := make([]byte, 1024)
	if _, err := conn.Read(out); err == nil {
		if bytes.Compare(out, r) == 0 {
			t.Error("response did match expected output")
		} else {
			log.Println("response ok")
		}
	} else {
		t.Error("could not read from connection")
	}

}
