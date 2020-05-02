package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/polygonledger/node/crypto"
)

var pw string

type PageData struct {
	PageTitle string
	Password  string
	Pubkey    string
	Privkey   string
	Address   string
}

func hello(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		pw = r.FormValue("pw")
		fmt.Fprintf(w, "pw = %s\n", pw)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	pw = "default"
	http.HandleFunc("/", hello)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
