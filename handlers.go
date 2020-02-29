package main

import "github.com/polygonledger/node/ntwk"

//--- request handler ---

func HandlePing() ntwk.Message {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_PONG, ntwk.EMPTY_DATA)
	return reply
	//msg_out_chan <- reply
}

func HandleHandshake(msg_out_chan chan ntwk.Message) {
	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_HANDSHAKE_STABLE, ntwk.EMPTY_DATA)
	msg_out_chan <- reply
}

func HandleReqMsg(msg ntwk.Message) ntwk.Message {
	nlog.Println("Handle ", msg.Command)

	switch msg.Command {

	case ntwk.CMD_PING:
		nlog.Println("PING PONG")
		return HandlePing()

		// case ntwk.CMD_HANDSHAKE_HELLO:
		// 	nlog.Println("handshake")
		// 	HandleHandshake(rep_chan)

		// case ntwk.CMD_BALANCE:
		// 	nlog.Println("Handle balance")

		// 	dataBytes := msg.Data
		// 	nlog.Println("data ", dataBytes)
		// 	var account block.Account

		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("get balance for account ", account)

		// 	balance := chain.Accounts[account]
		// 	//s := strconv.Itoa(balance)
		// 	data, _ := json.Marshal(balance)
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BALANCE, data)
		// 	log.Println(">> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_FAUCET:
		// 	//send money to specified address

		// 	dataBytes := msg.Data
		// 	var account block.Account
		// 	if err := json.Unmarshal(dataBytes, &account); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println("faucet for ... ", account)

		// 	randNonce := 0
		// 	amount := 10

		// 	keypair := chain.GenesisKeys()
		// 	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))
		// 	Genesis_Account := block.AccountFromString(addr)

		// 	tx := block.Tx{Nonce: randNonce, Amount: amount, Sender: Genesis_Account, Receiver: account}

		// 	tx = crypto.SignTxAdd(tx, keypair)
		// 	reply_string := chain.HandleTx(tx)
		// 	nlog.Println("resp > ", reply_string)

		// 	reply := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_FAUCET, reply_string)

		// 	rep_chan <- reply

		// case ntwk.CMD_BLOCKHEIGHT:

		// 	data, _ := json.Marshal(len(chain.Blocks))
		// 	reply := ntwk.EncodeMsgBytes(ntwk.REP, ntwk.CMD_BLOCKHEIGHT, data)
		// 	log.Println("CMD_BLOCKHEIGHT >> ", reply)

		// 	rep_chan <- reply

		// case ntwk.CMD_TX:
		// 	nlog.Println("Handle tx")

		// 	dataBytes := msg.Data

		// 	var tx block.Tx

		// 	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		// 		panic(err)
		// 	}
		// 	nlog.Println(">> ", tx)

		// 	resp := chain.HandleTx(tx)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		// 	rep_chan <- msg

		// case ntwk.CMD_GETTXPOOL:
		// 	nlog.Println("get tx pool")

		// 	//TODO
		// 	data, _ := json.Marshal(chain.Tx_pool)
		// 	msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_GETTXPOOL, string(data))
		// 	rep_chan <- msg

		//var Tx_pool []block.Tx

		// case ntwk.CMD_RANDOM_ACCOUNT:
		// 	nlog.Println("Handle random account")

		// 	txJson, _ := json.Marshal(chain.RandomAccount())

	default:
		// 	nlog.Println("unknown cmd ", msg.Command)
		resp := "ERROR UNKONWN CMD"

		msg := ntwk.EncodeMsg(ntwk.REP, ntwk.CMD_TX, resp)
		return msg

	}
}
