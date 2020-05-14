package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
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

func TestSignTxBasic(t *testing.T) {

	keypair := crypto.PairFromSecret("test")
	pub := crypto.PubKeyToHex(keypair.PubKey)
	//account := block.Account{AccountKey: crypto.Address(pub)}

	randNonce := 0
	amount := 10

	genkeypair := chain.GenesisKeys()
	addr := crypto.Address(crypto.PubKeyToHex(genkeypair.PubKey))
	//Genesis_Account := block.AccountFromString(addr)

	//{"Nonce":0,"Amount":0,"Sender":{"AccountKey":"Pa033f6528cc1"},"Receiver":{"AccountKey":"Pa033f6528cc1"},"SenderPubkey":"","Signature":"","id":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}
	r := crypto.Address(pub)
	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: addr, Receiver: r, SenderPubkey: "", Signature: ""}

	tx = crypto.SignTxAdd(tx, keypair)

	//log.Println(tx)

	verified := crypto.VerifyTxSig(tx)

	if !verified {
		t.Error("verify tx fail")
	}
}

func TestTxFile(t *testing.T) {
	//write tx.json

	keypair := crypto.PairFromSecret("test")

	pubk := crypto.PubKeyToHex(keypair.PubKey)
	addr := crypto.Address(pubk)

	if addr != "Pa033f6528cc1" {
		t.Error("address wrong ", addr)
	}

	keypair_recv := crypto.PairFromSecret("receive")
	addr_recv := crypto.Address(crypto.PubKeyToHex(keypair_recv.PubKey))

	tx := block.Tx{Nonce: 1, Amount: 10, Sender: addr, Receiver: addr_recv}

	signature := crypto.SignTx(tx, keypair)
	sighex := hex.EncodeToString(signature.Serialize())

	tx.Signature = sighex
	tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

	if !(tx.Amount == 10) {
		t.Error("amount wrong ")
	}

	txJson, _ := json.Marshal(tx)

	ioutil.WriteFile("tx_test.json", []byte(txJson), 0644)

	dat, _ := ioutil.ReadFile("tx_test.json")

	os.Remove("tx_test.json")

	var rtx block.Tx

	if err := json.Unmarshal(dat, &rtx); err != nil {
		panic(err)
	}

	//log.Println(rtx.SenderPubkey)

	if !(rtx.Amount == 10) {
		t.Error("amount wrong ")
	}

	verified := crypto.VerifyTxSig(tx)

	if !verified {
		t.Error("verify tx fail")
	}

}
