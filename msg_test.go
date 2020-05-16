package main

import (
	"encoding/json"
	"testing"

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
