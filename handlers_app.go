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

	countout := 0
	//TODO this should be in pub chan not REQ out
	//don't publish to author
	for _, subs := range t.ChatSubscribers {
		if subs == ntchan {
			fmt.Println("dont publish to self")
		} else {
			pub_msg_string := fmt.Sprintf("%v said: %v", ntchan, string(req_msg.Data))

			xjson, _ := json.Marshal(pub_msg_string)
			othermsg := netio.Message{MessageType: netio.PUB, Command: netio.CMD_CHAT, Data: []byte(xjson)}
			xmsg := netio.ToJSONMessage(othermsg)
			fmt.Println("send to ", subs, xmsg)
			subs.REQ_out <- xmsg
			countout++
		}
	}

	msgstring := fmt.Sprintf("message received. delievered to %v", countout)
	xjson, _ := json.Marshal(msgstring)
	msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_CHAT, Data: []byte(xjson)}
	reply_msg := netio.ToJSONMessage(msg)

	return reply_msg
}
