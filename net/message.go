package net

import (
	"fmt"

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
	CMD = "RANDOM_ACCOUNT"
	//PUB
	//SUB
	//HANDSHAKE
)

// type Direction int

// const (
//     North Direction = iota
//     East
//     South
//     West
// )

// func (d Direction) String() string {
//     return [...]string{"North", "East", "South", "West"}[d]
// }

//TODO proper enums
//generic message
type Message struct {
	MessageType string //type of message i.e. the communications protocol
	Command     string //Specific message command
	Data        []byte
	//Signature       btcec.Signature
}

func EmptyMsg() Message {
	return Message{}
}

func IsValidMsgType(msgType string) bool {
	fmt.Println("test ", msgType)
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
