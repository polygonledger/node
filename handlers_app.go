package main

import (
	"encoding/json"
	"fmt"

	"github.com/polygonledger/node/netio"
)

//TODO in handlers_app.go
func HandleChat(t *TCPNode, ntchan netio.Ntchan, req_msg netio.Message) string {

	xjson, _ := json.Marshal("hello chat world")
	msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_CHAT, Data: []byte(xjson)}
	reply_msg := netio.ToJSONMessage(msg)

	//TODO this should be in publoop
	for _, subs := range t.ChatSubscribers {
		pub_msg_string := fmt.Sprintf("%v said: publish chat world", ntchan)

		xjson, _ := json.Marshal(pub_msg_string)
		othermsg := netio.Message{MessageType: netio.PUB, Command: netio.CMD_CHAT, Data: []byte(xjson)}
		xmsg := netio.ToJSONMessage(othermsg)
		fmt.Println("send to ", subs, xmsg)
		subs.REQ_out <- xmsg
	}

	return reply_msg
}
