package block

//potential TransactionTypes
// VOTE_DELEGATE
// REGISTER_NAME

const (
	CASH_TRANSFER = "CASH_TRANSFER"
	//DELEGATE_REGISTER = "DELEGATE_REGISTER"
	CONCESSION_REG = "CONCESSION_REG"
)

type TxSigmap struct {
	SenderPubkey string `json:"senderPubkey"`
	Signature    string `json:"signature"`
}

type Tx struct {
	TxType   string `json:"txType"`
	Amount   int    `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	//TxSig   TxSig  `json:"txsig"`
	Nonce int `json:"nonce"`
	//TODO replace with txsig
	SenderPubkey string `json:"senderPubkey"`
	Signature    string `json:"signature"`
	//Id           [32]byte `edn:"id"`           //gets assigned when verified in a block

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}

// type TxE struct {
// 	TxType string `edn:"TxType"`
// 	//Transfer   string `edn:"Transfer"`
// 	//Sigmap
// }

// type SimpleTx struct {
// 	Amount   int    `edn:"amount"`
// 	Sender   string `edn:"sender"`   //[32]byte
// 	Receiver string `end:"receiver"` //[32]byte
// 	//Nonce        int    `edn:"Nonce"`
// }

type TxSigmapEdn struct {
	SenderPubkey string `edn:"senderPubkey"`
	Signature    string `edn:"signature"`
}

// type TxExpr struct {
// 	TxType   string   `edn:"TxType"`
// 	Transfer SimpleTx `edn:"TxTransfer"`
// 	Sigmap   TxSigmap `edn:"Sigmap"`
// }

type TxEdn struct {
	TxType   string `edn:"TxType"`
	Amount   int    `edn:"Amount"`
	Sender   string `edn:"Sender"`   //[32]byte
	Receiver string `end:"Receiver"` //[32]byte
	//TODO delete
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
