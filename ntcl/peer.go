package ntcl

type Peer struct {
	Address  string `json:"Address"`
	NodePort int
	NTchan   Ntchan

	//----------------
	//Name     string //can set name
	//Domain string
}
