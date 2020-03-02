package ntcl

import (
	"log"
	"net"
	"time"
)

const (
	EMPTY_MSG  = "EMPTY"
	ERROR_READ = "error_read"
)

type Ntchan struct {
	//Rw   *bufio.ReadWriter
	Conn net.Conn
	Name string
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

func logmsg(name string, src string, msg string, total int) {
	log.Printf("%s [%s] ### %v  %d", name, src, msg, total)
}

//wrap connection in Ntchan
func ConnNtchan(conn net.Conn) Ntchan {
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

func ReadLoop(ntchan Ntchan) {
	//vlog("init ReadLoop")
	d := 300 * time.Millisecond
	//msg_reader_total := 0
	for {
		//read from network and put in channel
		//vlog("iter ReadLoop")
		msg, err := MsgRead(ntchan.Conn)
		if err != nil {

		}
		//vlog("ntwk read => " + msg)
		//handle cases
		//currently can be empty or len, shoudl fix one style
		if len(msg) > 0 && msg != EMPTY_MSG {
			logmsgc(ntchan.Name, "ReadLoop", msg)
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
		log.Println(msgString)
	}

}
