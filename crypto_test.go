package main

import (
	"testing"

	cryptoutil "github.com/polygonledger/node/crypto"
)

func TestBasicSign(t *testing.T) {

	keypair := cryptoutil.SomeKeypair()
	message := "test message"
	signature := cryptoutil.Sign(keypair, message)
	messageHash := cryptoutil.MsgHash(message)
	verified := signature.Verify(messageHash, &keypair.PubKey)

	if !verified {
		t.Error("msg failed")
	}

	//cryptoutil.KeyExample()

	//btcec.PublicKey
	// s := cryptoutil.RandomPublicKey()
	// log.Printf("%s", s)

}
