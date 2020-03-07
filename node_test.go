package main

import (
	"log"
	"testing"

	"github.com/polygonledger/node/ntcl"
)

func TestReaderPing(t *testing.T) {
	log.Println("TestReaderPing")

	ntchan := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan)

	//TODO
	//go pinghandler(ntchan)

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)
	msg := ntcl.ParseMessage(req_msg)
	//ntchan.REQ_in <- msg

	reply_msg := HandlePing(msg)

	// time.Sleep(50 * time.Millisecond)

	// if !isEmpty(ntchan.REQ_in, 1*time.Second) {
	// 	t.Error("REQ_in not processed")
	// }

	// x := <-ntchan.REP_out
	// if x != ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_PONG, "") {
	// 	t.Error("not poing")
	// }

}
