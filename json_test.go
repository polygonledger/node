package main

import (
	"testing"
)

func TestMsgJson(t *testing.T) {

	// msg := netio.Message{MessageType: "REQ", Command: netio.CMD_PING}
	// msgJson, _ := json.Marshal(msg)
	// if string(msgJson) != `{"messagetype":"REQ","command":"PING"}` {
	// 	t.Error(string(msgJson))
	// }

	// msgstring := `{"messagetype":"REP","command":"PONG"}`
	// var repmsg netio.Message
	// json.Unmarshal([]byte(msgstring), &repmsg)
	// if repmsg.MessageType != "REP" {
	// 	t.Error(repmsg)
	// }
	// if repmsg.Command != "PONG" {
	// 	t.Error(repmsg)
	// }
}

func TestMsgJsonData(t *testing.T) {

	// msg := netio.Message{MessageType: "REQ", Command: netio.CMD_BALANCE, Data: []byte("abc")}
	// if !(bytes.Compare(msg.Data, []byte("abc")) == 0) {
	// 	i := 1
	// 	t.Error("data bytes not equal ", msg.Data[i], []byte("abc")[i])
	// }
	// msgJson, _ := json.Marshal(msg)
	// if string(msgJson) != `{messagetype:REQ","command":"BALANCE","data":"abc"}` {
	// 	t.Error("? ", string(msgJson), msgJson)
	// }

	// msgstring := `{"messagetype":"REQ","command":"BALANCE","data":"abc"}`
	// var repmsg netio.Message
	// json.Unmarshal([]byte(msgstring), &repmsg)
	// if repmsg.MessageType != "REQ" {
	// 	t.Error(repmsg)
	// }
	// if repmsg.Command != "BALANCE" {
	// 	t.Error(repmsg)
	// }
}
