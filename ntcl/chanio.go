package ntcl

// network layer (NTCL)

// NTCL -> semantics of channels
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
	Conn     net.Conn
	SrcName  string
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

	PUB_time_out  chan string
	SUB_time_out  chan string
	PUB_time_quit chan int
	verbose       bool
	// SUB_request_out   chan string
	// SUB_request_in    chan string
	// UNSUB_request_out chan string
	// UNSUB_request_in  chan string

	// Reader_processed int
	// Writer_processed int
}

func vlog(ntchan Ntchan, s string) {
	verbose := ntchan.verbose
	if verbose {
		log.Println(s)
	}
}

func logmsgd(ntchan Ntchan, src string, msg string) {
	s := fmt.Sprintf("[%s] ### %v", src, msg)
	vlog(ntchan, s)
}

func logmsgc(ntchan Ntchan, name string, src string, msg string) {
	s := fmt.Sprintf("%s [%s] ### %v", name, src, msg)
	vlog(ntchan, s)
}

func ConnNtchan(conn net.Conn, SrcName string, DestName string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.PUB_time_out = make(chan string)
	ntchan.PUB_time_quit = make(chan int)
	// ntchan.Reader_processed = 0
	// ntchan.Writer_processed = 0
	ntchan.Conn = conn
	ntchan.SrcName = SrcName
	ntchan.DestName = DestName
	ntchan.verbose = true

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
	ntchan.PUB_time_out = make(chan string)
	ntchan.PUB_time_quit = make(chan int)
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
		msg, err := MsgRead(ntchan)
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

//read from reader queue and process by forwarding to right channel
func ReadProcessor(ntchan Ntchan) {

	for {
		log.Println("ReadProcessor loop")
		msgString := <-ntchan.Reader_queue
		logmsgd(ntchan, "ReadProcessor", msgString)

		if len(msgString) > 0 {
			logmsgc(ntchan, ntchan.SrcName, "ReadProcessor", msgString) //, ntchan.Reader_processed)
			//TODO try catch

			msg := ParseMessage(msgString)

			if msg.MessageType == REQ {

				msg_string := MsgString(msg)
				logmsgd(ntchan, "ReadProcessor", "REQ_in")

				//TODO!
				ntchan.REQ_in <- msg_string
				// reply_string := "echo:" + msg_string
				// log.Println(">> ", reply_string)
				// ntchan.Writer_queue <- reply_string

			} else if msg.MessageType == REP {
				//TODO!
				//msg_string := MsgString(msg)
				msg_string := MsgString(msg)
				logmsgd(ntchan, "ReadProcessor", "REP_in")
				ntchan.REP_in <- msg_string

				// x := <-ntchan.REP_in
				// vlog(ntchan, "x "+x)
			}

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

func PublishTime(ntchan Ntchan) {
	timeFormat := "2006-01-02T15:04:05"
	limiter := time.Tick(1000 * time.Millisecond)
	pubcount := 0
	log.Println("PublishTime")

	for {
		t := time.Now()
		tf := t.Format(timeFormat)
		vlog(ntchan, "pub "+tf)
		ntchan.PUB_time_out <- tf
		<-limiter
		pubcount++
	}

}

//publication to writer queue. requires quit channel
func PubWriterLoop(ntchan Ntchan) {

	for {
		select {
		case msg := <-ntchan.PUB_time_out:
			vlog(ntchan, "sub "+msg)
			ntchan.Writer_queue <- msg
		case <-ntchan.PUB_time_quit:
			fmt.Println("stop pub")
			return
			// default:
			// 	fmt.Println("no message received")
		}
		time.Sleep(50 * time.Millisecond)

	}

}
