package ntwk

//protocol layer
//protocol operates on messages and channels not bytestreams

//TODO
//handshake message
//heartbeat

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
func RequestLoop(ntchan Ntchan, msg_in_chan chan Message, msg_out_chan chan Message) {
	for {
		//take from channel and perform request
		request := <-msg_in_chan
		//fmt.Println("request ", request)
		resp := RequestReplyMsg(ntchan, request)
		//fmt.Println("resp ", resp, resp.MessageType)
		//BUG MSG can be of type pub
		if resp.MessageType == REP {
			msg_out_chan <- resp
		}
	}
}

//request reply of messages
func RequestReplyMsg(ntchan Ntchan, req_msg Message) Message {
	req_msg_string := MsgString(req_msg)
	//log.Println("> ", req_msg_string)
	resp_msg_string := RequestReplyString(ntchan, req_msg_string)
	resp_msg := ParseMessage(resp_msg_string)
	return resp_msg
}

//request reply on network level
func RequestReplyString(ntchan Ntchan, req_msg_string string) string {
	NetworkWrite(ntchan, req_msg_string)
	resp_msg := NetworkRead(ntchan)
	return resp_msg
}

func ReplyNetwork(ntchan Ntchan, resp Message) {
	//rep_msg := EncodeReply(resp)
	resp_string := MsgString(resp)
	NetworkWrite(ntchan, resp_string)
}

func ReadMessage(ntchan Ntchan) Message {
	var msg Message
	msgString := NetworkReadMessage(ntchan)
	if msgString == EMPTY_MSG {
		return EmptyMsg()
	}
	msg = ParseMessage(msgString)
	return msg
}

func PubNetwork(ntchan Ntchan, resp Message) {
	//rep_msg := EncodeReply(resp)
	//resp_string := EncodePub(resp)
	resp_string := EncodePub("TEST", "")
	NetworkWrite(ntchan, resp_string)
}
