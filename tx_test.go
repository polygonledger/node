package main

import (
	"encoding/json"
	"testing"

	"github.com/polygonledger/node/block"
)

//basic block functions
func TestBasicAssign(t *testing.T) {
	var tx block.Tx
	tx = block.Tx{Nonce: 1}
	if tx.Nonce != 1 {
		t.Error("fail assign nonce")
	}
}

func TestTxJson(t *testing.T) {
	var tx block.Tx
	tx = block.Tx{Nonce: 1}
	txJson, _ := json.Marshal(tx)
	if txJson[0] != '{' {
		t.Error("start json")
	}
	i := len(txJson) - 1
	if txJson[i] != '}' {
		t.Error("end json")
	}

	var newtx block.Tx
	if err := json.Unmarshal(txJson, &newtx); err != nil {
		panic(err)
	}
	if newtx.Nonce != tx.Nonce {
		t.Error("json marshal failed")
	}
	if newtx.Sender != tx.Sender {
		t.Error("json marshal failed")
	}
}
