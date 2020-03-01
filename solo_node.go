package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/polygonledger/node/ntwk"
)

//simple node that runs standalone without peers

var srv Server

//const DELIM = '\n'
const DELIM = '|'

// minimum TCP server
type Server interface {
	Run() error
	Close() error
}

type TCPServer struct {
	addr   string
	server net.Listener
}

// start listening on tcp and handle connection through channels
func (t *TCPServer) Run() (err error) {

	//TODO! nodeport int
	//TODO! hearbeart, check if peers are alive
	//TODO! handshake

	log.Println("listen ", t.addr)
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", t.addr)
	}
	//run forever and don't close
	//defer t.Close()

	for {
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			break
		}
		if conn == nil {
			err = errors.New("could not create connection")
			break
		}
		strRemoteAddr := conn.RemoteAddr().String()
		log.Println("accepted conn ", strRemoteAddr)

		//?
		//setupPeer(strRemoteAddr, node_port, conn)
		//conn.Close()
		//t.handleConnections()
		go t.handleConnection(conn)
	}
	return
}

// handleConnections deals with the logic of
// each connection and their requests
func (t *TCPServer) handleConnection(conn net.Conn) {
	tr := 100 * time.Millisecond
	defer conn.Close()
	log.Println("handleConnection")

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {

		log.Println("read with delim ", DELIM)
		//req, err := rw.ReadString(DELIM)
		req, err := ntwk.NtwkRead(conn, DELIM)

		if err != nil {
			rw.WriteString("failed to read input")
			rw.Flush()
			return
		}

		if len(req) > 0 {
			log.Println("=> ", req, len(req))
			req = strings.Trim(req, string(DELIM))
			resp := "Echo: " + req
			log.Println("resp => ", resp)
			ntwk.NtwkWrite(conn, resp)

		} else {
			tr += 100 * time.Millisecond
		}

		time.Sleep(tr)
		//on empty reads increase time, but max at 800
		if tr > 800*time.Millisecond {
			tr = 800 * time.Millisecond
		}

	}
}

// NewServer creates a new Server using given protocol
// and addr
func NewServer(addr string) (Server, error) {
	return &TCPServer{
		addr: addr,
	}, nil

	return nil, errors.New("Invalid protocol given")
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

const node_port = 8888

func main() {

	srv, err := NewServer(":" + strconv.Itoa(node_port))

	if err != nil {
		log.Println("error starting TCP server")
		return
	}

	srv.Run()

}
