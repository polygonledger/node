package net

import (
	"bufio"
	"encoding/gob"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/polygon/block"
	chain "github.com/polygon/chain"
	cryptoutil "github.com/polygon/crypto"
)

const (
	// Port is the port number that the server listens to.
	Server_address string = "127.0.0.1"
	Port                  = ":8888"
	CMD_GOB               = "GOB"
	CMD_TX                = "TX"
)

//pickRandomAccount

//storeBalance

func GenesisTx() block.Tx {
	emptySender := chain.AccountFromString("") //sender is empty

	//genesisSender := "" //genesisSender is the bootstrap account

	//log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	r := cryptoutil.RandomPublicKey()
	genesisAmount := 20 //just a number for now
	//TODO id
	gTx := block.Tx{Nonce: randNonce, Sender: emptySender, Receiver: chain.AccountFromString(r), Amount: genesisAmount}
	return gTx
}

func RandomTx() block.Tx {
	s := cryptoutil.RandomPublicKey()
	log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	r := cryptoutil.RandomPublicKey()
	//TODO make sure the amount is covered by sender
	randomAmount := rand.Intn(100)
	testTx := block.Tx{Nonce: randNonce, Sender: chain.AccountFromString(s), Receiver: chain.AccountFromString(r), Amount: randomAmount}
	return testTx
}

/*
The client function connects to the server and sends GOB requests.
*/
func SendTx(rw *bufio.ReadWriter) error {

	// Send a GOB request.
	// Create an encoder that directly transmits to `rw`.
	// Send the request name.
	// Send the GOB.

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
