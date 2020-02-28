package ntwk

import (
	"log"
	"time"
)

type Peer struct {
	Address  string `json:"Address"`
	NodePort int

	//TODO! remove
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

func MakePingOld(peer Peer) bool {
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
