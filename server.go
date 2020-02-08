package main

//basic server based protocol
//server receives tx messages
//adds tx messages to a pool
//block gets created every 10 secs

//TODO
//package polygonledger/node
//basic signatures

//Tx, sender receiver
//var hash = sha256("secret")
//var keypair = MakeKeypair(hash)

//Delegates
//rounds
//slotTime = getSlotNumber(currentBlockData.time))
//if slotTime generate block

//newWallet

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"encoding/json"

	block "github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
	protocol "github.com/polygonledger/node/net"
)

/*
Incoming connections
*/

// starts listening
func ListenAll() error {
	log.Println("listen all")
	var err error
	var listener net.Listener
	listener, err = net.Listen("tcp", protocol.Port)
	if err != nil {
		log.Println(err)
		return errors.Wrapf(err, "Unable to listen on port %s\n", protocol.Port)
	}

	log.Println("Listen on", listener.Addr().String())
	for {
		log.Println("Accept a connection request")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages")
		go handleMessagesChan(conn)
	}
}

func handleMessagesChan(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//could add max listen
	//timeoutDuration := 5 * time.Second
	//conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	defer conn.Close()
	for {

		msg, err := rw.ReadString(protocol.DELIM)
		if err != nil {
			if len(msg) == 0 {
				break
			} else {
				log.Println("Failed ", err)
				//log.Println(err.)
				break
			}
		}

		// read
		log.Print("Receive message")

		log.Print("msg ", msg, len(msg))

		msg = strings.Trim(msg, string(protocol.DELIM))
		//msg = strings.Trim(msg, string("\n"))
		s := strings.Split(msg, "#")
		mtype := s[0]
		cmd := s[1]
		fmt.Println(mtype, cmd)

		fmt.Println("valid msg type ", protocol.IsValidMsgType(mtype))

		// Add the handle funcs
		//	protocol.CMD_TX, handleTxRequest)
		//	protocol.CMD_RANDOM_ACCOUNT, handleRandomAccountRequest)

		if mtype == protocol.REQ {
			log.Println("Request type")
			if cmd == protocol.CMD_TX {
				log.Println("Handle tx")
				dataJson := s[2]
				//dataBytes := byte[](dataJson)
				dataBytes := []byte(dataJson)

				log.Println("data ", dataJson)

				var tx block.Tx

				if err := json.Unmarshal(dataBytes, &tx); err != nil {
					panic(err)
				}
				log.Println(tx, tx.Amount, tx.Nonce)
				//protocol.CMD_TX, handleTxRequest)
				resp := chain.HandleTx(tx)

				//log.Println("amount ", tx.Amount)
				//n, err := rw.WriteString("response " + strconv.Itoa(tx.Amount) + string(protocol.DELIM))
				n, err := rw.WriteString("response " + resp + string(protocol.DELIM))
				if err != nil {
					log.Println(err, n)
					//err:= errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
				}
				rw.Flush()
			}
		}
	}
}

// handle ranaccount request
func handleRandomAccountRequest(rw *bufio.ReadWriter) {
	protocol.SendAccount(rw)
}

func serverNode() {
	// Start listening
	ListenAll()
}

//basic threading helper
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

//HTTP
func loadContent() string {
	content := ""

	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(chain.Tx_pool))

	for i := 0; i < len(chain.Tx_pool); i++ {
		content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
	}

	content += fmt.Sprintf("<h2>Accounts</h2>%d<br>", len(chain.Accounts))

	for k, v := range chain.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(chain.Blocks))

	for i := 0; i < len(chain.Blocks); i++ {
		t := chain.Blocks[i].Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", chain.Blocks[i].Height, tsf, chain.Blocks[i].Hash, chain.Blocks[i].Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(chain.Blocks[i].Txs))
		for j := 0; j < len(chain.Blocks[i].Txs); j++ {
			ctx := chain.Blocks[i].Txs[j]
			content += fmt.Sprintf("%x, %d from %s to %s<br>", ctx.Id, ctx.Amount, ctx.Sender, ctx.Receiver)
		}
	}

	return content
}

func runweb() {
	//webserver to access node state through browser
	// HTTP
	log.Println("start webserver")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := loadContent()
		//log.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1><div>%s</div>", p)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

/*
start server listening for incoming requests
*/
func main() {

	log.Println("run server")

	//account := block.Account{AccountKey: "test"}
	//accountJson, _ := json.Marshal(account)
	//fmt.Println(string(accountJson))

	/////

	chain.InitAccounts()

	genBlock := chain.MakeGenesisBlock()
	chain.ApplyBlock(genBlock)
	chain.AppendBlock(genBlock)

	// //create block every 10sec
	blockTime := 10000 * time.Millisecond
	go doEvery(blockTime, chain.MakeBlock)

	// //node server
	//go serverNode()
	go ListenAll()
	//log.Println("error ", err)

	// if err != nil {
	// 	log.Println("Error:", errors.WithStack(err))
	// }

	go runweb()
	log.Println("Server running")

}
