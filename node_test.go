package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/netio"
)

func TestBasicCommand(t *testing.T) {

	//log.Println("TestBasicCommand")

	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	node.Mgr = &mgr

	// req_msg := netio.EdnConstructMsgMap(netio.REQ, netio.CMD_PING)
	// //ntchan.REQ_in <- msg
	// reply_msg := HandlePing(msg)

	// ntchan := netio.ConnNtchanStub("")
	// go RequestHandlerTel(node, ntchan)
	// ntchan.REQ_in <- req_msg
	// reply := <-ntchan.REP_out

	// //if reply != "REP#PONG#|" {
	// if reply != "out" {
	// 	t.Error("reply_msg ", reply)
	// }

}

func TestPing(t *testing.T) {
	req_msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING}
	//jsonmsgtype, _ := netio.NewJSONMessage(req_msg)

	if req_msg.MessageType != "REQ" || req_msg.Command != "PING" {
		t.Error("req ", req_msg)
	}
	reply := HandlePing(req_msg)
	if reply.MessageType != "REP" || reply.Command != "PONG" {
		t.Error("reply type ", reply)
	}
}

func TestAccountmsg(t *testing.T) {

	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node.Mgr = &mgr

	req_msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_ACCOUNTS}
	//msg := netio.EdnParseMessageMap(req_msg)

	ntchan := netio.ConnNtchanStub("")

	p := netio.Peer{NTchan: ntchan}

	reply_msg := RequestReply(node, &p, req_msg)

	if reply_msg != `{"messagetype":"REP","command":"ACCOUNTS","data":{"P2e2bfb58c9db":400}}` {
		t.Error(reply_msg)
	} else {
		fmt.Println(reply_msg)
	}

}

func TestBalance(t *testing.T) {

	node, _ := NewNode()
	//defer node.Close()
	node.addr = ":" + strconv.Itoa(8888)
	node.Loglevel = LOGLEVEL_OFF
	mgr := chain.CreateManager()
	mgr.InitAccounts()
	node.Mgr = &mgr

	//req_msg := netio.EdnConstructMsgMapData(netio.REQ, netio.CMD_BALANCE, "abc")
	//fmt.Println(req_msg)
	balJson, _ := json.Marshal("abc")
	msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_BALANCE, Data: []byte(balJson)}

	//msg := netio.EdnParseMessageMap(req_msg)

	reply_msg := HandleBalance(node, msg)

	target := `{"messagetype":"REP","command":"BALANCE","data":0}`
	if reply_msg != target {
		t.Error("reply_msg ", reply_msg, target)
	}

	var reply netio.Message
	json.Unmarshal([]byte(reply_msg), &reply)

	//msg = netio.EdnParseMessageMap(reply_msg)

	if reply.MessageType != netio.REP {
		t.Error("balance msg", reply_msg)
	}

	//TODO with chain setup

	ra := mgr.RandomAccount()
	mgr.SetAccount(ra, 10)

	//log.Println(ra.AccountKey)
	//req_msg_string := netio.EdnConstructMsgMapData(netio.REQ, netio.CMD_BALANCE, ra)

	acc := "P2e2bfb58c9db"
	xJson, _ := json.Marshal(acc)
	msg2 := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(xJson)}
	jsonmsg := netio.ToJSONMessage(msg2)

	// msg2 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_BALANCE, Data: []byte("P2e2bfb58c9db")}
	// req_msg_string := netio.ToJSONMessage(msg2)
	if string(jsonmsg) != `{"messagetype":"REP","command":"BALANCE","data":"`+acc+`}` {
		fmt.Println("? ", jsonmsg)
		//t.Error("req string ", req_msg_string, msg2)

	}

	//req_msg_balance := netio.EdnParseMessageMapData(req_msg_string)

	reply_msg = HandleBalance(node, msg2)
	if reply_msg != `{"messagetype":"REP","command":"BALANCE","data":0}` {
		t.Error("reply_msg ", reply_msg)
	} else {
		fmt.Println(mgr.State.Accounts["P2e2bfb58c9db"])
	}

	if mgr.State.Accounts["P2e2bfb58c9db"] != 10 {
		t.Error("balance")
	}

	//log.Println(ra)

	// b := block.Block{}
	// tx := tx.Tx{}
	// mgr.ApplyBlock(b)
}

func TestFaucetTx(t *testing.T) {

	// kp := crypto.PairFromSecret("test")
	// pubk := crypto.PubKeyToHex(kp.PubKey)
	// addr := crypto.Address(pubk)
	// req_msg_string := netio.EdnConstructMsgMapData(netio.REQ, netio.CMD_FAUCET, addr)
	// req_msg := netio.EdnParseMessageMap(req_msg_string)

	// node, _ := NewNode()
	// //defer node.Close()
	// node.addr = ":" + strconv.Itoa(8888)
	// node.Loglevel = LOGLEVEL_OFF
	// mgr := chain.CreateManager()
	// mgr.InitAccounts()
	// node.Mgr = &mgr

	//reply_msg := HandleFaucet(node, req_msg)
	// if reply_msg != "{:REP FAUCET :data ok}" {
	// 	t.Error("reply_msg ", reply_msg)
	// }

	// chain.MakeBlock(&mgr)

	// time.Sleep(2000 * time.Millisecond)

	// req_msg = netio.EncodeMConstructMsgxx(netio.REQ, netio.CMD_BALANCE, addr)
	// msg = netio.ParseMessag(req_msg)

	// log.Println(mgr.Accounts)

	// reply_msg_string := HandleBalance(node, msg)
	// log.Println(reply_msg_string)
	// msg = netio.ParseMessag(reply_msg_string)
	// // if reply_msg_string != "REP#BALANCE#1|" {
	// // 	t.Error(msg)
	// // }

	// // bal := netio.ParseMessageBalance(reply_msg)
	// if msg.MessageType != "REP" || msg.Command != netio.CMD_BALANCE {
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
// 	req_msg := netio.ConstructMsgxx(netio.REQ, netio.CMD_FAUCET, addr)
// 	msg := netio.ParseMessag(req_msg)

// 	reply_msg := HandleFaucet(node, msg)
// 	if reply_msg != "REP#FAUCET#ok|" {
// 		t.Error("reply_msg ", reply_msg)
// 	}
// 	chain.MakeBlock(&mgr)
// 	time.Sleep(100 * time.Millisecond)
// 	req_msg = netio.ConstructMsgxx(netio.REQ, netio.CMD_BALANCE, addr)
// 	msg = netio.ParseMessag(req_msg)
// 	reply_msg = HandleBalance(node, msg)
// 	msg = netio.ParseMessa(reply_msg)

// 	if msg.MessageType != "REP" || msg.Command != netio.CMD_BALANCE {
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
// 	req_msg = netio.EncodeMessageTx(txJson)
// 	msg = netio.ParseMessa(req_msg)

// 	reply_msg = HandleTx(node, msg)
// 	// //TODO!
// 	if reply_msg != "REP#TX#ok|" {
// 		t.Error("reply_msg ", reply_msg)
// 	}

// 	chain.MakeBlock(&mgr)
// 	time.Sleep(100 * time.Millisecond)

// 	req_msg = netio.ConstructMsgxx(netio.REQ, netio.CMD_BALANCE, addr2)
// 	msg = netio.ParseMessae(req_msg)
// 	reply_msg = HandleBalance(node, msg)
// 	msg = netio.ParseMessa(reply_msg)
// 	//if reply_msg != "REP#BALANCE#5|" {
// 	//bal := netio.ParseMeseBalance(reply_msg)

// 	if msg.MessageType != "REP" || msg.Command != netio.CMD_BALANCE {
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
