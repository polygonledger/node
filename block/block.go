package block

import "time"

type Block struct {
	Hash            [32]byte
	Prev_Block_Hash [32]byte
	Height          int
	Txs             []Tx
	Timestamp       time.Time
}
