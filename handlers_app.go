package main

//this namespace is application specific will ultimately live separated from core e.g. polygon-services

//trivial simplified chat service
//clients need to subscribe with {"MESSAGETYPE":"REQ","COMMAND":"SUBTO"}
//clients send a REQ message

//more abstract note
//in ETH contracts are handled without knowledge and executed blindly
//we want everyone to be able to deploy contracts permissionsless. however, this completely ignored the problem of competing claims
//if say e.g. google, facebook, amazon or any public entity would deploy a contract, reserve a name etc. how we deal with conflicts?
//we want to be able to handle different trust paradigms and deal with real world property

import (
	"encoding/json"
	"fmt"

	"github.com/polygonledger/node/netio"
)

func HandleChat(t *TCPNode, ntchan netio.Ntchan, req_msg netio.Message) string {

	xjson, _ := json.Marshal("hello chat world")
	msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_CHAT, Data: []byte(xjson)}
	reply_msg := netio.ToJSONMessage(msg)

	//TODO this should be in pub chan not REQ out
	//don't publish to author
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
