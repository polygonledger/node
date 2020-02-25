package ntwk

//functions relating to network stack
//using bufio.ReadWriter stream

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//network channel, wrapper around ReadWriter
type Ntchan struct {
	Rw *bufio.ReadWriter
	//chan?
}

//read a message from network
func NetworkRead(nt Ntchan) string {
	//TODO handle err
	msg, _ := nt.Rw.ReadString(DELIM)
	msg = strings.Trim(msg, string(DELIM))
	return msg
}

//continous loop of processing reads
func ReaderLoop(nt Ntchan, msg_in_chan chan Message, msg_out_chan chan Message) {
	//
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

//wrap connection in Ntchan
func ConnNtchan(conn net.Conn) Ntchan {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return Ntchan{Rw: rw}
}

func OpenNtchan(addr string) Ntchan {
	conn := OpenConn(addr)
	return ConnNtchan(conn)
}
