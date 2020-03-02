package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/polygonledger/node/ntcl"
)

//simple node that runs standalone without peers

var srv Server

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
		ntchan := ntcl.ConnNtchan(conn)

		go t.handleConnection(ntchan)
	}
	return
}

func echohandler(ins string) string {
	resp := "Echo:" + ins
	return resp
}

func (t *TCPServer) handleConnection(ntchan ntcl.Ntchan) {
	//tr := 100 * time.Millisecond
	//defer ntchan.Conn.Close()
	log.Println("handleConnection")

	go ntcl.ReadLoop(ntchan)
	go ntcl.ReadProcessor(ntchan)
}

//deal with the logic of each connection
func (t *TCPServer) handleConnection1(conn net.Conn) {
	tr := 100 * time.Millisecond
	defer conn.Close()
	log.Println("handleConnection")

	for {

		log.Println("read with delim ", ntcl.DELIM)
		req, err := ntcl.NtwkRead(conn, ntcl.DELIM)

		if err != nil {
			log.Println(err)
		}

		if len(req) > 0 {
			log.Println("=> ", req, len(req))
			req = strings.Trim(req, string(ntcl.DELIM))
			resp := echohandler(req)

			log.Println("resp => ", resp)
			ntcl.NtwkWrite(conn, resp)

		} else {
			//empty read next read slower
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
