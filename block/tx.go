package block

type Tx struct {
	Id       [32]byte `json:"id"` //gets assigned when verified in a block
	Nonce    int      `json:"Nonce"`
	Amount   int      `json:"Amount"`
	Sender   string   //[32]byte
	Receiver string   //[32]byte

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
