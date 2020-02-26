package ntwk

// functions relating to network stack

// network layer (NTL)

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
// socket we need to coordinate the reads

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//network channel
type Ntchan struct {
	Rw           *bufio.ReadWriter
	Name         string
	Reader_queue chan string
	Writer_queue chan string
}

func ReaderWriterConnector(ntchan Ntchan) {
	//func (ntchan Ntchan) ReaderWriterConnector() {

	read_loop_time := 300 * time.Millisecond
	read_time_chan := 300 * time.Millisecond
	write_loop_time := 300 * time.Millisecond
	write_processor_time := 300 * time.Millisecond

	//reader
	go ReadLoop(ntchan, read_loop_time)

	go ReadProcessor(ntchan, read_time_chan)

	//writer
	go WriteLoop(ntchan, write_loop_time)

	go WriteProducer(ntchan, write_processor_time)
}

func logmsg(name string, src string, msg string, total int) {
	log.Printf("%s [%s] ### %v  %d", name, src, msg, total)
}

//continous network reads with sleep
func ReadLoop(ntchan Ntchan, d time.Duration) {
	msg_reader_total := 0
	for {
		logmsg(ntchan.Name, "ReadLoop", " read ", msg_reader_total)
		//read from network and put in channel
		msg := NetworkReadMessage(ntchan)
		//handle cases
		if len(msg) > 0 {
			ntchan.Reader_queue <- msg
		}

		time.Sleep(d)
		msg_reader_total++
	}
}

//read from reader queue and process
func ReadProcessor(ntchan Ntchan, d time.Duration) {

	msg_reader_processed := 0
	for {
		log.Println("ReadProcessor in")
		msgString := <-ntchan.Reader_queue
		//log.Println("got msg on reader ", msg)
		if len(msgString) > 0 {
			logmsg(ntchan.Name, "ReadProcessor", msgString, msg_reader_processed)
			msg := ParseMessage(msgString)
			log.Println(msg.MessageType)

			if msg.MessageType == REQ {
				//TODO proper handler
				reply := EncodeMsgString(REP, CMD_PONG, EMPTY_DATA)
				ntchan.Writer_queue <- reply
			}

			msg_reader_processed++
		} else {
			//empty message
			logmsg(ntchan.Name, "ReadProcessor", "empty", msg_reader_processed)
		}

		//TODO! handle

		time.Sleep(d)
	}
}

func WriteLoop(ntchan Ntchan, d time.Duration) {
	msg_writer_total := 0
	for {
		//log.Println("loop writer")
		//TODO! bug both reading
		//take from channel and write
		msg := <-ntchan.Writer_queue
		logmsg(ntchan.Name, "WriteLoop", msg, msg_writer_total)

		NetworkWrite(ntchan, msg)

		time.Sleep(d)
		msg_writer_total++
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
	msg, err := nt.Rw.ReadString(DELIM)
	//log.Println("msg > ", msg)
	if err != nil {
		//issue
		//special case is empty message if client disconnects?
		if len(msg) == 0 {
			//log.Println("empty message")
			return EMPTY_MSG
		} else {
			log.Println("Failed ", err)
			//log.Println(err.)
			return ERROR_READ
		}
	}
	return msg
}

func NetworkWrite(nt Ntchan, message string) error {
	//log.Println("network -> write ", message)
	n, err := nt.Rw.WriteString(message)
	if err != nil {
		return errors.Wrap(err, "Could not write data ("+strconv.Itoa(n)+" bytes written)")
	} else {
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

//wrap connection in Ntchan
func ConnNtchan(conn net.Conn, name string) Ntchan {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	return Ntchan{Rw: rw, Name: name, Reader_queue: make(chan string), Writer_queue: make(chan string)}
}

func OpenNtchan(addr string) Ntchan {
	conn := OpenConn(addr)
	name := addr
	return ConnNtchan(conn, name)
}
