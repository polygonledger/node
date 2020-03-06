package main

import (
	"encoding/json"
	"testing"

	block "github.com/polygonledger/node/block"
	"github.com/polygonledger/node/ntcl"
)

func TestMessageBasic(t *testing.T) {

	msg := ntcl.Message{MessageType: ntcl.REQ, Command: "CMD"}
	if !ntcl.IsValidMsgType(msg.MessageType) {
		t.Error("msg type invalid")
	}
}

func TestMessageJson(t *testing.T) {

	msg := ntcl.Message{MessageType: ntcl.REQ, Command: "CMD"}
	msgJson, _ := json.Marshal(msg)

	var msgUn ntcl.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}

	if msgUn.MessageType != ntcl.REQ {
		t.Error("JSON marshal failed")
	}
}

func TestMessageType(t *testing.T) {
	msg := ntcl.RequestMessage()
	if msg.MessageType != ntcl.REQ {
		t.Error("msg failed")
	}

	msg = ntcl.ReplyMessage()
	if msg.MessageType != ntcl.REP {
		t.Error("msg failed")
	}

}

func TestMessageAccount(t *testing.T) {
	a := block.Account{AccountKey: "test"}
	msg := ntcl.AccountMessage(a)
	if msg.Account != a {
		t.Error("msg failed")
	}

	msgJson, _ := json.Marshal(msg)

	var msgUn ntcl.MessageAccount
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Account.AccountKey != a.AccountKey {
		t.Error("JSON marshal failed")
	}

	//var genericmsg protocol.Message
	//genericmsg = protocol.Message(msgUn)

}
