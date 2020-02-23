package ntwk

//protocol layer
//protocol operates on messages and channels not bytestreams

//TODO
//handshake message
//heartbeat

import (
	"bufio"
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
	fmt.Println(msg)
	return msg
}

func EncodeReply(resp string) string {
	//TODO header missing
	msg := EncodeMsgString(REP, resp, "")
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
	//delim := "\n"
	msg := cmd + string(DELIM)
	return msg
}

//generic request<->reply
func RequestReplyChan(request_string string, msg_in_chan chan Message, msg_out_chan chan Message) Message {
	request := ParseMessage(request_string)
	msg_in_chan <- request
	resp := <-msg_out_chan
	return resp
}

//continous loop of processing requests
func RequestLoop(rw *bufio.ReadWriter, msg_in_chan chan Message, msg_out_chan chan Message) {
	for {
		//take from channel and perform request
		request := <-msg_in_chan
		fmt.Println("request ", request)
		resp := RequestReplyMsg(rw, request)
		fmt.Println("resp ", resp)
		msg_out_chan <- resp
	}
}

//request reply of messages
func RequestReplyMsg(rw *bufio.ReadWriter, req_msg Message) Message {
	req_msg_string := MsgString(req_msg)
	resp_msg_string := RequestReplyString(rw, req_msg_string)
	resp_msg := ParseMessage(resp_msg_string)
	return resp_msg
}

//request reply on network level
func RequestReplyString(rw *bufio.ReadWriter, req_msg_string string) string {
	NetworkWrite(rw, req_msg_string)
	resp_msg := NetworkRead(rw)
	return resp_msg
}

func ReplyNetwork(rw *bufio.ReadWriter, resp Message) {
	//rep_msg := EncodeReply(resp)
	resp_string := MsgString(resp)
	NetworkWrite(rw, resp_string)
}

func ReadMessage(rw *bufio.ReadWriter) Message {
	var msg Message
	msgString := NetworkReadMessage(rw)
	if msgString == EMPTY_MSG {
		return EmptyMsg()
	}
	msg = ParseMessage(msgString)
	return msg
}
