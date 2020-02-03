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
)

const (
	// Port is the port number that the server listens to.
	Server_address string = "127.0.0.1"
	Port                  = ":8888"
	CMD_GOB               = "GOB\n"
)

/*
The client function connects to the server and sends GOB requests.
*/
func SendTx(rw *bufio.ReadWriter) error {

	// Send a GOB request.
	// Create an encoder that directly transmits to `rw`.
	// Send the request name.
	// Send the GOB.

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	testTx := block.Tx{Nonce: randNonce}

	log.Println("Send a struct as GOB:")
	log.Printf("testTx: \n%#v\n", testTx)

	enc := gob.NewEncoder(rw)
	//Command

	n, err := rw.WriteString(CMD_GOB)
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
