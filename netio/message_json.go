package netio

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
