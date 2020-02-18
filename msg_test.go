package main

import (
	"encoding/json"
	"testing"

	block "github.com/polygonledger/node/block"
	ntwk "github.com/polygonledger/node/ntwk"
)

func TestMessageBasic(t *testing.T) {

	msg := ntwk.Message{MessageType: ntwk.REQ, Command: "CMD"}
	if !ntwk.IsValidMsgType(msg.MessageType) {
		t.Error("msg type invalid")
	}
}

func TestMessageJson(t *testing.T) {

	msg := ntwk.Message{MessageType: ntwk.REQ, Command: "CMD"}
	msgJson, _ := json.Marshal(msg)

	var msgUn ntwk.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}

	if msgUn.MessageType != ntwk.REQ {
		t.Error("JSON marshal failed")
	}
}

func TestMessageType(t *testing.T) {
	msg := ntwk.RequestMessage()
	if msg.MessageType != ntwk.REQ {
		t.Error("msg failed")
	}

	msg = ntwk.ReplyMessage()
	if msg.MessageType != ntwk.REP {
		t.Error("msg failed")
	}

}

func TestMessageAccount(t *testing.T) {
	a := block.Account{AccountKey: "test"}
	msg := ntwk.AccountMessage(a)
	if msg.Account != a {
		t.Error("msg failed")
	}

	msgJson, _ := json.Marshal(msg)

	var msgUn ntwk.MessageAccount
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Account.AccountKey != a.AccountKey {
		t.Error("JSON marshal failed")
	}

	//var genericmsg protocol.Message
	//genericmsg = protocol.Message(msgUn)

}
