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

//var srv Server

const node_port = 8888

type TCPServer struct {
	Name          string
	addr          string
	server        net.Listener
	accepting     bool
	ConnectedChan chan net.Conn
	//TODO! list of peers
	Peers []ntcl.Peer
}

func (t *TCPServer) GetPeers() []ntcl.Peer {
	if &t.Peers == nil {
		return nil
	}
	return t.Peers
}

// start listening on tcp and handle connection through channels
func (t *TCPServer) Run() (err error) {

	log.Println("listen ", t.addr)
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		//return errors.Wrapf(err, "Unable to listen on port %s\n", t.addr)
	}
	//run forever and don't close
	//defer t.Close()

	for {
		t.accepting = true
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			break
		}
		if conn == nil {
			err = errors.New("could not create connection")
			break
		}

		log.Println("put on out ", conn)
		t.ConnectedChan <- conn

	}
	log.Println("end run")
	return
}

func (t *TCPServer) HandleDisconnect() {

}

//handle new connection
func (t *TCPServer) HandleConnect() {

	//TODO! hearbeart, check if peers are alive
	//TODO! handshake

	for {
		newpeerConn := <-t.ConnectedChan
		strRemoteAddr := newpeerConn.RemoteAddr().String()
		log.Println("accepted conn ", strRemoteAddr, t.accepting)
		log.Println("new peer ", newpeerConn)
		// log.Println("> ", t.Peers)
		// log.Println("# peers ", len(t.Peers))

		ntchan := ntcl.ConnNtchan(newpeerConn, "server", strRemoteAddr)

		p := ntcl.Peer{Address: strRemoteAddr, NodePort: node_port, NTchan: ntchan}
		t.Peers = append(t.Peers, p)

		go t.handleConnection(ntchan)
		//go ChannelPeerNetwork(conn, peer)
		//setupPeer(strRemoteAddr, node_port, conn)

		//conn.Close()

	}
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
//simple readwriter
func (t *TCPServer) handleConnectionReadWriter(ntchan ntcl.Ntchan) {
	tr := 100 * time.Millisecond
	defer ntchan.Conn.Close()
	log.Println("handleConnection")

	for {

		log.Println("read with delim ", ntcl.DELIM)
		req, err := ntcl.NtwkRead(ntchan, ntcl.DELIM)

		if err != nil {
			log.Println(err)
		}

		if len(req) > 0 {
			log.Println("=> ", req, len(req))
			req = strings.Trim(req, string(ntcl.DELIM))
			resp := echohandler(req)

			log.Println("resp => ", resp)
			ntcl.NtwkWrite(ntchan, resp)

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
func NewServer(addr string) (*TCPServer, error) {
	return &TCPServer{
		addr:          addr,
		accepting:     false,
		ConnectedChan: make(chan net.Conn),
		//Peers: make([]ntcl.Peer)
	}, nil

}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

func main() {

	srv, err := NewServer(":" + strconv.Itoa(node_port))

	if err != nil {
		log.Println("error creating TCP server")
		return
	}

	// if err2 != nil {
	// 	log.Println("error starting TCP server ", err2)
	// 	return
	// }

	go srv.HandleConnect()

	srv.Run()
}
