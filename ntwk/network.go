package ntwk

//functions relating to network stack
//using bufio.ReadWriter stream

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//read a message from network
func NetworkRead(rw *bufio.ReadWriter) string {
	//TODO handle err
	msg, _ := rw.ReadString(DELIM)
	msg = strings.Trim(msg, string(DELIM))
	return msg
}

//continous loop of processing reads
func ReaderLoop(rw *bufio.ReadWriter, msg_in_chan chan Message, msg_out_chan chan Message) {
	//
}

//given a sream read from it
//TODO proper error handling
func NetworkReadMessage(rw *bufio.ReadWriter) string {
	msg, err := rw.ReadString(DELIM)
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

func NetworkWrite(rw *bufio.ReadWriter, message string) error {
	n, err := rw.WriteString(message)
	if err != nil {
		return errors.Wrap(err, "Could not write data ("+strconv.Itoa(n)+" bytes written)")
	} else {
		//TODO log trace
		//log.Println(strconv.Itoa(n) + " bytes written")
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}
