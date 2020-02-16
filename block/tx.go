package block

//chain "github.com/polygonledger/node/chain"

//potential TransactionTypes
// SEND_CASH
// REGISTER_DELEGATE
// VOTE_DELEGATE
// REGISTER_USERNAME

const (
	SEND_CASH = iota
	REGISTER_DELEGATE
)

type Tx struct {
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

//type TxSig struct {
