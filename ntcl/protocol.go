package ntcl

import (
	"fmt"
)

const (
	// Port is the port number that the server listens to.
	//TODO move to message type
	Genesis_Address string = "P0614579c42f2"
	EMPTY_MSG              = "EMPTY"
	ERROR_READ             = "error_read"
)

func trace(msg string) {
	fmt.Println(msg)
}

//parse message string into a message struct

func EncodeReply(resp string) string {
	//TODO header missing
	msg := EncodeMsgMap(REP, resp)
	return msg
}

func EncodeRequest(req_string string) string {
	//TODO header missing
	msg := EncodeMsgMap(REQ, req_string)
	return msg
}

func EncodePub(resp string, name string) string {
	//TODO header missing
	msg := EncodeMsgMap(PUB, resp)
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
	msg := EncodeMsgMapData(REQ, msgCmd, string(txJson))
	return msg
}
