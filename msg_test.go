package main

import (
	"encoding/json"
	"testing"

	block "github.com/polygonledger/node/block"
	net "github.com/polygonledger/node/net"
)

func TestMessageBasic(t *testing.T) {

	msg := net.Message{MessageType: net.REQ, Command: "CMD"}
	if !net.IsValidMsgType(msg.MessageType) {
		t.Error("msg type invalid")
	}
}

func TestMessageJson(t *testing.T) {

	msg := net.Message{MessageType: net.REQ, Command: "CMD"}
	msgJson, _ := json.Marshal(msg)

	var msgUn net.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}

	if msgUn.MessageType != net.REQ {
		t.Error("JSON marshal failed")
	}
}

func TestMessageType(t *testing.T) {
	msg := net.RequestMessage()
	if msg.MessageType != net.REQ {
		t.Error("msg failed")
	}

	msg = net.ReplyMessage()
	if msg.MessageType != net.REP {
		t.Error("msg failed")
	}

}

func TestMessageAccount(t *testing.T) {
	a := block.Account{AccountKey: "test"}
	msg := net.AccountMessage(a)
	if msg.Account != a {
		t.Error("msg failed")
	}

	msgJson, _ := json.Marshal(msg)

	var msgUn net.MessageAccount
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Account.AccountKey != a.AccountKey {
		t.Error("JSON marshal failed")
	}

	//var genericmsg protocol.Message
	//genericmsg = protocol.Message(msgUn)

}
