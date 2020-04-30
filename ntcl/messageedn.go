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

func EncodeMsgMapS(msg Message) string {
	m := map[string]string{msg.MessageType: msg.Command}
	msgstring := parser.MakeMap(m)
	return msgstring
}

func EncodeMsgMapData(msgType string, cmd string, data string) string {
	m := map[string]string{msgType: cmd, "data": data}
	msg := parser.MakeMap(m)
	return msg
}


func ParseMessageMap(msgString string) Message {
	//msgString = strings.Trim(msgString, string(DELIM))
	//s := strings.Split(msgString, string(DELIM_HEAD))
	//ERROR handling of malformed messages

	fmt.Println("msgString ",msgString)
	v, k := parser.ReadMap(msgString)
	// fmt.Println(v)
	// fmt.Println(k)

	var msg Message
	msg.MessageType = k[0]
	msg.Command = v[0]
	// dataJson := s[2] //data can empty but still we expect the delim to be there

	// msg.Data = []byte(dataJson)
	// //trace(msg)
	return msg
}

func ParseMessageMapData(msgString string) Message {
	//msgString = strings.Trim(msgString, string(DELIM))
	//s := strings.Split(msgString, string(DELIM_HEAD))
	//ERROR handling of malformed messages

	fmt.Println("msgString ",msgString)
	v, k := parser.ReadMap(msgString)
	
	var msg Message
	msg.MessageType = k[0]
	msg.Command = v[0]
	msg.Data = []byte(v[1])
	// dataJson := s[2] //data can empty but still we expect the delim to be there

	// msg.Data = []byte(dataJson)
	// //trace(msg)
	return msg
}
