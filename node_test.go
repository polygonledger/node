package main

import (
	"log"
	"testing"

	"github.com/polygonledger/node/ntcl"
)

func TestBasicCommand(t *testing.T) {

	log.Println("TestBasicCommand")

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)
	msg := ntcl.ParseMessage(req_msg)
	//ntchan.REQ_in <- msg
	reply_msg := HandlePing(msg)
	if reply_msg != "REP#PONG#|" {
		t.Error(reply_msg)
	}

	ntchan := ntcl.ConnNtchanStub("")
	go RequestHandlerTel(ntchan)
	ntchan.REQ_in <- req_msg
	reply := <-ntchan.REP_out

	if reply != "REP#PONG#|" {
		t.Error(reply_msg)
	}

}

func TestBalance(t *testing.T) {

	log.Println("TestBalance")

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_BALANCE, "abc")
	msg := ntcl.ParseMessage(req_msg)

	reply_msg := HandleBalance(msg)
	if reply_msg != "REP#BALANCE#0|" {
		t.Error(reply_msg)
	}

	//TODO with chain setup

	// ntchan := ntcl.ConnNtchanStub("")
	// go RequestHandlerTel(ntchan)
	// ntchan.REQ_in <- req_msg
	// reply := <-ntchan.REP_out

	// if reply != "REP#PONG#|" {
	// 	t.Error(reply_msg)
	// }
}
