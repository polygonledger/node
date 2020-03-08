package main

import (
	"log"
	"testing"

	chain "github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/ntcl"
)

func TestBasicCommand(t *testing.T) {

	log.Println("TestBasicCommand")

	mgr := chain.CreateManager()
	mgr.InitAccounts()

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)
	msg := ntcl.ParseMessage(req_msg)
	//ntchan.REQ_in <- msg
	reply_msg := HandlePing(msg)
	if reply_msg != "REP#PONG#|" {
		t.Error(reply_msg)
	}

	ntchan := ntcl.ConnNtchanStub("")
	go RequestHandlerTel(&mgr, ntchan)
	ntchan.REQ_in <- req_msg
	reply := <-ntchan.REP_out

	if reply != "REP#PONG#|" {
		t.Error(reply_msg)
	}

}

func TestBalance(t *testing.T) {

	log.Println("TestBalance")

	mgr := chain.CreateManager()
	mgr.InitAccounts()

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_BALANCE, "abc")
	msg := ntcl.ParseMessage(req_msg)

	reply_msg := HandleBalance(&mgr, msg)
	if reply_msg != "REP#BALANCE#0|" {
		t.Error(reply_msg)
	}

	//TODO with chain setup

	//log.Println(mgr.Accounts)
	ra := mgr.RandomAccount()
	mgr.SetAccount(ra, 10)

	log.Println(ra.AccountKey)
	req_msg = ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_BALANCE, ra.AccountKey)
	msg = ntcl.ParseMessage(req_msg)

	reply_msg = HandleBalance(&mgr, msg)
	if reply_msg != "REP#BALANCE#10|" {
		t.Error(reply_msg)
	}

	log.Println(mgr.Accounts)

	//log.Println(ra)

	// b := block.Block{}
	// tx := tx.Tx{}
	// mgr.ApplyBlock(b)
}

func TestRanaccount(t *testing.T) {
	// REQ#RANACC#|
	// REP#RANACC#{"AccountKey":"Pe2c32a4f7e8b"}|
	// REQ#BALANCE#Pe2c32a4f7e8b/
	// REP#BALANCE#0|
	// REQ#BALANCE#Pe2c32a4f7e8b|
	// REP#BALANCE#20|
}

func TestGenesis(t *testing.T) {

	// genBlock := chain.MakeGenesisBlock()
	// mgr.ApplyBlock(genBlock)
	// //chain.SetAccount()

	// for k, v := range mgr.Accounts {
	// 	fmt.Println(k, v)
	// 	if !mgr.IsTreasury(k) {
	// 		if v != 20 {
	// 			t.Error("...")
	// 		}
	// 	} else {
	// 		if v != 200 {
	// 			t.Error("...")
	// 		}
	// 	}

}
