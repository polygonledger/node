package block

//potential TransactionTypes
// VOTE_DELEGATE
// REGISTER_NAME

const (
	CASH_TRANSFER = "CASH_TRANSFER"
	//DELEGATE_REGISTER = "DELEGATE_REGISTER"
	CONCESSION_REG = "CONCESSION_REG"
)

type Tx struct {
	TxType       string `edn:"TxType"`
	Amount       int    `edn:"Amount"`
	Sender       string `edn:"Sender"`       //[32]byte
	Receiver     string `end:"Receiver"`     //[32]byte
	SenderPubkey string `edn:"SenderPubkey"` //hex string
	Signature    string `edn:"Signature"`    //hex string
	Nonce        int    `edn:"Nonce"`
	//Id           [32]byte `edn:"id"`           //gets assigned when verified in a block

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}

type TxSig struct {
	SenderPubkey string `edn:"SenderPubkey"` //hex string
	Signature    string `edn:"Signature"`    //hex string
}
