package ntcl

import (
	"log"
	"net"
	"time"

	"github.com/polygonledger/node/ntwk"
)

const (
	EMPTY_MSG  = "EMPTY"
	ERROR_READ = "error_read"
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

	// Reader_processed int
	// Writer_processed int
}

func logmsgd(src string, msg string) {
	log.Printf("[%s] ### %v", src, msg)
}

func logmsgc(name string, src string, msg string) {
	log.Printf("%s [%s] ### %v", name, src, msg)
}

//wrap connection in Ntchan
func ConnNtchan(conn net.Conn, SrcName string, DestName string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	// ntchan.Reader_processed = 0
	// ntchan.Writer_processed = 0
	ntchan.Conn = conn
	ntchan.SrcName = SrcName
	ntchan.DestName = DestName

	return ntchan
}

func ConnNtchanStub(name string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	//ntchan.Reader_processed = 0
	//ntchan.Writer_processed = 0

	return ntchan
}

func vlog(s string) {
	log.Println(s)
}

func ReadLoop(ntchan Ntchan) {
	vlog("init ReadLoop " + ntchan.SrcName + " " + ntchan.DestName)
	d := 300 * time.Millisecond
	//msg_reader_total := 0
	for {
		//read from network and put in channel
		vlog("iter ReadLoop " + ntchan.SrcName + " " + ntchan.DestName)
		msg, err := MsgRead(ntchan)
		if err != nil {

		}
		//handle cases
		//currently can be empty or len, shoudl fix one style
		if len(msg) > 0 { //&& msg != EMPTY_MSG {
			vlog("ntwk read => " + msg)
			logmsgc(ntchan.SrcName, "ReadLoop", msg)
			log.Println("put ", msg)
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
		msgString := <-ntchan.Reader_queue
		logmsgd("ReadProcessor", msgString)

		if len(msgString) > 0 {
			logmsgc(ntchan.SrcName, "ReadProcessor", msgString) //, ntchan.Reader_processed)
			//TODO try catch

			msg := ntwk.ParseMessage(msgString)

			if msg.MessageType == ntwk.REQ {

				msg_string := ntwk.MsgString(msg)
				logmsgd("ReadProcessor", "REQ_in")

				//TODO!
				ntchan.REQ_in <- msg_string
				// reply_string := "echo:" + msg_string
				// log.Println(">> ", reply_string)
				// ntchan.Writer_queue <- reply_string

			} else if msg.MessageType == ntwk.REP {
				//TODO!
				//msg_string := MsgString(msg)
				msg_string := ntwk.MsgString(msg)
				logmsgd("ReadProcessor", "REP_in")
				ntchan.REP_in <- msg_string

				x := <-ntchan.REP_in
				log.Println("x ", x)
			}

			//ntchan.Reader_processed++
			//log.Println(" ", ntchan.Reader_processed, ntchan)
		}
	}

}

func WriteLoop(ntchan Ntchan, d time.Duration) {
	//msg_writer_total := 0
	for {
		//log.Println("loop writer")
		//TODO!
		//

		//msg_string := <-ntchan.REQ_out

		//take from channel and write
		msg := <-ntchan.Writer_queue
		log.Println("writeloop ", msg)
		NtwkWrite(ntchan, msg)
		//logmsg(ntchan.Name, "WriteLoop", msg, msg_writer_total)
		//NetworkWrite(ntchan, msg)

		time.Sleep(d)
		//msg_writer_total++
	}
}
