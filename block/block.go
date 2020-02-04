package block

type Block struct {
	Hash []byte
	// PrevHash []byte
	Height int
	Txs    []Tx
}
