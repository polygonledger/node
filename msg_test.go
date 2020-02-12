package main

import (
	"encoding/json"
	"testing"

	block "github.com/polygonledger/node/block"
	protocol "github.com/polygonledger/node/net"
)

func TestMessageJson(t *testing.T) {

	msg := protocol.Message{MessageType: "msg", Command: "CMD"}
	msgJson, _ := json.Marshal(msg)

	var msgUn protocol.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}

	if msgUn.MessageType != "msg" {
		t.Error("JSON marshal failed")
	}
}

func TestMessageType(t *testing.T) {
	msg := protocol.RequestMessage()
	if msg.MessageType != protocol.REQ {
		t.Error("msg failed")
	}

	msg = protocol.ReplyMessage()
	if msg.MessageType != protocol.REP {
		t.Error("msg failed")
	}

}

func TestMessageAccount(t *testing.T) {
	a := block.Account{AccountKey: "test"}
	msg := protocol.AccountMessage(a)
	if msg.Account != a {
		t.Error("msg failed")
	}

	msgJson, _ := json.Marshal(msg)

	var msgUn protocol.MessageAccount
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}

	if msgUn.Account.AccountKey != a.AccountKey {
		t.Error("JSON marshal failed")
	}

	//var genericmsg protocol.Message
	//genericmsg = protocol.Message(msgUn)

}
