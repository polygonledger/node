package block

import (
	chain "github.com/polygon/chain"
)

type Tx struct {
	Id       [32]byte      `json:"id"` //gets assigned when verified in a block
	Nonce    int           `json:"Nonce"`
	Amount   int           `json:"Amount"`
	Sender   chain.Account //[32]byte
	Receiver chain.Account //[32]byte

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
