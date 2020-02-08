package net

import block "github.com/polygonledger/node/block"

//Message Types
//Request <--> Reply
const (
	REQ = "REQ"
	REP = "REP"
)

//generic message
type Message struct {
	MessageType string
	Command     string //Specific message command
	//Data        []byte
	//Signature       btcec.Signature
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
