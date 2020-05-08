package main

import (
	"testing"

	"github.com/polygonledger/edn"
	"github.com/polygonledger/node/block"
)

func TestEDN(t *testing.T) {

	data := `{:TxType "test",
			  :Sender "abc",
			  :Receiver "xyz",
		      :amount 42,
			  :nonce 1}`

	var tx block.Tx
	err := edn.Unmarshal([]byte(data), &tx)

	if err != nil {
		t.Error("err ", err)
	}

	if tx.Nonce != 1 || string(tx.TxType) != "test" || tx.Sender != "abc" || tx.Receiver != "xyz" {
		t.Error("tx edn parse")
	}

}
