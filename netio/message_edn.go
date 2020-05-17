package netio

import "github.com/polygonledger/node/parser"

//edn based messages

func EdnConstructMsgMap(msgType string, cmd string) string {
	m := map[string]string{msgType: cmd}
	msg := parser.MakeMap(m)
	return msg
}

func EdnConstructMsgMapS(msg Message) string {
	m := map[string]string{msg.MessageType: msg.Command}
	msgstring := parser.MakeMap(m)
	return msgstring
}

func EdnConstructMsgMapData(msgType string, cmd string, data string) string {
	m := map[string]string{msgType: cmd, "data": data}
	msg := parser.MakeMap(m)
	return msg
}

func EdnParseMessageMap(msgString string) Message {
	//msgString = strings.Trim(msgString, string(DELIM))
	//s := strings.Split(msgString, string(DELIM_HEAD))
	//ERROR handling of malformed messages
	//fmt.Println(msgString)
	v, k := parser.ReadMap(msgString)

	var msg Message
	msg.MessageType = k[0]
	msg.Command = v[0]
	// dataJson := s[2] //data can empty but still we expect the delim to be there

	// msg.Data = []byte(dataJson)
	// //trace(msg)
	return msg
}

func EdnParseMessageMapData(msgString string) Message {

	v, k := parser.ReadMap(msgString)
	// fmt.Println("values ", v)
	// fmt.Println("keys ", k)

	var msg Message
	msg.MessageType = k[0]
	msg.Command = v[0]

	msg.Data = []byte(v[1])
	return msg
}
