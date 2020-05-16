package main

import (
	"encoding/json"
	"testing"

	"github.com/polygonledger/node/netio"
)

func TestMessageBasicParseJson(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_ACCOUNTS}
	jsonmsgtype, _ := netio.NewJSONMessage(msg)
	jsonmsg, err := json.Marshal(jsonmsgtype)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(b))

	if string(jsonmsg) != `{"messagetype":"REQ","command":"ACCOUNTS"}` {
		t.Error("unmarshal")
	}

	jsonmsg2 := netio.ToJSONMessage(msg)
	if string(jsonmsg2) != `{"messagetype":"REQ","command":"ACCOUNTS"}` {
		t.Error("unmarshal")
	}

}

func TestMessageBasicParseJsonData(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_ACCOUNTS, Data: []byte(`{"test":"test"}`)}
	jsonmsgtype, _ := netio.NewJSONMessage(msg)
	jsonmsg, err := json.Marshal(jsonmsgtype)
	if err != nil {
		panic(err)
	}

	if string(jsonmsg) != `{"messagetype":"REQ","command":"ACCOUNTS","data":{"test":"test"}}` {
		t.Error("unmarshal ", string(jsonmsg))
	}

	jsonmsg2 := netio.ToJSONMessage(msg)
	if string(jsonmsg2) != `{"messagetype":"REQ","command":"ACCOUNTS","data":{"test":"test"}}` {
		t.Error("unmarshal ", string(jsonmsg2))
	}

}

// jsonmsg, _ := json.Marshal(msg)
// if string(jsonmsg) != `{"messagetype":"REQ","command":"CMD"}` {
// 	t.Error("jsonmsg ", string(jsonmsg))
// }

// var JSONMsgData = []byte(`{
// 	"messagetype": "REQ",
// 	"command": "CMD_ACCOUNT",
// 	"data": {
// 		 "account": "abc"
// 	}}`)

// var m netio.Message
// json.Unmarshal(JSONMsgData, &m)

// 	fmt.Println(err)
// }
// fmt.Printf("JSON type is [%v] and [%v] \n", m.Type, e.ABC)

// if !netio.IsValidMsgType(msg.MessageType) {
// 	t.Error("msg type invalid")

// func TestMsgJson(t *testing.T) {
// 	msg := netio.Message{MessageType: "REQ", Command: netio.CMD_PING}
// 	msgJson, _ := json.Marshal(msg)
// 	if string(msgJson) != `{"messagetype":"REQ","command":"PING"}` {
// 		t.Error(string(msgJson))
// 	}

// 	msgstring := `{"messagetype":"REP","command":"PONG"}`
// 	var repmsg netio.Message
// 	json.Unmarshal([]byte(msgstring), &repmsg)
// 	if repmsg.MessageType != "REP" {
// 		t.Error(repmsg)
// 	}
// 	if repmsg.Command != "PONG" {
// 		t.Error(repmsg)
// 	}
// }

// func TestMsgJsonData(t *testing.T) {

// 	msg := netio.Message{MessageType: "REQ", Command: netio.CMD_BALANCE, Data: []byte("abc")}
// 	if !(bytes.Compare(msg.Data, []byte("abc")) == 0) {
// 		i := 1
// 		t.Error("data bytes not equal ", msg.Data[i], []byte("abc")[i])
// 	}
// 	msgJson, _ := json.Marshal(msg)
// 	if string(msgJson) != `{messagetype:REQ","command":"BALANCE","data":"abc"}` {
// 		t.Error("? ", string(msgJson), msgJson)
// 	}

// 	msgstring := `{"messagetype":"REQ","command":"BALANCE","data":"abc"}`
// 	var repmsg netio.Message
// 	json.Unmarshal([]byte(msgstring), &repmsg)
// 	if repmsg.MessageType != "REQ" {
// 		t.Error(repmsg)
// 	}
// 	if repmsg.Command != "BALANCE" {
// 		t.Error(repmsg)
// 	}
// }
