package ntwk

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Peer struct {
	Address  string `json:"Address"`
	NodePort int

	Req_chan     chan Message
	Rep_chan     chan Message
	Out_req_chan chan Message
	Out_rep_chan chan Message
	//wip
	Pub_chan chan Message
	//----------------
	//Name     string //can set name
	//Domain string
}

//peer functions
//get peers
//onReceiveBlock
//validateBlockSlot
//generateBlock
//loadBlocksFromPeer
//loadBlocksOffset
//getCommonBlock //Performs chain comparison with remote peer

func CreatePeer(ipAddress string, nodeport int) Peer {
	//addr := ip
	//NodePort: NodePort,
	p := Peer{Address: ipAddress, NodePort: nodeport, Req_chan: make(chan Message), Rep_chan: make(chan Message), Out_req_chan: make(chan Message), Out_rep_chan: make(chan Message), Pub_chan: make(chan Message)}
	return p
}

func MakePing(peer Peer) bool {
	emptydata := ""
	req_msg := EncodeMsgString(REQ, CMD_PING, emptydata)
	resp := RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	//log.Println("resp ", resp)
	if string(resp.Command) == "PONG" {
		log.Println("ping success")
		return true
	} else {
		log.Println("ping failed ", string(resp.Data))
		return false
	}
}

func HearbeatOnce(peer Peer) bool {
	emptydata := ""
	req_msg := EncodeMsgString(REQ, CMD_PING, emptydata)
	resp := RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	//log.Println("resp ", resp)
	if string(resp.Command) == "PONG" {
		log.Println("ping success")
		return true
	} else {
		log.Println("ping failed ", string(resp.Data))
		return false
	}
}

func Hearbeat(peer Peer) {

	hTime := 1000 * time.Millisecond

	for tt := range time.Tick(hTime) {
		log.Println("heartbeat ", tt)
		HearbeatOnce(peer)
	}

}

func MakeHandshake(peer Peer) bool {
	emptydata := ""
	req_msg := EncodeMsgString(REQ, CMD_HANDSHAKE_HELLO, emptydata)
	resp := RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	//log.Println("resp ", resp)
	if string(resp.Command) == CMD_HANDSHAKE_STABLE {
		log.Println("handshake success")
		return true
	} else {
		log.Println("handshake failed ", string(resp.Data))
		return false
	}
}

//--------network--------

func OpenConn(addr string) net.Conn {
	// Dial the remote process
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	return conn
}

// connects to a TCP Address
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
		log.Println("error ", err)
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func OpenOut(ip string, Port int) (*bufio.ReadWriter, error) {
	addr := ip + ":" + strconv.Itoa(Port)
	log.Println("> open out address ", addr)
	rw, err := Open(addr)
	return rw, err
}
