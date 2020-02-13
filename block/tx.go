package block

import "github.com/btcd/btcec"

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
	Id        [32]byte `json:"id"` //gets assigned when verified in a block
	Nonce     int      `json:"Nonce"`
	Amount    int      `json:"Amount"`
	Sender    Account  //[32]byte
	Receiver  Account  //[32]byte
	Signature btcec.Signature
	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
