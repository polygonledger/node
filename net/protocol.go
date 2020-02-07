package net

import (
	"bufio"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	block "github.com/polygonledger/node/block"
	cryptoutil "github.com/polygonledger/node/crypto"
)

const (
	// Port is the port number that the server listens to.
	Server_address     string = "127.0.0.1"
	Port                      = ":8888"
	CMD_GOB                   = "GOB"
	CMD_TX                    = "TX"
	CMD_RANDOM_ACCOUNT        = "RANACC"
	Genesis_Address    string = "P0614579c42f2"
)

//pickRandomAccount

//storeBalance

func RandomTx() block.Tx {
	s := cryptoutil.RandomPublicKey()
	address_s := cryptoutil.Address(s)
	account_s := block.AccountFromString(address_s)
	log.Printf("%s", s)

	//FIX
	//doesn't work on client side
	//account_s := chain.RandomAccount()
	//account_r := chain.RandomAccount()

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)

	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	account_r := block.AccountFromString(address_r)

	//TODO make sure the amount is covered by sender
	rand.Seed(time.Now().UnixNano())
	randomAmount := rand.Intn(20)

	log.Printf("randomAmount ", randomAmount)
	log.Printf("randNonce ", randNonce)
	testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	return testTx
}

/*
request account address
*/
func RequestAccount(rw *bufio.ReadWriter) error {
	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)
	n, err := rw.WriteString(CMD_RANDOM_ACCOUNT + "\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	} else {
		log.Println(strconv.Itoa(n) + " bytes written")
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

/*
sends account address
*/
func SendAccount(rw *bufio.ReadWriter) error {

	a := block.Account{AccountKey: "test"}

	log.Println("Send a struct as GOB:")
	log.Printf("account: \n%#v\n", a)

	enc := gob.NewEncoder(rw)
	//Command

	n, err := rw.WriteString(CMD_RANDOM_ACCOUNT + "\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(a)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", a)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

/*
sends GOB requests (client to server)
*/
func SendTx(rw *bufio.ReadWriter) error {

	// Send a GOB request
	// Create an encoder that directly transmits to `rw`.
	// Send the request name
	// Send the GOB data

	testTx := RandomTx()

	log.Println("Send a struct as GOB:")
	log.Printf("testTx: \n%#v\n", testTx)

	enc := gob.NewEncoder(rw)
	//Command

	n, err := rw.WriteString(CMD_TX + "\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(testTx)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", testTx)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}
