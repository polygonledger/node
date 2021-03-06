package netio

// network communication layer (netio)

// netio -> semantics of channels
// TCP/IP -> golang net

// see whitepaper for details

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Ntchan struct {
	//TODO is only single connection
	Conn net.Conn
	//Name     string
	SrcName  string //TODO doc
	DestName string
	//TODO message type
	Reader_queue chan string
	Writer_queue chan string
	//inflow
	REQ_in  chan string
	REP_out chan string
	//outflow
	REQ_out chan string
	REP_in  chan string

	PUB_out chan string
	SUB_out chan string
	//PUB_time_quit chan int
	verbose bool
	// SUB_request_out   chan string
	// SUB_request_in    chan string
	// UNSUB_request_out chan string
	// UNSUB_request_in  chan string

	// Reader_processed int
	// Writer_processed int
}

func vlog(ntchan Ntchan, s string) {
	//verbose := ntchan.verbose
	//fmt.Println(s)
	// if verbose {
	// 	log.Println("vlog ", s)
	// }
	log.Println("vlog ", s)
}

func logmsgd(ntchan Ntchan, src string, msg string) {
	s := fmt.Sprintf("[%s] ### %v", src, msg)
	vlog(ntchan, s)
}

func logmsgc(ntchan Ntchan, name string, src string, msg string) {
	s := fmt.Sprintf("%s [%s] ### %v", name, src, msg)
	vlog(ntchan, s)
}

func ConnNtchan(conn net.Conn, SrcName string, DestName string, verbose bool) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.PUB_out = make(chan string)
	//ntchan.PUB_time_quit = make(chan int)
	// ntchan.Reader_processed = 0
	// ntchan.Writer_processed = 0
	ntchan.Conn = conn
	ntchan.SrcName = SrcName
	ntchan.DestName = DestName

	return ntchan
}

//for testing
func ConnNtchanStub(name string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.PUB_out = make(chan string)
	//ntchan.PUB_time_quit = make(chan int)
	//ntchan.Reader_processed = 0
	//ntchan.Writer_processed = 0

	return ntchan
}

//all major processes to operate
func NetConnectorSetup(ntchan Ntchan) {

	// read_loop_time := 800 * time.Millisecond
	// read_time_chan := 300 * time.Millisecond
	// write_loop_time := 300 * time.Millisecond

	//reads from the actual "physical" network
	go ReadLoop(ntchan)
	//process of reads
	go ReadProcessor(ntchan)
	//processor of REQ_out REP_out
	go WriteProcessor(ntchan)
	//write to network whatever is in writer queue

	go WriteLoop(ntchan, 300*time.Millisecond)

	//TODO
	//go WriteProducer(ntchan)
}

func ReadLoop(ntchan Ntchan) {
	vlog(ntchan, "init ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
	d := 300 * time.Millisecond
	//msg_reader_total := 0
	for {
		//read from network and put in channel
		vlog(ntchan, "iter ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
		msg, err := NetMsgRead(ntchan)
		if err != nil {

		}
		//handle cases
		//currently can be empty or len, shoudl fix one style
		if len(msg) > 0 { //&& msg != EMPTY_MSG {
			vlog(ntchan, "ntwk read => "+msg)
			logmsgc(ntchan, ntchan.SrcName, "ReadLoop", msg)
			vlog(ntchan, "put on Reader queue "+msg)
			//put in the queue to process
			ntchan.Reader_queue <- msg
		}

		time.Sleep(d)
		//fix: need ntchan to be a pointer
		//msg_reader_total++
	}
}

//only output the message read
func ProcessorEchoMock(ntchan Ntchan) {

}

//read from reader queue and process by forwarding to right channel
func ReadProcessor(ntchan Ntchan) {

	for {
		logmsgd(ntchan, "ReadProcessor", "loop")
		msgString := <-ntchan.Reader_queue
		logmsgd(ntchan, "ReadProcessor", msgString)

		if len(msgString) > 0 {
			logmsgc(ntchan, ntchan.SrcName, "ReadProcessor", msgString) //, ntchan.Reader_processed)
			//TODO try catch

			//msg := EdnParseMessageMap(msgString)
			msg := FromJSON(msgString)

			if msg.MessageType == REQ {

				//msg_string := MsgString(msg)
				//msg_string := EdnConstructMsgMapS(msg)
				logmsgd(ntchan, "ReadProcessor", "REQ_in")

				//TODO!
				ntchan.REQ_in <- msgString
				// reply_string := "echo:" + msg_string
				// ntchan.Writer_queue <- reply_string

			} else if msg.MessageType == REP {
				//TODO!
				//msg_string := MsgString(msg)
				//msg_string := EdnConstructMsgMapS(msg)

				logmsgc(ntchan, "ReadProcessor", "REP_in", msgString)
				ntchan.REP_in <- msgString

				// x := <-ntchan.REP_in
				// vlog(ntchan, "x "+x)
			}
			//  else if msg.MessageType == SUB {

			// }

			//ntchan.Reader_processed++
			//log.Println(" ", ntchan.Reader_processed, ntchan)
		}
	}

}

//process from higher level chans into write queue
func WriteProcessor(ntchan Ntchan) {
	for {

		select {
		case msg := <-ntchan.REP_out:
			//read from REQ_out
			//log.Println("[Writeprocessor]  REP_out", msg)
			logmsgc(ntchan, "WriteProcessor", "REP_out", msg)
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.REQ_out:
			logmsgc(ntchan, "Writeprocessor", "REQ_out ", msg)
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.PUB_out:
			ntchan.Writer_queue <- msg
			//PUB?
		}
	}
}

func WriteLoop(ntchan Ntchan, d time.Duration) {
	//msg_writer_total := 0
	for {
		//log.Println("loop writer")
		//TODO!
		//

		//take from channel and write
		msg := <-ntchan.Writer_queue
		vlog(ntchan, "writeloop "+msg)
		NetWrite(ntchan, msg)
		//logmsg(ntchan.Name, "WriteLoop", msg, msg_writer_total)
		//NetworkWrite(ntchan, msg)

		time.Sleep(d)
		//msg_writer_total++
	}
}

//simple echo net
func NetConnectorSetupEcho(ntchan Ntchan) {

	//read loop
	go func() {
		d := 300 * time.Millisecond
		for {
			//vlog(ntchan, "iter ReadLoop "+ntchan.SrcName+" "+ntchan.DestName)
			msg, _ := NetMsgRead(ntchan)
			//currently can be empty or len, shoudl fix one style
			if len(msg) > 0 { //&& msg != EMPTY_MSG {
				vlog(ntchan, "ntwk read => "+msg)
				logmsgc(ntchan, ntchan.SrcName, "ReadLoop", msg)
				vlog(ntchan, "put on Reader queue "+msg)
				reply := "echo >>> " + msg
				ntchan.Reader_queue <- reply
			}

			time.Sleep(d)
		}
	}()

	//echo back
	go func() {
		for {
			msgString := <-ntchan.Reader_queue
			NetWrite(ntchan, msgString)
		}
	}()
}
