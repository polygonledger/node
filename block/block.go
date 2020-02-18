package block

import (
	"time"

	"github.com/btcsuite/btcd/btcec"
)

type Block struct {
	Hash            [32]byte
	Prev_Block_Hash [32]byte
	Height          int
	Txs             []Tx
	Timestamp       time.Time
	Signature       btcec.Signature
}
