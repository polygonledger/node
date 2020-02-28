package main

import (
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

// func TestReader(t *testing.T) {

// 	ntchan := ntwk.ConnNtchanStub("")

// 	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)
// 	if ntchan.Reader_processed != 0 {
// 		t.Error("reader processed")
// 	}

// 	if len(ntchan.Reader_queue) != 0 {
// 		t.Error("reader queue not empty")
// 	}

// 	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

// 	// ntchan.Reader_queue <- req_msg
// 	// time.Sleep(5 * time.Millisecond)
// 	// if ntchan.Reader_processed != 1 {
// 	// 	t.Error("reader not processed")
// 	// }

// }

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

func SimulateRequest(ntchan *ntwk.Ntchan) {

	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)
	ntchan.Reader_queue <- req_msg
	log.Println(len(ntchan.Reader_queue))
	time.Sleep(100 * time.Millisecond)

}

func isEmpty(c chan string, d time.Duration) bool {
	select {
	case ret := <-c:
		log.Println("got ", ret)
		return false

	case <-time.After(d):
		//ch <- false
		log.Println("timeout")
		return true
	}
}

// in reader should bet forwarded to req_chan
func TestRequest(t *testing.T) {
	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 100*time.Millisecond)

	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("channel full")
	}

	//put 1 request in reader
	go SimulateRequest(&ntchan)

	//wait
	time.Sleep(100 * time.Millisecond)

	//reader empty
	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("channel empty")
	}

	//req_in
	if isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("req channel empty")
	}

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
