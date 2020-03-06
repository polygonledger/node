package ntcl

type Peer struct {
	Address  string `json:"Address"`
	NodePort int
	NTchan   Ntchan

	//----------------
	//Name     string //can set name
	//Domain string
}

func CreatePeer(ipAddress string, nodeport int) Peer {
	//addr := ip
	//NodePort: NodePort,
	p := Peer{Address: ipAddress, NodePort: nodeport}
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
