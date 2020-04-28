package ntcl

import (
	"fmt"

	"github.com/polygonledger/node/parser"
)

func EncodeMsgMap(msgType string, cmd string) string {
	m := map[string]string{msgType: cmd}
	msg := parser.MakeMap(m)
	return msg
}

func ParseMessageMap(msgString string) Message {
	//msgString = strings.Trim(msgString, string(DELIM))
	//s := strings.Split(msgString, string(DELIM_HEAD))
	//ERROR handling of malformed messages

	v, k := parser.ReadMap1(msgString)
	fmt.Println(v)
	fmt.Println(k)

	var msg Message
	// msg.MessageType = s[0]
	// msg.Command = s[1]
	// dataJson := s[2] //data can empty but still we expect the delim to be there

	// msg.Data = []byte(dataJson)
	// //trace(msg)
	return msg
}
