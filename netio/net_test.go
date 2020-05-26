package netio

import (
	"fmt"
	"testing"
	"time"

	"github.com/polygonledger/node/xutils"
)

func SimulateNetworkInput(ntchan *Ntchan) {

	for {
		ntchan.Reader_queue <- "test"
		//log.Println(len(ntchan.Reader_queue))
		time.Sleep(100 * time.Millisecond)
	}
}

func TestReaderin(t *testing.T) {

	ntchan := ConnNtchanStub("")
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

func SimulateRequest(ntchan *Ntchan) {

	msg := Message{MessageType: REQ, Command: CMD_PING}
	jsonmsg := ToJSONMessage(msg)
	//fmt.Println("sim ", req_msg)
	ntchan.Reader_queue <- jsonmsg
	//log.Println(len(ntchan.Reader_queue))
	time.Sleep(100 * time.Millisecond)

}

// in reader should bet forwarded to req_chan
func TestRequestIn(t *testing.T) {
	fmt.Println("TestRequestIn")
	ntchan := ConnNtchanStub("")

	go ReadProcessor(ntchan)

	if !xutils.IsEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("channel full")
	}

	//put 1 request in reader
	fmt.Println("SimulateRequest")
	go SimulateRequest(&ntchan)

	//wait
	time.Sleep(100 * time.Millisecond)

	//check if reader empty
	fmt.Println("isempty?")
	if !xutils.IsEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("Reader_queue not empty")
	}

	fmt.Println("wait for req")
	req_in := <-ntchan.REQ_in
	msg := Message{MessageType: REQ, Command: CMD_PING}
	jsonmsg := ToJSONMessage(msg)

	if req_in != jsonmsg {
		t.Error("req not>> ", req_in, jsonmsg)
	} else {
		fmt.Println("ok")
	}

	//req_in
	if !xutils.IsEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("req channel empty")
	}

}

func TestRequestOut(t *testing.T) {

	ntchan := ConnNtchanStub("")

	//req_out_msg := EdnConstructMsgMap(REQ, CMD_PING)

	select {
	case msg := <-ntchan.REQ_out:
		fmt.Println("received message", msg)
		t.Error("should not contain")
	case <-time.After(100 * time.Millisecond):
		//trace("no message received")
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

	ntchan := ConnNtchanStub("")

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

func basicPingReqProcessor(ntchan Ntchan, t *testing.T) {

	x := <-ntchan.REQ_in
	if x != "ping" {
		t.Error("req in")
	}
	ntchan.REP_out <- "pong"
}

func TestPingAll(t *testing.T) {

	ntchan := ConnNtchanStub("")

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

	ntchan := ConnNtchanStub("")

	go ReadProcessor(ntchan)

	//req_msg := EdnConstructMsgMap(REQ, CMD_PING)
	msg := Message{MessageType: REQ, Command: CMD_PING}
	req_msg := ToJSONMessage(msg)

	ntchan.Reader_queue <- req_msg
	time.Sleep(50 * time.Millisecond)

	//reader queue should be empty
	if !xutils.IsEmpty(ntchan.Reader_queue, 1*time.Second) {
		t.Error("reader not processed")
	}

	req_in := <-ntchan.REQ_in
	if req_in != req_msg {
		t.Error("req not equal")
	}

}

func pinghandler(ntchan Ntchan) {
	for {
		//REQUEST
		req_msg := <-ntchan.REQ_in
		//if ping
		if req_msg == "" {
			//<-req_msg
		}
		rep_msg := EdnConstructMsgMap(REP, CMD_PONG)
		ntchan.REP_out <- rep_msg
		//log.Println("REP_out >> ", rep_msg)
	}
}

func TestReaderPing(t *testing.T) {

	ntchan := ConnNtchanStub("")

	go ReadProcessor(ntchan)

	go pinghandler(ntchan)

	ntchan.REQ_in <- EdnConstructMsgMap(REQ, CMD_PING)

	time.Sleep(50 * time.Millisecond)

	if !xutils.IsEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	x := <-ntchan.REP_out
	if x != EdnConstructMsgMap(REP, CMD_PONG) {
		t.Error("not poing")
	}

}

//test entire loop from reader to writer
func TestAllPingPoingIn(t *testing.T) {

	ntchan := ConnNtchanStub("")

	go ReadProcessor(ntchan)
	go WriteProcessor(ntchan)
	go pinghandler(ntchan)

	ntchan.Reader_queue <- EdnConstructMsgMap(REQ, CMD_PING)

	time.Sleep(50 * time.Millisecond)

	if !xutils.IsEmpty(ntchan.REQ_in, 1*time.Second) {
		t.Error("REQ_in not processed")
	}

	//TODO! fix
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
func ConnectWrite(ntchan1 Ntchan, ntchan2 Ntchan) {
	for {
		xout := <-ntchan1.Writer_queue
		ntchan2.Reader_queue <- xout

		yout := <-ntchan2.Writer_queue
		ntchan1.Reader_queue <- yout
	}
}

func TestAllPingPongDuplex(t *testing.T) {
	//log.Println("TestAllPingPongDuplex")

	ntchan1 := ConnNtchanStub("")
	ntchan2 := ConnNtchanStub("")

	go ReadProcessor(ntchan1)
	go WriteProcessor(ntchan1)
	go ConnectWrite(ntchan1, ntchan2)

	ntchan1.Writer_queue <- "test"
	time.Sleep(600 * time.Millisecond)

	if !xutils.IsEmpty(ntchan1.Writer_queue, 300*time.Millisecond) {
		t.Error("writer not empty")
	}

	x := <-ntchan2.Reader_queue
	if x != "test" {
		t.Error("did not connect")
	}
	//go pinghandler(ntchan)

}

//timed test
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
