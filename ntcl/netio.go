package ntcl

// netio contains functions relating to network stack

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/polygonledger/node/ntcl"
)

const DELIM = '|'

func EncodeMsg(content string) string {
	return content + string(DELIM)
}

//TODO! replace old in ntwk package
func NtwkWrite(ntchan Ntchan, content string) (int, error) {
	//READLINE uses \n
	NEWLINE := '\n'
	//respContent := fmt.Sprintf("%s%c%c", content, DELIM, NEWLINE)
	respContent := fmt.Sprintf("%s%c", content, NEWLINE)
	//log.Println("write > ", content, respContent)
	writer := bufio.NewWriter(ntchan.Conn)
	n, err := writer.WriteString(respContent)
	if err == nil {
		err = writer.Flush()
	}
	s := fmt.Sprintf("bytes written", n, " ", ntchan.SrcName, ntchan.DestName)
	vlog(s)
	return n, err
}

func NtwkRead(ntchan Ntchan, delim byte) (string, error) {
	//log.Println("NtwkRead ", ntchan.SrcName, ntchan.DestName)
	reader := bufio.NewReader(ntchan.Conn)
	var buffer bytes.Buffer
	for {
		//READLINE uses \n
		ba, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		buffer.Write(ba)
		if !isPrefix {
			break
		}
	}
	return buffer.String(), nil
}

func MsgRead(ntchan Ntchan) (string, error) {
	msg_string, err := NtwkRead(ntchan, DELIM)
	msg_string = strings.Trim(msg_string, string(DELIM))
	return msg_string, err
}

func initClient() Ntchan {
	addr := ":" + strconv.Itoa(node_port)
	log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("cant connect to ", addr)
		//return
	}

	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr)

	go ReadLoop(ntchan)
	go ReadProcessor(ntchan)
	go WriteProcessor(ntchan, 100*time.Millisecond)
	go WriteLoop(ntchan, 300*time.Millisecond)
	return ntchan

}
