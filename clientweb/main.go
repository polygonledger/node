package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
	"github.com/polygonledger/node/parser"
)

var connect_peer ntcl.Peer

var pw string
var node_status string

type PageData struct {
	PageTitle string
	Password  string
	Pubkey    string
	Privkey   string
	Address   string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreatePubKeypairFormat(pubkey_string string, address string) string {
	mp := map[string]string{"pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
	m := parser.MakeMap(mp)
	return m
}

func CreateKeypairFormat(privkey string, pubkey_string string, address string) string {
	mp := map[string]string{"privkey": parser.StringWrap(privkey), "pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
	m := parser.MakeMap(mp)
	return m
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("basic.html"))

	kp := crypto.PairFromSecret(pw)
	pubkeyHex := crypto.PubKeyToHex(kp.PubKey)
	privHex := crypto.PrivKeyToHex(kp.PrivKey)
	address := crypto.Address(pubkeyHex)

	data := PageData{
		PageTitle: "Polygon client",
		Password:  pw,
		Pubkey:    pubkeyHex,
		Privkey:   privHex,
		Address:   address,
	}

	//tmpl.Execute(w, data)

	switch r.Method {
	case "GET":
		//http.ServeFile(w, r, "form.html")
		tmpl.Execute(w, data)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "r.PostFrom = %v\n", r.PostForm)
		pw = r.FormValue("pw")
		fmt.Fprintf(w, "pw = %s\n", pw)
	default:
		fmt.Fprintf(w, "only GET and POST methods are supported.")
	}
}

func postpw(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		//http.ServeFile(w, r, "form.html")
		rand.Seed(time.Now().UnixNano())
		pw = randSeq(12)
		fmt.Fprintf(w, pw)
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			//fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "r.PostFrom = %v\n", r.PostForm)
		gotpw := r.FormValue("pw")
		pw = gotpw
		//fmt.Fprintf(w, "pw = %s\n", pw)
	default:
		fmt.Fprintf(w, "only GET and POST methods are supported.")
	}
}

func sendnode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			//fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "r.PostFrom = %v\n", r.PostForm)
		message := r.FormValue("message")
		fmt.Println("send message to node ", message)
		reply := sendmsg(connect_peer, message)
		fmt.Println("reply ", reply)
		fmt.Fprintf(w, reply)

	default:
		fmt.Fprintf(w, "only POST method is supported.")
	}
}

func initClient(mainPeerAddress string, verbose bool) ntcl.Ntchan {
	const node_port = 8888
	addr := mainPeerAddress + ":" + strconv.Itoa(node_port)
	log.Println("dial ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("cant run")
		//return
	}

	log.Println("connected")
	ntchan := ntcl.ConnNtchan(conn, "client", addr, verbose)

	go ntcl.ReadLoop(ntchan)
	go ntcl.ReadProcessor(ntchan)
	go ntcl.WriteProcessor(ntchan)
	go ntcl.WriteLoop(ntchan, 300*time.Millisecond)
	return ntchan

}

func sendmsg(peer ntcl.Peer, req_msg string) string {

	peer.NTchan.REQ_out <- req_msg

	time.Sleep(1000 * time.Millisecond)

	reply := <-peer.NTchan.REP_in
	log.Println("reply ", reply)
	return reply

}

func ping(peer ntcl.Peer) {
	ping_msg := ntcl.EncodeMsgMap(ntcl.REQ, ntcl.CMD_PING)
	reply := sendmsg(peer, ping_msg)
	success := reply == "{:REP PONG}"
	log.Println("success ", success)
}

func getStatus(peer ntcl.Peer) {
	node_status = sendmsg(connect_peer, "{:REQ STATUS}")
	fmt.Println(node_status)
	//return node_status
}

func nodestatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getStatus(connect_peer)
		fmt.Fprintf(w, node_status)

	default:
		fmt.Fprintf(w, "only POST method is supported.")
	}
}

func main() {
	peerAddress := "localhost"
	ntchan := initClient(peerAddress, true)

	const node_port = 8888
	connect_peer = ntcl.CreatePeer(peerAddress, peerAddress, node_port, ntchan)

	ping(connect_peer)

	http.HandleFunc("/", index)
	http.HandleFunc("/pw", postpw)
	http.HandleFunc("/nodestatus", nodestatus)
	http.HandleFunc("/sendnode", sendnode)
	http.HandleFunc("/wallet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=wallet.wfe")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		kp := crypto.PairFromSecret(pw)
		pubkeyHex := crypto.PubKeyToHex(kp.PubKey)
		privHex := crypto.PrivKeyToHex(kp.PrivKey)
		address := crypto.Address(pubkeyHex)
		s := CreateKeypairFormat(privHex, pubkeyHex, address)
		ioutil.WriteFile("wallet.wfe", []byte(s), 0644)

		http.ServeFile(w, r, "./wallet.wfe")
	})

	fmt.Printf("Starting client web...\n")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
}
