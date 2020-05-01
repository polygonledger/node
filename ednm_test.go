package main

import (
	"testing"

	//"fmt"

	"github.com/polygonledger/node/ntcl"
	"github.com/polygonledger/node/parser"
)

//basic block functions
func TestMap(t *testing.T) {

	m := map[string]string{"test": "value"}

	mstr := parser.MakeMap(m)

	if mstr != "{:test value}" {
		t.Error("error creating map")
	}

}

func TestEncodeReq(t *testing.T) {
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

func TestDecodeMap(t *testing.T) {
	s := "{:REQ BALANCE :data P2e2bfb58c9db}"

	v, k := parser.ReadMap(s)

	if len(v) != 2 {
		t.Error("parsing map ", v)
	}

	if len(k) != 2 {
		t.Error("parsing map ", v)
	}

	if v[0] != "BALANCE" || v[1] != "P2e2bfb58c9db" {
		t.Error(s)
	}

	if k[0] != "REQ" || k[1] != "data" {
		t.Error(k)
	}

	//req_msg_string := ntcl.EncodeMsgMapData(ntcl.REQ, ntcl.CMD_BALANCE, ra)
	//req_msg_balance := ntcl.ParseMessageMapData(req_msg_string)
	// msg := ntcl.ParseMessageMapData(s)

	// //if x.Data != []byte("P2e2bfb58c9db") {
	// if msg.Data == nil {
	// 	t.Error("msg ", msg)
	// }

}
