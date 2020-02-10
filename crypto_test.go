package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/btcd/btcec"
	"github.com/btcd/chaincfg/chainhash"
	cryptoutil "github.com/polygonledger/node/crypto"
)

func TestBasicSign(t *testing.T) {

	//keypair := cryptoutil.SomeKeypair()
	keypair := cryptoutil.PairFromSecret("test")
	message := "test message"

	log.Println(keypair.PubKey)

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

func TestDecode(t *testing.T) {
	pubKeyBytes, err := hex.DecodeString("02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode hex-encoded serialized signature.
	sigBytes, err := hex.DecodeString("30450220090ebfb3690a0ff115bb1b38b" +
		"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
		"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479")

	if err != nil {
		fmt.Println(err)
		return
	}
	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		return
	}

	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)

	if !verified {
		t.Error("signature decoding failed")
	}
}
