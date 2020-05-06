package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/parser"
)

var pw string

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

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/pw", postpw)
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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
