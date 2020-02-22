package block

//potential TransactionTypes
// VOTE_DELEGATE
// REGISTER_NAME

const (
	TX_SEND_CASH      = "SEND_CASH" //iota
	REGISTER_DELEGATE = "REGISTER_DELEGATE"
)

type Tx struct {
	TxType       string   `json:"TxType"`
	Nonce        int      `json:"Nonce"`
	Amount       int      `json:"Amount"`
	Sender       Account  `json:"Sender"`       //[32]byte
	Receiver     Account  `json:"Receiver"`     //[32]byte
	SenderPubkey string   `json:"SenderPubkey"` //hex string
	Signature    string   `json:"Signature"`    //hex string
	Id           [32]byte `json:"id"`           //gets assigned when verified in a block

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
