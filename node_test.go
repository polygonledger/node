package main

import (
	"testing"

	protocol "github.com/polygonledger/node/ntwk"
)

func TestBasicPing(t *testing.T) {

	msg_out_chan := make(chan string)

	go HandlePing(msg_out_chan)
	msg := <-msg_out_chan

	if !(msg == "PONG") {
		t.Error("ping failed")
	}
}

func TestBasicPingMsg(t *testing.T) {

	msg_in := make(chan string)
	msg_out := make(chan string)

	emptydata := ""
	req_msg := protocol.EncodeMsg(protocol.REQ, protocol.CMD_PING, emptydata)
	go func() {
		msg_in <- req_msg
	}()
	go HandleMsg(msg_in, msg_out)
	msg := <-msg_out

	if !(msg == "PONG") {
		t.Error("ping failed")
	}
}
