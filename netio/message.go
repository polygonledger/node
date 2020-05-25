package netio

import (
	"encoding/json"
)

//messages. currently uses json and some hacks to make it somewhat flexible mechanism
//needs to be ported to edn for which the basics exist

//Message Types
//Message Types exist in a context of communication
//i.e. is this message an intial request or a reply or part of a stream of events etc
//Request <--> Reply

type Message struct {
	//type of message i.e. the communications protocol
	MessageType string
	//Specific message command
	Command string
	//any data, can be empty. gets interpreted downstream to other structs
	Data json.RawMessage
	//timestamp
	Layer string
}

const (
	REQ = "REQ"
	REP = "REP"
	PUB = "PUB"
	SUB = "SUB"
	//HANDSHAKE
)

//TODO define "namespace/grammar/protocol file"
const (
	CMD_ACCOUNTS       = "ACCOUNTS"
	CMD_BALANCE        = "BALANCE" //get balance of account
	CMD_BLOCKHEIGHT    = "BLOCKHEIGHT"
	CMD_FAUCET         = "FAUCET"
	CMD_GETTXPOOL      = "GETTXPOOL"
	CMD_GETPEERS       = "GETPEERS"
	CMD_GETBLOCKS      = "GETBLOCKS"
	CMD_HEARTBEAT      = "HEARTBEAT"
	EMPTY_DATA         = "EMPTY"
	CMD_LOGOFF         = "LOGOFF"
	CMD_NUMCONN        = "NUMCONN"     //number of connected
	CMD_NUMACCOUNTS    = "NUMACCOUNTS" //number of accounts
	CMD_STATUS         = "STATUS"
	CMD_SUB            = "SUBTO"
	CMD_SUBUN          = "SUBUN" //unsuscribe
	CMD_PING           = "PING"
	CMD_PONG           = "PONG"
	CMD_RANDOM_ACCOUNT = "RANACC" //get some random account
	CMD_TX             = "TX"     //send transaction
	//app layer. this is a hack right now
	CMD_CHAT         = "CHAT"
	CMD_ERROR        = "ERROR"
	CMD_REGISTERNAME = "REGISTERNAME"
	//CMD_HANDSHAKE_HELLO  = "HELLO"
	//CMD_HANDSHAKE_STABLE = "STABLE"
	//CMD_LOGIN            = "LOGIN"
	//GET

)

var CMDS = []string{
	CMD_PING,
	CMD_PONG,
	CMD_NUMACCOUNTS,
	CMD_ACCOUNTS,
	CMD_BALANCE,
	CMD_GETPEERS,
	CMD_BLOCKHEIGHT,
	CMD_FAUCET,
	CMD_GETTXPOOL,
	CMD_GETBLOCKS,
	CMD_RANDOM_ACCOUNT,
	CMD_STATUS,
	CMD_NUMCONN,
	CMD_SUB,
	CMD_SUBUN,
	CMD_LOGOFF,
	CMD_HEARTBEAT,
	CMD_TX,
	CMD_CHAT,
	CMD_ERROR,
	CMD_REGISTERNAME,
	// CMD_FAUCET,

}

// func (m MsgType) String() string {
// 	return [...]string{"REQ", "REP"}
// }

//generic message

func IsValidMsgType(msgType string) bool {
	switch msgType {
	case
		REQ,
		REP:
		return true
	}
	return false
}

func RequestMessage() Message {
	return Message{MessageType: REQ}
}

func ReplyMessage() Message {
	return Message{MessageType: REP}
}

func ConstructMsg(msgType string, cmd string, data string) Message {
	m := Message{MessageType: msgType, Command: cmd, Data: []byte(data)}
	return m
}

func ConstructMsgBytes(msgType string, cmd string, data []byte) Message {
	m := Message{MessageType: msgType, Command: cmd, Data: data}
	return m
}

func validCMD(cmd string) bool {
	for _, a := range CMDS {
		if a == cmd {
			return true
		}
	}
	return false
}

func FromJSON(msg_string string) Message {
	var msgu Message
	json.Unmarshal([]byte(msg_string), &msgu)
	return msgu
}
