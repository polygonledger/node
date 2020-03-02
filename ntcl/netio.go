package ntcl

// netio contains functions relating to network stack

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

const DELIM = '|'

func EncodeMsg(content string) string {
	return content + string(DELIM)
}

//TODO! factor and replace old
func NtwkWrite(conn net.Conn, content string) (int, error) {
	//READLINE uses \n
	NEWLINE := '\n'
	respContent := fmt.Sprintf("%s%c%c", content, DELIM, NEWLINE)
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(respContent)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func NtwkRead(conn net.Conn, delim byte) (string, error) {
	reader := bufio.NewReader(conn)
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

func MsgRead(conn net.Conn) (string, error) {
	return NtwkRead(conn, DELIM)
}
