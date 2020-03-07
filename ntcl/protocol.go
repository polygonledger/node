package ntcl

import (
	"fmt"
	"strings"
)

const (
	// Port is the port number that the server listens to.
	//TODO move to message type
	Genesis_Address string = "P0614579c42f2"
	DELIM           byte   = '|'
	DELIM_HEAD      byte   = '#'
	EMPTY_MSG              = "EMPTY"
	ERROR_READ             = "error_read"
)

func trace(msg string) {
	fmt.Println(msg)
}

//parse message string into a message struct
//delimiters of two kind:
//* DELIM for delimiting the entire message
//* DELIM_HEAD
//currently we employ delimiters instead of byte encoding, so the size of messages is unlimited
//can however easily fix by adding size to header and reject messages larger than maximum size
func ParseMessage(msgString string) Message {
	msgString = strings.Trim(msgString, string(DELIM))
	s := strings.Split(msgString, string(DELIM_HEAD))

	//ERROR handling of malformed messages

	var msg Message
	msg.MessageType = s[0]
	msg.Command = s[1]
	dataJson := s[2] //data can empty but still we expect the delim to be there

	msg.Data = []byte(dataJson)
	//trace(msg)
	return msg
}

func EncodeReply(resp string) string {
	//TODO header missing
	msg := EncodeMsgString(REP, resp, "")
	return msg
}

func EncodeRequest(req_string string) string {
	//TODO header missing
	msg := EncodeMsgString(REQ, req_string, "")
	return msg
}

func EncodePub(resp string, name string) string {
	//TODO header missing
	msg := EncodeMsgString(PUB, resp, name)
	return msg
}

func EncodeHeartbeat(name string) string {
	//TODO time
	msg := EncodePub(CMD_HEARTBEAT, name)
	return msg
}

func EncodeMessageTx(txJson []byte) string {
	//emptyData := ""
	msgCmd := "TX"
	//TODO check
	msg := EncodeMsgString(REQ, msgCmd, string(txJson))
	return msg
}

func ConstructMessage(cmd string) string {
	msg := cmd + string(DELIM)
	return msg
}
