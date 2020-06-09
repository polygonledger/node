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
	if string(jsonmsg) != `{"messagetype":"REQ","command":"ACCOUNTS"}` {
		t.Error("unmarshal")
	}

	jsonmsg2 := netio.ToJSONMessage(msg)
	if string(jsonmsg2) != `{"messagetype":"REQ","command":"ACCOUNTS"}` {
		t.Error("unmarshal")
	}

}

func TestMessagePing(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING} //, Data: []byte("test")}
	jsonmsg := netio.ToJSONMessage(msg)

	if jsonmsg != `{"messagetype":"REQ","command":"PING"}` {
		t.Error(jsonmsg)
	}

	//msg2 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING, Data: []byte("test")}
	msg2 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING, Data: []byte(`{"test":"test"}`)}
	jsonmsg2 := netio.ToJSONMessage(msg2)

	if jsonmsg2 != `{"messagetype":"REQ","command":"PING","data":{"test":"test"}}` {
		t.Error(jsonmsg2)
	}

	msg3 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING, Data: []byte("20")}
	jsonmsg3 := netio.ToJSONMessage(msg3)

	if string(jsonmsg3) != `{"messagetype":"REQ","command":"PING","data":20}` {
		t.Error(jsonmsg3)
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

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMessageBalance(t *testing.T) {

	balJson, _ := json.Marshal(20)
	msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
	jsonmsg := netio.ToJSONMessage(msg)

	if string(jsonmsg) != `{"messagetype":"REP","command":"BALANCE","data":20}` {
		t.Error("unmarshal ", string(jsonmsg))
	}

	var msgu netio.Message
	json.Unmarshal([]byte(`{"messagetype":"REP","command":"BALANCE","data":20}`), &msgu)

	if msgu.MessageType != "REP" || msgu.Command != netio.CMD_BALANCE {
		t.Error(msgu)
	}

	if !Equal(msg.Data, []byte("20")) {
		t.Error(msgu.Data)
	}

	xJson, _ := json.Marshal("test")
	msg = netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(xJson)}
	jsonmsg = netio.ToJSONMessage(msg)

	if string(jsonmsg) != `{"messagetype":"REP","command":"BALANCE","data":"test"}` {
		t.Error("unmarshal ", string(jsonmsg))
	}

}

func TestMessagePeers(t *testing.T) {
	p := netio.Peer{Address: "test", Name: "test", NodePort: 80}
	if p.Address != "test" {
		t.Error(p)
	}
	peer_json, _ := json.Marshal(p)
	sj := string(peer_json)
	if sj != `{"Address":"test","Name":"test","NodePort":80}` {
		t.Error("peer json ", sj, len(sj))
	}

}

func TestMessageChat(t *testing.T) {

	msg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_CHAT} //, Data: []byte("test")}
	jsonmsg := netio.ToJSONMessage(msg)

	if jsonmsg != `{"messagetype":"REQ","command":"CHAT"}` {
		t.Error(jsonmsg)
	}

	//msg2 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_PING, Data: []byte("test")}
	msg2 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_CHAT, Data: []byte(`{"test":"test"}`)}
	jsonmsg2 := netio.ToJSONMessage(msg2)

	if jsonmsg2 != `{"messagetype":"REQ","command":"CHAT","data":{"test":"test"}}` {
		t.Error(jsonmsg2)
	}

	msg3 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_CHAT, Data: []byte("20")}
	jsonmsg3 := netio.ToJSONMessage(msg3)

	if string(jsonmsg3) != `{"messagetype":"REQ","command":"CHAT","data":20}` {
		t.Error(jsonmsg3)
	}

	//TODO! fails
	// msg4 := netio.Message{MessageType: netio.REQ, Command: netio.CMD_CHAT, Data: []byte("test")}
	// jsonmsg4 := netio.ToJSONMessage(msg4)

	// if string(jsonmsg4) != `{"messagetype":"REQ","command":"CHAT","data":"test"}` {
	// 	t.Error(jsonmsg4)
	// }
}

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
