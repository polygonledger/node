package block

//chain "github.com/polygonledger/chain"

type Tx struct {
	Id       [32]byte `json:"id"` //gets assigned when verified in a block
	Nonce    int      `json:"Nonce"`
	Amount   int      `json:"Amount"`
	Sender   Account  //[32]byte
	Receiver Account  //[32]byte

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
