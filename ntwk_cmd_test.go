package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

func TestRequestBalance(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	req_out_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_BALANCE, ntwk.EMPTY_DATA)
	log.Println(req_out_msg)

	select {
	case msg := <-ntchan.REQ_out:
		fmt.Println("received message", msg)
		t.Error("should not contain")
	case <-time.After(100 * time.Millisecond):
		//fmt.Println("no message received")
	}

	go func() {
		ntchan.REQ_out <- "test"
	}()

	go func() {
		x := <-ntchan.REQ_out
		if x != "test" {
			t.Error("should receive test")
		}
		//}
	}()

	go func() {
		ntchan.REQ_out <- "test"
	}()

	go func() {
		select {
		case msg := <-ntchan.REQ_out:
			if msg != "test" {
				t.Error("wrong")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("should contain")
		}
	}()

}
