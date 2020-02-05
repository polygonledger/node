package block

type Tx struct {
	Id     [32]byte `json:"id"` //gets assigned when verified in a block
	Nonce  int      `json:"Nonce"`
	Amount int      `json:"Amount"`
	//sender
	//receiver

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
