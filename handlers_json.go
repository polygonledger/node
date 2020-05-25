package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/netio"
)

//--- handlers ---

//alternative
// m := map[string]interface{}{
// 	"handlef": f,
// }
// 	switch fcall {
// 	case "f":
// 		func()
// 	}
// }

func HandleEcho(ins string) string {
	resp := "Echo:" + ins
	return resp
}

func HandlePing(msg netio.Message) netio.Message {
	// validRequest := msg.MessageType == netio.REQ && msg.Command == "PING"
	// if !validRequest{
	// 	//error
	// }
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_PONG}
	return reply_msg
}

func HandleBlockheight(t *TCPNode, msg netio.Message) netio.Message {
	bh := len(t.Mgr.Blocks)
	//data := strconv.Itoa(bh)
	djson, _ := json.Marshal(bh)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BLOCKHEIGHT, Data: []byte(djson)}
	//reply_msg := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_BLOCKHEIGHT, data)
	//("BLOCKHEIGHT ", reply_msg)
	return reply_msg
}

//Standard Tx handler
//InteractiveTx also possible
//client requests tranaction <=> server response with challenge <=> client proves
func HandleTx(t *TCPNode, msg netio.Message) netio.Message {
	dataBytes := msg.Data

	var tx block.Tx

	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		panic(err)
	}
	t.log(fmt.Sprintf("tx %v ", tx))

	reply := chain.HandleTx(t.Mgr, tx)
	return reply
}

func HandleBalance(t *TCPNode, msg netio.Message) string {
	dataBytes := msg.Data
	t.log(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))
	fmt.Println(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))

	//a := block.Account{AccountKey: string(msg.Data)}

	// var account block.Account

	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	a := string(msg.Data)
	balance := t.Mgr.State.Accounts[a]
	fmt.Println("balance for ", a, balance, t.Mgr.State.Accounts)

	balJson, _ := json.Marshal(balance)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
	reply_msg_json := netio.ToJSONMessage(reply_msg)
	//reply_msg := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_BALANCE, data)

	return reply_msg_json
}

func HandlePeers(t *TCPNode, msg netio.Message) netio.Message {
	peersJson, _ := json.Marshal(t.Peers)
	reply_msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_GETPEERS, Data: []byte(peersJson)}
	return reply_msg
}

func HandleFaucet(t *TCPNode, msg netio.Message) netio.Message {
	//t.log(fmt.Sprintf("HandleFaucet"))

	// dataBytes := msg.Data
	// var account block.Account
	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	//account := block.Account{AccountKey: string(msg.Data)}
	//t.log(fmt.Sprintf("faucet for ... %v", account.AccountKey))

	randNonce := 0
	amount := rand.Intn(10)

	keypair := chain.GenesisKeys()
	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
	//Genesis_Account := block.AccountFromString(addr)

	//tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}
	a := string(msg.Data)
	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: addr, Receiver: a}

	tx = crypto.SignTxAdd(tx, keypair)
	reply := chain.HandleTx(t.Mgr, tx)
	//t.log(fmt.Sprintf("resp > %s", reply_string))

	//reply := netio.EdnConstructMsgMapData(netio.REP, netio.CMD_FAUCET, reply_string)
	return reply
}

func RequestReply(t *TCPNode, ntchan netio.Ntchan, msg netio.Message) string {

	var reply_msg string
	//var reply_msg netio.Message

	t.log(fmt.Sprintf("Handle cmd %v", msg.Command))

	switch msg.Command {

	case netio.CMD_PING:
		reply := HandlePing(msg)
		//msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_BALANCE, Data: []byte(balJson)}
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_NUMACCOUNTS:

		numacc := len(t.Mgr.State.Accounts)
		t.log(fmt.Sprintf("numacc %v", numacc))
		nJson, _ := json.Marshal(numacc)
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_NUMACCOUNTS, Data: []byte(nJson)}
		reply_msg = netio.ToJSONMessage(msg)
		t.log(fmt.Sprintf("reply %s", reply_msg))

	case netio.CMD_ACCOUNTS:
		nJson, _ := json.Marshal(t.Mgr.State.Accounts)
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_ACCOUNTS, Data: []byte(nJson)}
		reply_msg = netio.ToJSONMessage(msg)

	case netio.CMD_GETPEERS:
		reply := HandlePeers(t, msg)
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_BALANCE:
		reply_msg = HandleBalance(t, msg)

	case netio.CMD_FAUCET:
		//send money to specified address
		reply := HandleFaucet(t, msg)
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_BLOCKHEIGHT:
		reply := HandleBlockheight(t, msg)
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_GETTXPOOL:
		t.log("get tx pool")

		//TODO!
		data, _ := json.Marshal(t.Mgr.Tx_pool)
		reply_msg = netio.EdnConstructMsgMapData(netio.REP, netio.CMD_GETTXPOOL, string(data))

	case netio.CMD_GETBLOCKS:
		t.log("get tx pool")

		//TODO!
		data, _ := json.Marshal(t.Mgr.Blocks)
		reply_msg = netio.EdnConstructMsgMapData(netio.REP, netio.CMD_GETBLOCKS, string(data))

		//Login would be challenge response protocol
		// case netio.CMD_LOGIN:
		// 	log.Println("> ", msg.Data)

	case netio.CMD_TX:
		t.log("Handle tx")
		reply := HandleTx(t, msg)
		//reply_msg = netio.EdnConstructMsgMapData(netio.REP, netio.CMD_TX, reply_msg)
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_RANDOM_ACCOUNT:
		t.log("Handle random account")

		//TODO!
		txJson, _ := json.Marshal(t.Mgr.RandomAccount())
		reply_msg = netio.EdnConstructMsgMapData(netio.REP, netio.CMD_RANDOM_ACCOUNT, string(txJson))

	case netio.CMD_STATUS:
		status := StatusContent(t.Mgr, t)
		data, _ := json.Marshal(status)
		reply := netio.Message{MessageType: netio.REP, Command: netio.CMD_STATUS, Data: data}
		reply_msg = netio.ToJSONMessage(reply)

	case netio.CMD_NUMCONN:
		pn := len(t.GetPeers())
		nJson, _ := json.Marshal(pn)
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_NUMCONN, Data: []byte(nJson)}
		reply_msg = netio.ToJSONMessage(msg)

	//TODO separate handle process
	//PUBSUB
	case netio.CMD_SUB:
		t.log(fmt.Sprintf("subscribe to topic %v", msg.Data))

		t.ChatSubscribers = append(t.ChatSubscribers, ntchan)
		t.log(fmt.Sprintf("subscribers %v", t.ChatSubscribers))

		/////////////////////
		//EXAMPLE publishtime
		//quitpub := make(chan int)
		//go netio.PublishTime(ntchan)
		// go netio.PubWriterLoop(ntchan)
		//TODO reply sub ok

	case netio.CMD_SUBUN:
		t.log(fmt.Sprintf("unsubscribe from topic %v", msg.Data))

		go func() {
			//time.Sleep(5000 * time.Millisecond)
			close(ntchan.PUB_out)
		}()

		//TODO reply unsub ok

	case netio.CMD_LOGOFF:
		reply_msg = "{:REP BYE}"
		ntchan.Writer_queue <- reply_msg
		time.Sleep(500 * time.Millisecond)
		ntchan.Conn.Close()

	// app layer
	case netio.CMD_CHAT:
		xjson, _ := json.Marshal("hello chat world")
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_CHAT, Data: []byte(xjson)}
		reply_msg = netio.ToJSONMessage(msg)

		//TODO this should be in publoop
		for _, subs := range t.ChatSubscribers {
			pub_msg_string := fmt.Sprintf("%v said: publish chat world", ntchan)

			xjson, _ := json.Marshal(pub_msg_string)
			othermsg := netio.Message{MessageType: netio.PUB, Command: netio.CMD_CHAT, Data: []byte(xjson)}
			xmsg := netio.ToJSONMessage(othermsg)
			fmt.Println("send to ", subs, xmsg)
			subs.REQ_out <- xmsg
		}

	default:
		errormsg := "Error: not found command"
		fmt.Println(errormsg)
		xjson, _ := json.Marshal("")
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_ERROR, Data: []byte(xjson)}
		reply_msg = netio.ToJSONMessage(msg)
	}

	return reply_msg
}

//handle requests in telnet style. messages are edn based
func RequestHandlerTel(t *TCPNode, ntchan netio.Ntchan) {
	for {
		msg_string := <-ntchan.REQ_in
		t.log(fmt.Sprintf(">> handle request %s ", msg_string))

		//msg := netio.EdnParseMessageMap(msg_string)
		msg := netio.FromJSON(msg_string)
		//json.Unmarshal([]byte(msg_string), &msg)
		t.log(fmt.Sprintf(">> handle request msg %s ", msg))

		reply_msg := RequestReply(t, ntchan, msg)

		//TODO parse out, i.e not return just a string

		t.log(fmt.Sprintf("reply_msg %s", reply_msg))
		ntchan.Writer_queue <- reply_msg

	}
}
