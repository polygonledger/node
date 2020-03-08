package main

import (
	"log"
	"math/rand"
	"testing"

	"github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
)

//basic chain functions
func TestChainsetup(t *testing.T) {

	mgr := chain.CreateManager()
	mgr.InitAccounts()
	log.Println(mgr.Accounts)
	ra := mgr.RandomAccount()
	log.Println(ra)

	Genesis_Account := block.AccountFromString(chain.Treasury_Address)
	randNonce := rand.Intn(100)
	r := crypto.RandomPublicKey()
	address_r := crypto.Address(r)
	r_account := block.AccountFromString(address_r)
	amount := 10

	someTx := block.Tx{Nonce: randNonce, Sender: Genesis_Account, Receiver: r_account, Amount: amount}
	log.Println(someTx)

	b := block.Block{}
	b.Txs = []block.Tx{}
	b.Txs = append(b.Txs, someTx)
	mgr.ApplyBlock(b)

	if mgr.Accounts[r_account] != 10 {
		t.Error("wrong amount")
	}

	//chain.InitAccounts()
}
