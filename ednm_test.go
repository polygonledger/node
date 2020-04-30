package main

import (
	"testing"
	//"fmt"

	"github.com/polygonledger/node/parser"
	"github.com/polygonledger/node/ntcl"
)

//basic block functions
func TestMap(t *testing.T) {

	m := map[string]string{"test": "value"}

	mstr := parser.MakeMap(m)

	if mstr != "{:test value}" {
		t.Error("error creating map")
	}

}

func TestEncodeReq(t *testing.T){
	msgs := ntcl.EncodeMsgMap(ntcl.REQ, "PING")
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
