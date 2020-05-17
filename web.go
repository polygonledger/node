package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/polygonledger/node/chain"
)

func StatusContent(mgr *chain.ChainManager, t *TCPNode) Status {

	servertime := time.Now()
	uptimedur := time.Now().Sub(t.Starttime)
	uptime := int64(uptimedur / time.Second)
	lastblocktime := t.Mgr.LastBlock().Timestamp
	timebehind := int64(servertime.Sub(lastblocktime) / time.Second)
	status := Status{Blockheight: len(mgr.Blocks), Starttime: t.Starttime, Uptime: uptime, Servertime: servertime, LastBlocktime: lastblocktime, Timebehind: timebehind}
	return status
}

func BlockContent(mgr *chain.ChainManager) string {
	content := ""

	content += fmt.Sprintf("<br><h2>Blocks</h2><i>number of blocks %d</i><br>", len(mgr.Blocks))

	for i := 0; i < len(mgr.Blocks); i++ {
		current_block := mgr.Blocks[i]

		t := current_block.Timestamp
		tsf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		//summary
		content += fmt.Sprintf("<br><h3>Block %d</h3>timestamp %s<br>hash %x<br>prevhash %x\n", current_block.Height, tsf, current_block.Hash, current_block.Prev_Block_Hash)

		content += fmt.Sprintf("<h4>Number of Tx %d</h4>", len(current_block.Txs))
		for j := 0; j < len(current_block.Txs); j++ {
			ctx := current_block.Txs[j]
			content += fmt.Sprintf("%d from %s to %s <br>", ctx.Amount, ctx.Sender, ctx.Receiver)
		}
	}
	return content
}

func AccountContent(mgr *chain.ChainManager) string {

	content := ""
	content += fmt.Sprintf("<h2>Accounts</h2>number of accounts: %d<br><br>", len(mgr.State.Accounts))

	for k, v := range mgr.State.Accounts {
		content += fmt.Sprintf("%s %d<br>", k, v)
	}
	return content
}

func Txpoolcontent(mgr *chain.ChainManager) string {
	content := ""
	content += fmt.Sprintf("<h2>TxPool</h2>%d<br>", len(mgr.Tx_pool))

	for i := 0; i < len(mgr.Tx_pool); i++ {
		//content += fmt.Sprintf("Nonce %d, Id %x<br>", chain.Tx_pool[i].Nonce, chain.Tx_pool[i].Id[:])
		ctx := mgr.Tx_pool[i]
		content += fmt.Sprintf("%d from %s to %s<br>", ctx.Amount, ctx.Sender, ctx.Receiver)
	}
	return content
}

//HTTP
func LoadContent(mgr *chain.ChainManager) string {
	content := ""

	// content += fmt.Sprintf("<h2>Peers</h2>Peers: %d<br>", len(peers))
	// for i := 0; i < len(peers); i++ {
	// 	content += fmt.Sprintf("peer ip address: %s<br>", peers[i].Address)
	// }

	content += Txpoolcontent(mgr)
	content += "<br>"

	content += "<a href=\"/blocks\">blocks</a><br>"
	content += "<a href=\"/accounts\">accounts</a><br>"

	//content += BlockContent(mgr)

	return content
}

func runWeb(t *TCPNode) {
	//webserver to access node state through browser
	// HTTP
	log.Printf("start webserver %d", t.Config.WebPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := LoadContent(t.Mgr)
		statusdata := StatusContent(t.Mgr, t)
		sd, _ := json.Marshal(statusdata)
		//nlog.Print(p)
		fmt.Fprintf(w, "<h1>Polygon chain</h1>Status%s<br><div>%s </div>", sd, p)
	})

	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		p := BlockContent(t.Mgr)
		//nlog.Print(p)
		fmt.Fprintf(w, "<div>%s</div>", p)
	})

	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		p := AccountContent(t.Mgr)
		//nlog.Print(p)
		fmt.Fprintf(w, "<div>%s</div>", p)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {

		status := StatusContent(t.Mgr, t)
		data, _ := json.Marshal(status)

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

		//w.WriteHeader(http.StatusCreated)
		//json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {

		dat, _ := ioutil.ReadFile("node.log")
		fmt.Fprintf(w, "%s", dat)
	})

	//log.Fatal(http.ListenAndServe(":"+strconv.Itoa(webport), nil))
	http.ListenAndServe(":"+strconv.Itoa(t.Config.WebPort), nil)

}
