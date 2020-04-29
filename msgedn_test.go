package main

import (
	"testing"
	"github.com/polygonledger/node/ntcl"
)



func TestMessageMapBasic(t *testing.T) {
	// msg := ntcl.Message{MessageType: ntcl.REQ, Command: "CMD"}
	// if !ntcl.IsValidMsgType(msg.MessageType) {
	// 	t.Error("msg type invalid")
	// }

	msgs := ntcl.EncodeMsgMap("REQ", "PING")
	//fmt.Println(msgs)
	if msgs != "{:REQ PING}" {
		t.Error("wrong encoding ", msgs)
	}

	msg := ntcl.ParseMessageMap(msgs)
	if msg.MessageType != "REQ" {
		t.Error("type")
	}
	if msg.Command != "PING" {
		t.Error("command")
	}

}
