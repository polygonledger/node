package netio

type Peer struct {
	Address  string `json:"Address"`
	Name     string
	NodePort int
	NTchan   Ntchan

	//----------------
	//Name     string //can set name
	//Domain string
}

func CreatePeer(name string, ipAddress string, nodeport int, ntchan Ntchan) Peer {
	//addr := ip
	//NodePort: NodePort,
	p := Peer{Name: name, Address: ipAddress, NodePort: nodeport, NTchan: ntchan}
	return p
}

//peer functions
//get peers
//onReceiveBlock
//validateBlockSlot
//generateBlock
//loadBlocksFromPeer
//loadBlocksOffset
//getCommonBlock //Performs chain comparison with remote peer
