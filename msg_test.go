package main

import (
	"encoding/json"
	"testing"

	block "github.com/polygonledger/node/block"
	"github.com/polygonledger/node/netio"
)

func TestMessageBasic(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: "CMD"}
	if !netio.IsValidMsgType(msg.MessageType) {
		t.Error("msg type invalid")
	}
}

func TestMessageJson(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: "CMD"}
	msgJson, _ := json.Marshal(msg)

	var msgUn netio.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}

	if msgUn.MessageType != netio.REQ {
		t.Error("JSON marshal failed")
	}
}

func TestMessageType(t *testing.T) {
	msg := netio.RequestMessage()
	if msg.MessageType != netio.REQ {
		t.Error("msg failed")
	}

	msg = netio.ReplyMessage()
	if msg.MessageType != netio.REP {
		t.Error("msg failed")
	}

}

func TestMessageAccount(t *testing.T) {
	a := block.Account{AccountKey: "test"}
	msg := netio.AccountMessage(a)
	if msg.Account != a {
		t.Error("msg failed")
	}

	msgJson, _ := json.Marshal(msg)

	var msgUn netio.MessageAccount
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Account.AccountKey != a.AccountKey {
		t.Error("JSON marshal failed")
	}

	//var genericmsg protocol.Message
	//genericmsg = protocol.Message(msgUn)

}
