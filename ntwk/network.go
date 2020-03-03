package ntwk

// functions relating to network stack

// network layer (NTL)

// NTL -> semantics of channels
// TCP/IP -> golang net

// golang has the native net package. as TCP only deals with byte streams we need some form
// to delinate distinct messages to implement the equivalent of actors. since we have
// channels as the major building block for the network we wrap the bufio readwriter in
// defined set of channels with equivalent set of messages.

// the struct ntchan is the struct wraps the native readwriter with reads and write queues
// as channels.

// network reads happen in distinct units of messages which are delimited by the DELIM byte
// messages have types to indicate the flow of message direction and commands
// an open question is use of priorities, timing etc.

// the P2P network or any network connection really has different behaviour based on the
// types of messages going through it. a request-reply for example will have a single read
// and single write in order, publish-subscribe will  push messages from producers to
// consumers, etc.

// since we always have only one single two-way channel available as we are on a single
// socket we need to coordinate the reads and writes. the network is essentialy a scarce
// resource and depending on the context and semantics messages will be sent/received in
// different style

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//network channel, abstraction of a network connection
type Ntchan struct {
	Rw *bufio.ReadWriter
	//Conn *net.conn
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

	Reader_processed int
	Writer_processed int
}

// --- NTL layer ---

func vlog(s string) {
	log.Println(s)
}

// main setup of all read and write processes for a single connection
func ReaderWriterConnector(ntchan Ntchan) {
	vlog("ReaderWriterConnector")
	//func (ntchan Ntchan) ReaderWriterConnector() {

	//timers
	read_loop_time := 800 * time.Millisecond
	read_time_chan := 300 * time.Millisecond
	write_loop_time := 300 * time.Millisecond
	//write_processor_time := 300 * time.Millisecond

	//any coordination between reader and writer

	//init reader
	//reads from the actual "physical" network
	go ReadLoop(ntchan, read_loop_time)

	//process of reads
	go ReadProcessor(ntchan, read_time_chan)

	//init writer
	//write to network whatever is the reader queue
	go WriteLoop(ntchan, write_loop_time)

	go Writeprocessor(ntchan, 200*time.Millisecond)

	//REQ processor
	//go Reqprocessor1(ntchan)

	// go func() {
	// 	ntchan.REQ_in <- "test"
	// 	//xchan <- "test"
	// }()

	//TODO
	//go WriteProducer(ntchan, write_processor_time)
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

//continous network reads with sleep
func ReadLoop(ntchan Ntchan, d time.Duration) {
	vlog("init ReadLoop")
	msg_reader_total := 0
	for {
		//read from network and put in channel
		//vlog("iter ReadLoop")
		msg := NetworkReadMessage(ntchan)
		//vlog("ntwk read => " + msg)
		//handle cases
		//currently can be empty or len, shoudl fix one style
		if len(msg) > 0 && msg != EMPTY_MSG {
			logmsg(ntchan.Name, "ReadLoop", msg, msg_reader_total)
			//put in the queue to process
			ntchan.Reader_queue <- msg
		}

		time.Sleep(d)
		//fix: need ntchan to be a pointer
		//msg_reader_total++
	}
}

// func RequestProcessor(ntchan *Ntchan, d time.Duration) {
// 	log.Println("init RequestProcessor ")
// 	for {
// 		request := <-ntchan.Req_chan
// 		log.Println("request ", request)

// 		//reply_string := EncodeMsgString(REP, CMD_PONG, EMPTY_DATA)
// 		reply := EncodeMsg(REP, CMD_PONG, EMPTY_DATA)
// 		//TODO reply goes through req_chan only
// 		ntchan.Rep_chan <- reply
// 		//ntchan.Writer_queue <- reply_string

// 	}
// }

//process replies
func ReplyProcessor(ntchan *Ntchan, d time.Duration) {
	vlog("init ReplyProcessor ")
	for {
		reply_string := <-ntchan.REP_out
		//log.Println("reply ", reply)
		//reply_string := MsgString(reply)
		ntchan.Writer_queue <- reply_string

	}
}

//read from reader queue and process by forwarding to right channel
func ReadProcessor(ntchan Ntchan, d time.Duration) {

	vlog("init ReadProcessor")

	//loop and basic fanout based on message type
	//can optimize for performance here based on channel select
	for {
		msgString := <-ntchan.Reader_queue
		ntchan.Reader_processed++
		//log.Println("got msg on reader ", msg)

		if len(msgString) > 0 {
			logmsg(ntchan.Name, "ReadProcessor", msgString, ntchan.Reader_processed)
			//TODO try catch
			msg := ParseMessage(msgString)

			if msg.MessageType == REQ {
				//TODO proper handler
				//log.Println("req ", msg.Command)
				//logmsg(ntchan.Name, "ReadProcessor", msg.Command, 0)

				msg_string := MsgString(msg)
				logmsgd("ReadProcessor", "REQ_in")
				ntchan.REQ_in <- msg_string

				//ntchan.Writer_queue <- reply_string

			} else if msg.MessageType == REP {
				//TODO!
				//msg_string := MsgString(msg)
				msg_string := MsgString(msg)
				logmsgd("ReadProcessor", "REP_in")
				ntchan.REP_in <- msg_string

				x := <-ntchan.REP_in
				log.Println("x ", x)
			}

			//ntchan.Reader_processed++
			//log.Println(" ", ntchan.Reader_processed, ntchan)
		} else {
			//empty message
			logmsg(ntchan.Name, "ReadProcessor", "empty", ntchan.Reader_processed)
		}

		//TODO handle

		time.Sleep(d)
	}
}

func WriteLoop(ntchan Ntchan, d time.Duration) {
	msg_writer_total := 0
	for {
		//log.Println("loop writer")
		//TODO!
		//

		//msg_string := <-ntchan.REQ_out

		//take from channel and write
		msg := <-ntchan.Writer_queue
		logmsg(ntchan.Name, "WriteLoop", msg, msg_writer_total)

		NetworkWrite(ntchan, msg)

		time.Sleep(d)
		msg_writer_total++
	}
}

func Writeprocessor(ntchan Ntchan, d time.Duration) {

	vlog("init Writeprocessor")

	//TODO!
	//put on Writer_queue
	for {
		vlog("loop Writeprocessor")

		//selectively read on write outputs
		select {
		case msg := <-ntchan.REP_out:
			//read from REQ_out
			//log.Println("[Writeprocessor]  REP_out", msg)
			logmsgc("WriteProcessor", "REP_out", msg)
			ntchan.Writer_queue <- msg

		case msg := <-ntchan.REQ_out:
			logmsgc("Writeprocessor", "REQ_out ", msg)
			ntchan.Writer_queue <- msg

			//PUB?
		}
	}
}

func WriteProducer(ntchan Ntchan, d time.Duration) {
	msg_write_processed := 0
	for {
		//TODO gather produced writes from other channels
		msg := "test"

		ntchan.Writer_queue <- msg
		//log.Println("got msg on reader ", msg)
		if len(msg) > 0 {
			logmsg(ntchan.Name, "WriteProducer", msg, msg_write_processed)
			msg_write_processed++
		} else {
			//empty message
		}

		//TODO! handle

		time.Sleep(d)
	}
}

// --- underlying network stack calls ---

//read a message from network
func NetworkRead(nt Ntchan) string {
	//TODO handle err
	msg, _ := nt.Rw.ReadString(DELIM)
	msg = strings.Trim(msg, string(DELIM))
	return msg
}

//given a sream read from it
//TODO proper error handling
func NetworkReadMessage(nt Ntchan) string {
	vlog("NetworkReadMessage")
	msg, err := nt.Rw.ReadString(DELIM)
	if err != nil {
		if !(err == io.EOF) {
			logmsgd("NetworkReadMessage err", string(err.Error()))
		}

		//issue
		//special case is empty message if client disconnects?

		// else {
		// 	log.Println("Failed ", err)
		// 	//log.Println(err.)
		// 	return ERROR_READ
		// }
	}
	if len(msg) > 0 {
		logmsgd("NetworkReadMessage", msg)
		//vlog("msg > " + msg)
		//log.Println("empty message")
	} else {
		return EMPTY_MSG
	}

	return msg
}

func NetworkWrite(nt Ntchan, message string) error {
	vlog("#network# -> write " + message)
	n, err := nt.Rw.WriteString(message)
	if err != nil {
		return errors.Wrap(err, "Could not write data ("+strconv.Itoa(n)+" bytes written)")
	} else {
		log.Println("??")
		//TODO log trace
		//log.Println(strconv.Itoa(n) + " bytes written")
	}
	err = nt.Rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

func OpenConn(addr string) net.Conn {
	// Dial the remote process
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	return conn
}

// connects to a TCP Address
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//return nil, errors.Wrap(err, "Dialing "+addr+" failed")
		log.Println("error ", err)
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
	//tmp
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	//TODO! return conn, not rw
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func OpenOut(ip string, Port int) (*bufio.ReadWriter, error) {
	addr := ip + ":" + strconv.Itoa(Port)
	log.Println("> open out address ", addr)
	rw, err := Open(addr)
	return rw, err
}

func OpenNtchanOut(ip string, Port int) Ntchan {
	fulladdr := ip + ":" + strconv.Itoa(Port)
	return OpenNtchan(fulladdr)
}

func ConnNtchanStub(name string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.Reader_processed = 0
	ntchan.Writer_processed = 0

	return ntchan
}

//wrap connection in Ntchan
func ConnNtchan(conn net.Conn, name string) Ntchan {
	var ntchan Ntchan
	ntchan.Reader_queue = make(chan string)
	ntchan.Writer_queue = make(chan string)
	ntchan.REQ_in = make(chan string)
	ntchan.REP_in = make(chan string)
	ntchan.REP_out = make(chan string)
	ntchan.REQ_out = make(chan string)
	ntchan.Reader_processed = 0
	ntchan.Writer_processed = 0

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	ntchan.Rw = rw

	return ntchan
}

func OpenNtchan(addr string) Ntchan {
	conn := OpenConn(addr)
	name := addr
	return ConnNtchan(conn, name)
}

//old
func Reqprocessor1(ntchan Ntchan) {
	x := <-ntchan.REQ_in
	//x := <-xchan
	logmsgd("REQ processor ", x)

	//reply_string := "reply"
	reply := EncodeMsg(REP, CMD_PONG, EMPTY_DATA)
	reply_string := MsgString(reply)

	ntchan.Writer_queue <- reply_string
}
