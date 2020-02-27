package main

import (
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

func TestBasicNtwk(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)
	if ntchan.Reader_processed != 0 {
		t.Error("reader processed")
	}

	if len(ntchan.Reader_queue) != 0 {
		t.Error("reader queue not empty")
	}

	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

	ntchan.Reader_queue <- req_msg
	time.Sleep(5 * time.Millisecond)
	if ntchan.Reader_processed != 1 {
		t.Error("reader not processed")
	}

	log.Println(req_msg)

}

func SimulateNetworkInput(ntchan *ntwk.Ntchan) {
	for {
		ntchan.Reader_queue <- "test"
		//log.Println(len(ntchan.Reader_queue))
		time.Sleep(100 * time.Millisecond)
	}
}

func TestReaderin(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")
	go SimulateNetworkInput(&ntchan)
	time.Sleep(300 * time.Millisecond)
	start := time.Now()

	maxt := 300 * time.Millisecond
	tt := time.Now()
	elapsed := tt.Sub(start)

	for ok := true; ok; ok = elapsed < maxt {
		t2 := time.Now()
		elapsed = t2.Sub(start)
		time.Sleep(10 * time.Millisecond)
	}
	x := <-ntchan.Reader_queue

	if x != "test" {
		t.Error("reader queue empty")
	}
}

func SimulateRequests(ntchan *ntwk.Ntchan) {
	for {
		req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)
		ntchan.Reader_queue <- req_msg
		log.Println(len(ntchan.Reader_queue))
		time.Sleep(100 * time.Millisecond)
	}
}

//ping pong
func TestRequest(t *testing.T) {
	ntchan := ntwk.ConnNtchanStub("")

	go SimulateRequests(&ntchan)
	// go ntwk.RequestProcessor(ntchan, 1*time.Second)
	// go ntwk.ReplyProcessor(&ntchan, 1*time.Second)
	// read_time_chan := 300 * time.Millisecond
	// go ntwk.ReadProcessor(ntchan, read_time_chan)
	// start := time.Now()

	// maxt := 300 * time.Millisecond
	// tt := time.Now()
	// elapsed := tt.Sub(start)

	// for ok := true; ok; ok = elapsed < maxt {
	// 	t2 := time.Now()
	// 	elapsed = t2.Sub(start)
	// 	time.Sleep(10 * time.Millisecond)
	// }
	// x := <-ntchan.Reader_queue

	// if x != "REQ#PING#EMPTY|" {
	// 	t.Error("TestRequest reader queue ", x)
	// }

	// resp_to_write := <-ntchan.Writer_queue

	// if resp_to_write != "REP#PONG#EMPTY|" {
	// 	t.Error("wrong reply ", resp_to_write)
	// }

}
