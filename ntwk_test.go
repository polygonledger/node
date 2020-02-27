package main

import (
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

func x(ntchan ntwk.Ntchan) {
	ntchan.Reader_processed = 1
}

func BasicChanCounter(t *testing.T) {

	t.Error("counter fails")

	// var ntchan ntwk.Ntchan
	// ntchan.Reader_queue = make(chan string)
	// ntchan.Writer_queue = make(chan string)
	// ntchan.Reader_processed = 0
	// ntchan.Writer_processed = 0
	// go x(ntchan)
	// time.Sleep(1000 * time.Second)
	// if ntchan.Reader_processed != 1 {
	// 	t.Error("counter fails")
	// }
	// t.Error("counter fails")
}

func SimulateNetworkInput(ntchan ntwk.Ntchan) {
	ntchan.Reader_queue <- "test"
}

func TestBasicNtwk(t *testing.T) {

	var ntchan ntwk.Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.Reader_processed = 0
	ntchan.Writer_processed = 0

	go ntwk.ReadProcessor(&ntchan, 1*time.Millisecond)
	if ntchan.Reader_processed != 0 {
		t.Error("reader processed")
	}

	if len(ntchan.Reader_queue) != 0 {
		t.Error("reader queue not empty")
	}

	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

	log.Println(req_msg)
	ntchan.Reader_queue <- req_msg
	time.Sleep(5 * time.Millisecond)
	if ntchan.Reader_processed != 1 {
		t.Error("reader not processed")
	}

	log.Println(req_msg)

}
