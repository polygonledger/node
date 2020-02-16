package net

import (
	"fmt"
	"time"

	block "github.com/polygonledger/node/block"
)

//any message is defined through the delimters, no size restrictions as of now
//MSG TYPE # CMD # DATA |

//Message Types
//Message Types exist in a context of communication
//i.e. is this message an intial request or a reply or part of a stream of events etc
//Request <--> Reply

const (
	REQ = "REQ"
	REP = "REP"
	//PUB
	//SUB
	//HANDSHAKE
)

const (
	CMD_PING           = "PING"    //ping
	CMD_BALANCE        = "BALANCE" //get balance of account
	CMD_FAUCET         = "FAUCET"
	CMD_TX             = "TX"     //send transaction
	CMD_RANDOM_ACCOUNT = "RANACC" //get some random account
	CMD_GETTXPOOL      = "GETTXPOOL"
	//CMD_GETBLOCKS      = "GETBLOCKS"
)

//TODO proper enums
// type MsgType int

// const (
// 	REQ = iota
// 	REP
// )

// func (m MsgType) String() string {
// 	return [...]string{"REQ", "REP"}
// }

//generic message
type Message struct {
	MessageType string //type of message i.e. the communications protocol
	Command     string //Specific message command
	Data        []byte
	//Signature       btcec.Signature
}

type TimeMessage struct {
	Timestamp time.Time
}

func EmptyMsg() Message {
	return Message{}
}

func IsValidMsgType(msgType string) bool {
	//fmt.Println("test ", msgType)
	switch msgType {
	case
		REQ,
		REP:
		return true
	}
	return false
}

func IsValidCmd(cmd string) bool {
	fmt.Println("test cmd ", cmd)
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

type MessageTx struct {
	MessageType string
	Command     string
	Tx          block.Tx
}

type MessageAccount struct {
	MessageType string
	Command     string
	Account     block.Account
}

func AccountMessage(account block.Account) MessageAccount {
	msg := MessageAccount{MessageType: REP, Command: "REP_ACCOUNT", Account: account}
	return msg
}

//encode a message
func EncodeMsg(msgType string, cmd string, data string) string {
	//TODO types
	msg := msgType + string(DELIM_HEAD) + cmd + string(DELIM_HEAD) + data + string(DELIM)
	return msg
}
