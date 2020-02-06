package net

import (
	"bufio"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	block "github.com/polygon/block"
	//chain "github.com/polygon/chain"
	cryptoutil "github.com/polygon/crypto"
)

const (
	// Port is the port number that the server listens to.
	Server_address  string = "127.0.0.1"
	Port                   = ":8888"
	CMD_GOB                = "GOB"
	CMD_TX                 = "TX"
	Genesis_Address string = "P0614579c42f2"
)

//pickRandomAccount

//storeBalance

func GenesisTx() block.Tx {
	Genesis_Account := block.AccountFromString(Genesis_Address)
	//block.AccountFromString("") //sender is empty

	//genesisSender := "" //genesisSender is the bootstrap account

	//log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	r_account := block.AccountFromString(address_r)
	genesisAmount := 20 //just a number for now
	//TODO id
	gTx := block.Tx{Nonce: randNonce, Sender: Genesis_Account, Receiver: r_account, Amount: genesisAmount}
	return gTx
}

func RandomTx() block.Tx {
	s := cryptoutil.RandomPublicKey()
	address_s := cryptoutil.Address(s)
	account_s := block.AccountFromString(address_s)

	log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	account_r := block.AccountFromString(address_r)

	//TODO make sure the amount is covered by sender
	randomAmount := rand.Intn(100)
	testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	return testTx
}

/*
The client function connects to the server and sends GOB requests.
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
