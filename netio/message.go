package netio

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
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
}

type MessageJSON struct {
	//type of message i.e. the communications protocol
	MessageType string `json:"messagetype"`
	//Specific message command
	Command string `json:"command"`
	//any data, can be empty. gets interpreted downstream to other structs
	Data json.RawMessage `json:"data,omitempty"`
	//timestamp
}

//marshal to json, check command
func NewJSONMessage(m Message) (MessageJSON, error) {
	//fmt.Println("NewJSONMessage")
	valid := validCMD(m.Command)
	if valid {
		//fmt.Println("?? ", m.Data)
		//		if m.Data != nil {
		return MessageJSON{
			m.MessageType,
			m.Command,
			m.Data,
		}, nil

	} else {
		fmt.Println("not valid cmd")
	}
	return MessageJSON{}, errors.New("not valid cmd")
}

func ToJSONMessage(m Message) string {
	jsonmsgtype, _ := NewJSONMessage(m)
	jsonmsg, _ := json.Marshal(jsonmsgtype)
	return string(jsonmsg)
}

const (
	REQ = "REQ"
	REP = "REP"
	PUB = "PUB"
	//SUB
	//HANDSHAKE
)

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
	// CMD_FAUCET,

}

// func (m MsgType) String() string {
// 	return [...]string{"REQ", "REP"}
// }

//generic message

type TimeMessage struct {
	Timestamp time.Time
}

func EmptyMsg() Message {
	return Message{}
}

func IsValidMsgType(msgType string) bool {
	switch msgType {
	case
		REQ,
		REP:
		return true
	}
	return false
}

func IsValidCmd(cmd string) bool {
	switch cmd {
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

// type MessageTx struct {
// 	MessageType string
// 	Command     string
// 	Tx          block.Tx
// }

//////////////////////

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

// func (m *Message) UnmarshalJSON(data []byte) error {
// 	var res Message

// 	// Assign value from json to Go struct
// 	m.MessageType = res.MessageType
// 	if validCMD(res.Command) {
// 		fmt.Println("valid command")
// 		m.Command = res.Command
// 	} else {
// 		fmt.Println("not valid command ", res.Command)
// 	}

// 	//test if Command is in accepted list

// 	// if err != nil {
// 	// 	return err
// 	// }

// 	//p.Password = password

// 	// Parse millisecond to string
// 	//p.CreatedAt = time.Unix(0, res.CreatedAt*int64(time.Millisecond)).UTC().Format(time.RFC3339)

// 	return nil
// }
