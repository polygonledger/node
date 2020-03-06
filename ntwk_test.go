package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/polygonledger/node/ntcl"
)

func SimulateNetworkInput(ntchan *ntcl.Ntchan) {
	for {
		ntchan.Reader_queue <- "test"
		//log.Println(len(ntchan.Reader_queue))
		time.Sleep(100 * time.Millisecond)
	}
}

func TestReaderin(t *testing.T) {
	log.Println("TestReaderin")

	ntchan := ntcl.ConnNtchanStub("")
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

func SimulateRequest(ntchan *ntcl.Ntchan) {

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)
	ntchan.Reader_queue <- req_msg
	//log.Println(len(ntchan.Reader_queue))
	time.Sleep(100 * time.Millisecond)

}

// in reader should bet forwarded to req_chan
func TestRequestIn(t *testing.T) {
	log.Println("TestRequestIn")
	ntchan := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan)

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
	if req_in != ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA) {
		t.Error("req not")
	}

	//req_in
	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("req channel empty")
	}

}

func TestRequestOut(t *testing.T) {
	log.Println("TestRequestOut")

	ntchan := ntcl.ConnNtchanStub("")

	req_out_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)
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
	log.Println("TestReplyIn")

	ntchan := ntcl.ConnNtchanStub("")

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

func basicPingReqProcessor(ntchan ntcl.Ntchan, t *testing.T) {

	x := <-ntchan.REQ_in
	if x != "ping" {
		t.Error("req in")
	}
	ntchan.REP_out <- "pong"
}

func TestPing(t *testing.T) {
	log.Println("TestPing")

	ntchan := ntcl.ConnNtchanStub("")

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
	log.Println("TestReaderRequest")

	ntchan := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan)

	req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)

	ntchan.Reader_queue <- req_msg
	time.Sleep(50 * time.Millisecond)

	//reader queue should be empty
	if !isEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("reader not processed")
	}

	req_in := <-ntchan.REQ_in
	if req_in != ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA) {
		t.Error("req not equal")
	}

}

func pinghandler(ntchan ntcl.Ntchan) {
	for {
		//REQUEST
		req_msg := <-ntchan.REQ_in
		//if ping
		if req_msg == "" {
			//<-req_msg
		}
		rep_msg := ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_PONG, "")
		ntchan.REP_out <- rep_msg
		//log.Println("REP_out >> ", rep_msg)
	}
}

func TestReaderPing(t *testing.T) {
	log.Println("TestReaderPing")

	ntchan := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan)

	go pinghandler(ntchan)

	ntchan.REQ_in <- ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)

	time.Sleep(50 * time.Millisecond)

	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	x := <-ntchan.REP_out
	if x != ntcl.EncodeMsgString(ntcl.REP, ntcl.CMD_PONG, "") {
		t.Error("not poing")
	}

}

//test entire loop from reader to writer
func TestAllPingPoingIn(t *testing.T) {
	log.Println("TestAllPingPoingIn")
	//TODO! fix

	ntchan := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan)
	go ntcl.WriteProcessor(ntchan)
	go pinghandler(ntchan)

	ntchan.Reader_queue <- ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_PING, ntcl.EMPTY_DATA)

	time.Sleep(50 * time.Millisecond)

	if !isEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	// // if isEmpty(ntchan.Writer_queue, time.Second) {
	// // 	t.Error("writer queue empty")
	// // }

	// write_out := <-ntchan.Writer_queue
	// if write_out != "REP#PONG#|" {
	// 	t.Error("pong out ", write_out)
	// }

}

//connect to ntchans to each other
//this simulates a real network connection
func ConnectWrite(ntchan1 ntcl.Ntchan, ntchan2 ntcl.Ntchan) {
	for {
		xout := <-ntchan1.Writer_queue
		ntchan2.Reader_queue <- xout

		yout := <-ntchan2.Writer_queue
		ntchan1.Reader_queue <- yout
	}
}

func TestAllPingPongDuplex(t *testing.T) {
	log.Println("TestAllPingPongDuplex")

	ntchan1 := ntcl.ConnNtchanStub("")
	ntchan2 := ntcl.ConnNtchanStub("")

	go ntcl.ReadProcessor(ntchan1)
	go ntcl.WriteProcessor(ntchan1)
	go ConnectWrite(ntchan1, ntchan2)

	ntchan1.Writer_queue <- "test"
	time.Sleep(600 * time.Millisecond)

	if !isEmpty(ntchan1.Writer_queue, 300*time.Millisecond) {
		t.Error("writer not empty")
	}

	x := <-ntchan2.Reader_queue
	if x != "test" {
		t.Error("did not connect")
	}
	//go pinghandler(ntchan)

}

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
