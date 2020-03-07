package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/ntcl"
)

func ReadKeys(keysfile string) crypto.Keypair {

	dat, _ := ioutil.ReadFile(keysfile)
	s := strings.Split(string(dat), string("\n"))

	pubkeyHex := s[0]
	log.Println("pub ", pubkeyHex)

	privHex := s[1]
	log.Println("privHex ", privHex)

	return crypto.Keypair{PubKey: crypto.PubKeyFromHex(pubkeyHex), PrivKey: crypto.PrivKeyFromHex(privHex)}
}

func WriteKeys(kp crypto.Keypair, keysfile string) {

	pubkeyHex := crypto.PubKeyToHex(kp.PubKey)
	log.Println("pub ", pubkeyHex)

	privHex := crypto.PrivKeyToHex(kp.PrivKey)
	log.Println("privHex ", privHex)

	address := crypto.Address(pubkeyHex)

	t := pubkeyHex + "\n" + privHex + "\n" + address
	//log.Println("address ", address)
	ioutil.WriteFile(keysfile, []byte(t), 0644)
}

func CreateKeys() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.Trim(pw, string('\n'))
	fmt.Println(pw)

	//check if exists
	//dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)

	kp := crypto.PairFromSecret(pw)
	log.Println("keypair ", kp)

	WriteKeys(kp, "keys.txt")

}

func Createtx() {
	kp := ReadKeys("keys.txt")

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter password: ")
	// pw, _ := reader.ReadString('\n')
	// pw = strings.Trim(pw, string('\n'))

	// keypair := crypto.PairFromSecret(pw)

	pubk := crypto.PubKeyToHex(kp.PubKey)
	addr := crypto.Address(pubk)
	s := block.AccountFromString(addr)
	log.Println("using account ", s)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter amount: ")
	amount, _ := reader.ReadString('\n')
	amount = strings.Trim(amount, string('\n'))
	amount_int, _ := strconv.Atoi(amount)

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter recipient: ")
	recv, _ := reader.ReadString('\n')
	recv = strings.Trim(recv, string('\n'))

	tx := block.Tx{Nonce: 1, Amount: amount_int, Sender: block.Account{AccountKey: addr}, Receiver: block.Account{AccountKey: recv}}
	log.Println("tx ", tx)

	signature := crypto.SignTx(tx, kp)
	sighex := hex.EncodeToString(signature.Serialize())

	tx.Signature = sighex
	tx.SenderPubkey = crypto.PubKeyToHex(kp.PubKey)
	log.Println("tx ", tx)

	txJson, _ := json.Marshal(tx)
	// //write to file
	// log.Println(txJson)

	ioutil.WriteFile("tx.json", []byte(txJson), 0644)
}

//
func MakeRandomTx(peer ntcl.Peer) {
	//make a random transaction by requesting random account from node
	//get random account

	// req_msg := ntcl.EncodeMsgString(ntcl.REQ, ntcl.CMD_RANDOM_ACCOUNT, "emptydata")

	// response := ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)

	// var a block.Account
	// dataBytes := []byte(response.Data)
	// if err := json.Unmarshal(dataBytes, &a); err != nil {
	// 	panic(err)
	// }
	// log.Print(" account key ", a.AccountKey)

	// //use this random account to send coins from

	// //send Tx
	// testTx := ntcl.RandomTx(a)
	// txJson, _ := json.Marshal(testTx)
	// log.Println("txJson ", txJson)

	// req_msg = ntcl.EncodeMessageTx(txJson)

	// response = ntcl.RequestReplyChan(req_msg, peer.Req_chan, peer.Rep_chan)
	// log.Print("response msg ", response)

	// return nil
}

//handlers TODO this is higher level and should be somewhere else
func RandomTx(account_s block.Account) block.Tx {
	// s := crypto.RandomPublicKey()
	// address_s := crypto.Address(s)
	// account_s := block.AccountFromString(address_s)
	// log.Printf("%s", s)

	//FIX
	//doesn't work on client side
	//account_r := chain.RandomAccount()

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)

	kp := crypto.PairFromSecret("test111??")
	log.Println("PUBKEY ", kp.PubKey)

	r := crypto.RandomPublicKey()
	address_r := crypto.Address(r)
	account_r := block.AccountFromString(address_r)

	//TODO make sure the amount is covered by sender
	rand.Seed(time.Now().UnixNano())
	randomAmount := rand.Intn(20)

	log.Printf("randomAmount ", randomAmount)
	log.Printf("randNonce ", randNonce)
	testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	sig := crypto.SignTx(testTx, kp)
	sighex := hex.EncodeToString(sig.Serialize())
	testTx.Signature = sighex
	log.Println(">> ran tx", testTx.Signature)
	return testTx
}

// func CreateTx(peer ntcl.Peer) {
// 	// keypair := crypto.PairFromSecret("test")
// 	// var tx block.Tx
// 	// s := block.AccountFromString("Pa033f6528cc1")
// 	// r := s //TODO
// 	// tx = block.Tx{Nonce: 0, Amount: 0, Sender: s, Receiver: r}
// 	// signature := crypto.SignTx(tx, keypair)
// 	// sighex := hex.EncodeToString(signature.Serialize())
// 	// tx.Signature = sighex
// 	// tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

// }
