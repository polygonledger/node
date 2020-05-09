package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/polygonledger/edn"
	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/chain"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
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

func HandlePing(msg ntcl.Message) ntcl.Message {
	// validRequest := msg.MessageType == ntcl.REQ && msg.Command == "PING"
	// if !validRequest{
	// 	//error
	// }
	reply_msg := ntcl.EncodeMsgMap(ntcl.REP, "PONG")
	m := ntcl.ParseMessageMap(reply_msg)
	return m
}

func HandleBlockheight(t *TCPNode, msg ntcl.Message) string {
	bh := len(t.Mgr.Blocks)
	data := strconv.Itoa(bh)
	reply_msg := ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_BLOCKHEIGHT, data)
	//("BLOCKHEIGHT ", reply_msg)
	return reply_msg
}

//Standard Tx handler
//InteractiveTx also possible
//client requests tranaction <=> server response with challenge <=> client proves
func HandleTx(t *TCPNode, msg ntcl.Message) string {
	dataBytes := msg.Data

	var tx block.Tx

	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		panic(err)
	}
	t.log(fmt.Sprintf("tx %v ", tx))

	resp := chain.HandleTx(t.Mgr, tx)
	reply := ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_TX, resp)
	return reply
}

func HandleBalance(t *TCPNode, msg ntcl.Message) string {
	dataBytes := msg.Data
	t.log(fmt.Sprintf("HandleBalance data %v %v", string(msg.Data), dataBytes))

	//a := block.Account{AccountKey: string(msg.Data)}

	// var account block.Account

	// if err := json.Unmarshal(dataBytes, &account); err != nil {
	// 	panic(err)
	// }

	a := string(msg.Data)
	//fmt.Println("balance for ", a, msg)
	balance := t.Mgr.Accounts[a]

	//s := strconv.Itoa(balance)
	// data, _ := json.Marshal(balance)
	data := strconv.Itoa(balance)
	reply_msg := ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_BALANCE, data)
	return reply_msg
}

func HandleFaucet(t *TCPNode, msg ntcl.Message) string {
	t.log(fmt.Sprintf("HandleFaucet"))
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
	reply_string := chain.HandleTx(t.Mgr, tx)
	t.log(fmt.Sprintf("resp > %s", reply_string))

	reply := ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_FAUCET, reply_string)
	return reply
}

//handle requests in telnet style. messages are edn based
func RequestHandlerTel(t *TCPNode, ntchan ntcl.Ntchan) {
	for {
		msg_string := <-ntchan.REQ_in
		t.log(fmt.Sprintf("?? handle request %s ", msg_string))

		msg := ntcl.ParseMessageMap(msg_string)

		var reply_msg string
		//var reply_msg ntcl.Message

		t.log(fmt.Sprintf("Handle cmd %v", msg.Command))

		switch msg.Command {

		case ntcl.CMD_PING:
			reply := HandlePing(msg)
			reply_msg = ntcl.EncodeMsgMapS(reply)

		case ntcl.CMD_NUMACCOUNTS:
			numacc := len(t.Mgr.Accounts)
			data := strconv.Itoa(numacc)
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_NUMACCOUNTS, data)

		case ntcl.CMD_ACCOUNTS:
			dat, _ := edn.Marshal(t.Mgr.Accounts)
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_ACCOUNTS, string(dat))

		case ntcl.CMD_STATUS:
			statusdata := string(StatusContent(t.Mgr, t))
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_STATUS, statusdata)

		case ntcl.CMD_NUMCONN:
			pn := len(t.GetPeers())
			data := strconv.Itoa(pn)
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_NUMCONN, data)

		case ntcl.CMD_BALANCE:
			reply_msg = HandleBalance(t, msg)

		case ntcl.CMD_FAUCET:
			//send money to specified address
			reply_msg = HandleFaucet(t, msg)

		case ntcl.CMD_BLOCKHEIGHT:
			reply_msg = HandleBlockheight(t, msg)

		case ntcl.CMD_GETTXPOOL:
			t.log("get tx pool")

			//TODO
			data, _ := json.Marshal(t.Mgr.Tx_pool)
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_GETTXPOOL, string(data))

		case ntcl.CMD_GETBLOCKS:
			t.log("get tx pool")

			//TODO
			data, _ := json.Marshal(t.Mgr.Blocks)
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_GETBLOCKS, string(data))

			//Login would be challenge response protocol
			// case ntcl.CMD_LOGIN:
			// 	log.Println("> ", msg.Data)

		case ntcl.CMD_TX:
			t.log("Handle tx")
			reply_msg = HandleTx(t, msg)

		case ntcl.CMD_RANDOM_ACCOUNT:
			t.log("Handle random account")

			txJson, _ := json.Marshal(t.Mgr.RandomAccount())
			reply_msg = ntcl.EncodeMsgMapData(ntcl.REP, ntcl.CMD_RANDOM_ACCOUNT, string(txJson))

		//TODO separate handle process
		//PUBSUB
		case ntcl.CMD_SUB:
			t.log(fmt.Sprintf("subscribe to topic %v", msg.Data))

			//quitpub := make(chan int)
			go ntcl.PublishTime(ntchan)
			go ntcl.PubWriterLoop(ntchan)
			//TODO reply sub ok

		case ntcl.CMD_SUBUN:
			t.log(fmt.Sprintf("unsubscribe from topic %v", msg.Data))

			go func() {
				//time.Sleep(5000 * time.Millisecond)
				close(ntchan.PUB_time_quit)
			}()

			//TODO reply unsub ok

		case ntcl.CMD_LOGOFF:
			reply_msg = "{:REP BYE}"
			ntchan.Writer_queue <- reply_msg
			time.Sleep(500 * time.Millisecond)
			ntchan.Conn.Close()

		}

		t.log(fmt.Sprintf("reply_msg %s", reply_msg))
		ntchan.Writer_queue <- reply_msg

	}
}
