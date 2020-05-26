package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/polygonledger/node/block"
)

func TestBasicSign(t *testing.T) {

	//sign newly created keypair should be valid signature
	keypair := PairFromSecret("test")
	message := "test"

	signature := SignMsgHash(keypair.PrivKey, message)
	verified := VerifyMessageSign(signature, keypair, message)
	if !verified {
		t.Error("msg failed")
	}

	messagefalse := "testshouldbefalse"
	verifiedfalse := VerifyMessageSign(signature, keypair, messagefalse)

	if verifiedfalse {
		t.Error("sign verify should fail")
	}

	otherkeypair := PairFromSecret("testother")
	verifiedother := VerifyMessageSign(signature, otherkeypair, message)
	if verifiedother {
		t.Error("sign other should fail")
	}

}

func TestDecode(t *testing.T) {

	pubKey := PubKeyFromHex("02a673638cb9587cb68ea08dbef685c6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")

	h := "30450220090ebfb3690a0ff115bb1b38b8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e7071cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479"
	signature := SignatureFromHex(h)

	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, &pubKey)

	if !verified {
		t.Error("signature decoding failed")
	}
}

func TestAddress(t *testing.T) {

	keypair := PairFromSecret("test")
	pubkey_string := PubKeyToHex(keypair.PubKey)
	if pubkey_string != "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" {
		t.Error("expected different hex of pubkey")
	}

	hexString := "a11b0a4e1a132305652ee7a8eb7848f6ad5ea381e3ce20a2c086a2e388230811"
	privKey := PrivKeyFromHex(hexString)
	privKeyHex := PrivKeyToHex(privKey)

	if privKeyHex != hexString {
		t.Error("privkey encoding")
	}

	addr := Address(pubkey_string)
	if addr[0] != 'P' {
		t.Error("address should start with P ", addr[0])
	}

	if len(addr) != 13 {
		t.Error("length of address should be 13 ", len(addr))
	}
}

func TestSignHardcoded(t *testing.T) {
	pub := "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31"
	pubkey := PubKeyFromHex(pub)

	keypair := PairFromSecret("test")
	h := PubKeyToHex(keypair.PubKey)

	if h != pub {
		t.Error("hardcoded pubkey wrong")
	}

	sig := "3045022100b6789781f032512fc9ae06e9621ca4ce40d317a83a6b6f4ee1ad35942a3c928602204d8f03b330efc416b822447862333140d0acb3323d4575f1efba6e5418a036f7"
	sign := SignatureFromHex(sig)
	msg := "test"
	verified := VerifyMessageSignPub(sign, pubkey, msg)
	if !verified {
		t.Error("should verify standard ", verified)
	}

}

func TestGenkeys(t *testing.T) {
	h := "22a47fa09a223f2aa079edf85a7c2d4f8720ee63e502ee2869afab7de234b80c"

	keypair := PairFromHex(h)

	if PubKeyToHex(keypair.PubKey) == "" {
		t.Error("keypair is nil")
	}

}

func TestSignTx(t *testing.T) {
	//sign
	keypair := PairFromSecret("test")
	var tx block.Tx
	//s := block.AccountFromString("Pa033f6528cc1")
	s := "Pa033f6528cc1"
	r := s //TODO
	tx = block.Tx{Nonce: 0, Amount: 0, Sender: s, Receiver: r}

	signature := SignTx(tx, keypair.PrivKey)
	sighex := hex.EncodeToString(signature.Serialize())

	if sighex == "" {
		t.Error("hex empty")
	}
	tx.Signature = sighex
	tx.SenderPubkey = PubKeyToHex(keypair.PubKey)

	//verify
	verified := VerifyTxSig(tx)

	if !verified {
		t.Error("verify tx fail")
	}

}

func TestAdress(t *testing.T) {
	//Pa033f6528cc1
	keypair := PairFromSecret("test")
	pub := PubKeyToHex(keypair.PubKey)

	//pub := "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31"
	a := Address(pub)

	if a != "Pa033f6528cc1" {
		t.Error("hardcoded wrong")
	}

	if a[0] != 'P' {
		t.Error("adress stars with P")
	}

}
