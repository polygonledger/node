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

	"github.com/pkg/errors"
)

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
	vlog(ntchan, s)
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
	// 	//timeoutDuration := 5 * time.Second
	// 	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	return conn
}

func OpenNtchanOut(src string, ip string, Port int) Ntchan {
	fulladdr := ip + ":" + strconv.Itoa(Port)
	return OpenNtchan(src, fulladdr)
}

func OpenNtchan(src string, dest string) Ntchan {
	conn := OpenConn(dest)
	//name := addr
	return ConnNtchan(conn, src, dest)
}
