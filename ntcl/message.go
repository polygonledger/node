package ntcl

import (
	"fmt"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/parser"
)

//TODO does not use edn yet

//any message is defined through the delimters, no size restrictions as of now
//MSG TYPE # CMD # DATA |

//Message Types
//Message Types exist in a context of communication
//i.e. is this message an intial request or a reply or part of a stream of events etc
//Request <--> Reply

const (
	REQ = "REQ"
	REP = "REP"
	PUB = "PUB"
	//SUB
	//HANDSHAKE
)

const (
	CMD_PING             = "PING"
	CMD_PONG             = "PONG"
	CMD_HANDSHAKE_HELLO  = "HELLO"
	CMD_HANDSHAKE_STABLE = "STABLE"
	CMD_HEARTBEAT        = "HEARTBEAT"
	CMD_BALANCE          = "BALANCE" //get balance of account
	CMD_FAUCET           = "FAUCET"
	CMD_TX               = "TX" //send transaction
	CMD_BLOCKHEIGHT      = "BLOCKHEIGHT"
	CMD_RANDOM_ACCOUNT   = "RANACC" //get some random account

	//CMD_LOGIN            = "LOGIN"
	//GET
	CMD_GETTXPOOL = "GETTXPOOL"
	CMD_GETPEERS  = "GETPEERS"
	CMD_GETBLOCKS = "GETBLOCKS"
	EMPTY_DATA    = "EMPTY"
	CMD_SUB       = "SUBTO"
	CMD_SUBUN     = "SUBUN" //unsuscribe
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

	//SRC
	//DEST
	//Signature       btcec.Signature
}

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

type MessageBalance struct {
	MessageType string
	Command     string
	Balance     int
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

func EncodeMsgMap(msgType string, cmd string) string {
	m := map[string]string{msgType: cmd}
	msg := parser.MakeMap(m)
	return msg
}

//////////////////////

func AccountMessage(account block.Account) MessageAccount {
	msg := MessageAccount{MessageType: REP, Command: "REP_ACCOUNT", Account: account}
	return msg
}

//encode a message
func EncodeMsgString(msgType string, cmd string, data string) string {
	//TODO types
	msg := msgType + string(DELIM_HEAD) + cmd + string(DELIM_HEAD) + data + string(DELIM)
	return msg
}

func EncodeMsg(msgType string, cmd string, data string) Message {
	m := Message{MessageType: msgType, Command: cmd, Data: []byte(data)}
	return m
}

func EncodeMsgBytes(msgType string, cmd string, data []byte) Message {
	m := Message{MessageType: msgType, Command: cmd, Data: data}
	return m
}

func MsgString(m Message) string {
	//TODO types
	msg := m.MessageType + string(DELIM_HEAD) + m.Command + string(DELIM_HEAD) + string(m.Data) + string(DELIM)
	return msg
}

func DecodeMsg(msg string) string {

	return msg
}
