package main

import (
	"encoding/json"
	"testing"

	protocol "github.com/polygonledger/node/net"
)

func TestAverage(t *testing.T) {

	msg := protocol.Message{MessageType: "msg", Command: "CMD"}
	msgJson, _ := json.Marshal(msg)
	//fmt.Println(string(msgJson))

	var msgUn protocol.Message
	if err := json.Unmarshal(msgJson, &msgUn); err != nil {
		panic(err)
	}
	//fmt.Println(msgUn.Command, msgUn.MessageType)

	if msgUn.Command != "CMD" {
		t.Error("JSON marshal failed")
	}
}
