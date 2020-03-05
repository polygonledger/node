package ntwk

import (
	"encoding/hex"
	"log"
	"math/rand"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
)

//request account address
// func RequestAccount(rw *bufio.ReadWriter) error {
// 	msg := ConstructMessage(CMD_RANDOM_ACCOUNT)

// func ReceiveAccount(rw *bufio.ReadWriter) error {
// 	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)

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
