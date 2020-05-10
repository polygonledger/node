package main

import (
	"log"
	"strconv"
	"testing"

	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
)

func TestBasicCommand(t *testing.T) {

	log.Println("TestBasicCommand")

	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	node.Mgr = &mgr

	// req_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_PING)
	// //ntchan.REQ_in <- msg
	// reply_msg := HandlePing(msg)

	// ntchan := ntcl.ConnNtchanStub("")
	// go RequestHandlerTel(node, ntchan)
	// ntchan.REQ_in <- req_msg
	// reply := <-ntchan.REP_out

	// //if reply != "REP#PONG#|" {
	// if reply != "out" {
	// 	t.Error("reply_msg ", reply)
	// }

}

func TestPing(t *testing.T) {
	reqstring := ntcl.EncodeMsgMap(ntcl.REQ, "PING")
	req := ntcl.ParseMessageMap(reqstring)
	if req.MessageType != "REQ" || req.Command != "PING" {
		t.Error("req ", req)
	}
	reply := HandlePing(req)
	if reply.MessageType != "REP" || reply.Command != "PONG" {
		t.Error("reply type ", reply)
	}
}

// func TestPingRequest(t *testing.T) {

// 	reqstring := ntcl.EncodeMsgMap(ntcl.REQ, "PING")
// 	req := ntcl.ParseMessageMap(reqstring)

// }

func TestBalance(t *testing.T) {

	log.Println("TestBalance")

	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node.Mgr = &mgr

	req_msg := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_BALANCE, "abc")
	//fmt.Println(req_msg)

	msg := ntcl.ParseMessageMap(req_msg)

	reply_msg := HandleBalance(node, msg)
	target := "{:REP BALANCE :data 0}"
	if reply_msg != target {
		t.Error("reply_msg ", reply_msg, target)
	}

	msg = ntcl.ParseMessageMap(reply_msg)

	if msg.MessageType != ntcl.REP {
		t.Error("balance msg")
	}

	//TODO with chain setup

	ra := mgr.RandomAccount()
	mgr.SetAccount(ra, 10)

	//log.Println(ra.AccountKey)
	req_msg_string := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_BALANCE, ra)
	if req_msg_string != "{:REQ BALANCE :data P2e2bfb58c9db}" {
		t.Error("req string")
	}

	req_msg_balance := ntcl.ParseMessageMapData(req_msg_string)

	reply_msg = HandleBalance(node, req_msg_balance)
	if reply_msg != "{:REP BALANCE :data 10}" {
		t.Error("reply_msg ", reply_msg)
	}

	if mgr.Accounts["P2e2bfb58c9db"] != 10 {
		t.Error("balance")
	}

	//log.Println(ra)

	// b := block.Block{}
	// tx := tx.Tx{}
	// mgr.ApplyBlock(b)
}

func TestFaucetTx(t *testing.T) {

	kp := crypto.PairFromSecret("test")
	pubk := crypto.PubKeyToHex(kp.PubKey)
	addr := crypto.Address(pubk)
	req_msg_string := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_FAUCET, addr)
	req_msg := ntcl.ParseMessageMap(req_msg_string)

	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node.Mgr = &mgr

	reply_msg := HandleFaucet(node, req_msg)
	if reply_msg != "{:REP FAUCET :data ok}" {
		t.Error("reply_msg ", reply_msg)
	}

	// chain.MakeBlock(&mgr)

	// time.Sleep(2000 * time.Millisecond)

	// req_msg = ntcl.EncodeMEncodeMsgxx(ntcl.REQ, ntcl.CMD_BALANCE, addr)
	// msg = ntcl.ParseMessag(req_msg)

	// log.Println(mgr.Accounts)

	// reply_msg_string := HandleBalance(node, msg)
	// log.Println(reply_msg_string)
	// msg = ntcl.ParseMessag(reply_msg_string)
	// // if reply_msg_string != "REP#BALANCE#1|" {
	// // 	t.Error(msg)
	// // }

	// // bal := ntcl.ParseMessageBalance(reply_msg)
	// if msg.MessageType != "REP" || msg.Command != ntcl.CMD_BALANCE {
	// 	t.Error("msg ", msg)
	// }

}

//TODO fix
// func TestTx(t *testing.T) {

// 	node, _ := NewNode()
// 	//defer node.Close()
// 	node.addr = ":" + strconv.Itoa(8888)
// 	mgr := chain.CreateManager()
// 	mgr.InitAccounts()
// 	node.Mgr = &mgr
// 	node.Loglevel = LOGLEVEL_OFF
// 	kp := crypto.PairFromSecret("test")
// 	pubk := crypto.PubKeyToHex(kp.PubKey)
// 	addr := crypto.Address(pubk)
// 	req_msg := ntcl.EncodeMsgxx(ntcl.REQ, ntcl.CMD_FAUCET, addr)
// 	msg := ntcl.ParseMessag(req_msg)

// 	reply_msg := HandleFaucet(node, msg)
// 	if reply_msg != "REP#FAUCET#ok|" {
// 		t.Error("reply_msg ", reply_msg)
// 	}
// 	chain.MakeBlock(&mgr)
// 	time.Sleep(100 * time.Millisecond)
// 	req_msg = ntcl.EncodeMsgxx(ntcl.REQ, ntcl.CMD_BALANCE, addr)
// 	msg = ntcl.ParseMessag(req_msg)
// 	reply_msg = HandleBalance(node, msg)
// 	msg = ntcl.ParseMessa(reply_msg)

// 	if msg.MessageType != "REP" || msg.Command != ntcl.CMD_BALANCE {
// 		t.Error("msg ", msg)
// 	}

// 	sender := block.AccountFromString(addr)

// 	kp2 := crypto.PairFromSecret("test2")
// 	pubk2 := crypto.PubKeyToHex(kp2.PubKey)
// 	addr2 := crypto.Address(pubk2)
// 	recv := block.AccountFromString(addr2)

// 	amount := 1
// 	tx := block.Tx{Nonce: 1, Amount: amount, Sender: sender, Receiver: recv}
// 	signature := crypto.SignTx(tx, kp)
// 	sighex := hex.EncodeToString(signature.Serialize())
// 	tx.Signature = sighex
// 	tx.SenderPubkey = crypto.PubKeyToHex(kp.PubKey)

// 	verified := crypto.VerifyTxSig(tx)

// 	if !verified {
// 		t.Error("not verified")
// 	}

// 	valid := chain.TxValid(&mgr, tx)

// 	if !valid {
// 		t.Error("not valid")
// 	}

// 	txJson, _ := json.Marshal(tx)
// 	req_msg = ntcl.EncodeMessageTx(txJson)
// 	msg = ntcl.ParseMessa(req_msg)

// 	reply_msg = HandleTx(node, msg)
// 	// //TODO!
// 	if reply_msg != "REP#TX#ok|" {
// 		t.Error("reply_msg ", reply_msg)
// 	}

// 	chain.MakeBlock(&mgr)
// 	time.Sleep(100 * time.Millisecond)

// 	req_msg = ntcl.EncodeMsgxx(ntcl.REQ, ntcl.CMD_BALANCE, addr2)
// 	msg = ntcl.ParseMessae(req_msg)
// 	reply_msg = HandleBalance(node, msg)
// 	msg = ntcl.ParseMessa(reply_msg)
// 	//if reply_msg != "REP#BALANCE#5|" {
// 	//bal := ntcl.ParseMeseBalance(reply_msg)

// 	if msg.MessageType != "REP" || msg.Command != ntcl.CMD_BALANCE {
// 		t.Error("reply_msg ", reply_msg)
// 	}

// }

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
