package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

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

	maxt := 100 * time.Millisecond
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
	//log.Println(len(ntchan.Reader_queue))
	time.Sleep(100 * time.Millisecond)

}

//check if a channel is empty
func isEmpty(c chan string, d time.Duration) bool {
	select {
	case ret := <-c:
		log.Println("got ", ret)
		return false

	case <-time.After(d):
		//log.Println("timeout")
		return true
	}
}

// in reader should bet forwarded to req_chan
func TestRequestIn(t *testing.T) {
	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 100*time.Millisecond)

	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("channel full")
	}

	//put 1 request in reader
	go SimulateRequest(&ntchan)

	//wait
	time.Sleep(100 * time.Millisecond)

	//check if reader empty
	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("Reader_queue not empty")
	}

	req_in := <-ntchan.REQ_in
	if req_in != ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA) {
		t.Error("req not")
	}

	//req_in
	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("req channel empty")
	}

}

func TestRequestOut(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	req_out_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)
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

func TestReplyIn(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	go func() {
		ntchan.REP_in <- "test"
	}()

	go func() {
		select {
		case msg := <-ntchan.REP_in:
			if msg != "test" {
				t.Error("wrong")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("should contain")
		}
	}()

}

func basicPingReqProcessor(ntchan ntwk.Ntchan, t *testing.T) {

	x := <-ntchan.REQ_in
	if x != "ping" {
		t.Error("req in")
	}
	ntchan.REP_out <- "pong"
}

func TestPing(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	//REQ in
	//reply request ping with pong
	go func() {
		ntchan.REQ_in <- "ping"
	}()

	go basicPingReqProcessor(ntchan, t)

	go func() {
		x := <-ntchan.REP_out
		if x != "pong" {
			t.Error("expect pong")
		}
	}()

	//REQ out
	//request ping should return pong
	go func() {
		ntchan.REQ_out <- "ping"
	}()

	go func() {
		x := <-ntchan.REP_in
		if x != "pong" {
			t.Error("expect pong")
		}
	}()

	go func() {
		x := <-ntchan.REQ_out
		if x != "ping" {
			t.Error("req out")
		}

		ntchan.REP_in <- "pong"
	}()

}

func TestReaderRequest(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)

	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

	ntchan.Reader_queue <- req_msg
	time.Sleep(50 * time.Millisecond)

	//reader queue should be empty
	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("reader not processed")
	}

	req_in := <-ntchan.REQ_in
	if req_in != ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA) {
		t.Error("req not equal")
	}

}

func pinghandler(ntchan ntwk.Ntchan) {
	for {
		//REQUEST
		req_msg := <-ntchan.REQ_in
		//if ping
		if req_msg == "" {
			//<-req_msg
		}
		rep_msg := ntwk.EncodeMsgString(ntwk.REP, ntwk.CMD_PONG, "")
		ntchan.REP_out <- rep_msg
		//log.Println("REP_out >> ", rep_msg)
	}
}

func TestReaderPing(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)

	go pinghandler(ntchan)

	ntchan.REQ_in <- ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

	time.Sleep(50 * time.Millisecond)

	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	x := <-ntchan.REP_out
	if x != ntwk.EncodeMsgString(ntwk.REP, ntwk.CMD_PONG, "") {
		t.Error("not poing")
	}

}

func TestAllPingPoing(t *testing.T) {

	ntchan := ntwk.ConnNtchanStub("")

	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)
	go ntwk.Writeprocessor(ntchan, 1*time.Millisecond)

	go pinghandler(ntchan)

	ntchan.Reader_queue <- ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, ntwk.EMPTY_DATA)

	time.Sleep(50 * time.Millisecond)

	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	// if isEmpty(ntchan.Writer_queue, time.Second) {
	// 	t.Error("writer queue empty")
	// }

	write_out := <-ntchan.Writer_queue
	if write_out != "REP#PONG#|" {
		t.Error("pong out ", write_out)
	}

}

//
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
