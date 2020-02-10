package main

import (
	"testing"

	"github.com/btcd/chaincfg/chainhash"
	cryptoutil "github.com/polygonledger/node/crypto"
)

func TestBasicSign(t *testing.T) {

	//sign newly created keypair should be valid signature
	keypair := cryptoutil.PairFromSecret("test")
	message := "test"

	signature := cryptoutil.Sign(keypair, message)

	messageHash := cryptoutil.MsgHash(message)
	verified := signature.Verify(messageHash, &keypair.PubKey)

	if !verified {
		t.Error("msg failed")
	}

}

func TestRecordPrivkey(t *testing.T) {
	// Decode the hex-encoded pubkey of the recipient.
	// pubKeyBytes, err := hex.DecodeString("04115c42e757b2efb7671c578530ec191a1" +
	// 	"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
	// 	"21010db67ac11b1b51b651953d2") // uncompressed pubkey

	// pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())

	// // Encrypt a message decryptable by the private key corresponding to pubKey
	// // message := "test message"
	// // ciphertext, err := btcec.Encrypt(pubKey, []byte(message))

	// // Decode the hex-encoded private key
	// pkBytes, err := hex.DecodeString("a11b0a4e1a132305652ee7a8eb7848f6ad" +
	// 	"5ea381e3ce20a2c086a2e388230811")

	// note that we already have corresponding pubKey
	//privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	// Try decrypting and verify if it's the same message
	// plaintext, err := btcec.Decrypt(privKey, ciphertext)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(string(plaintext))
}

func TestDecode(t *testing.T) {

	pubKey := cryptoutil.PubKeyFromHex("02a673638cb9587cb68ea08dbef685c6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")

	h := "30450220090ebfb3690a0ff115bb1b38b8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e7071cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479"
	signature := cryptoutil.SignatureFromHex(h)

	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, &pubKey)

	if !verified {
		t.Error("signature decoding failed")
	}
}

func TestAddress(t *testing.T) {

	keypair := cryptoutil.PairFromSecret("test")
	pubkey_string := cryptoutil.PubKeyToHex(keypair.PubKey)
	if pubkey_string != "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" {
		t.Error("expected different hex of pubkey")
	}

	hexString := "a11b0a4e1a132305652ee7a8eb7848f6ad5ea381e3ce20a2c086a2e388230811"
	privKey := cryptoutil.PrivKeyFromHex(hexString)
	privKeyHex := cryptoutil.PrivKeyToHex(privKey)

	if privKeyHex != hexString {
		t.Error("privkey encoding")
	}

	// pubKeyBytes, err := hex.DecodeString(pubkey_string)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	//FAIL
	// if *pubKey != keypair.PubKey {
	// 	log.Println(*pubKey, keypair.PubKey)
	// 	t.Error("error recoding pubkey")
	// }

	addr := cryptoutil.Address(pubkey_string)
	if addr[0] != 'P' {
		t.Error("address should start with P ", addr[0])
	}

	if len(addr) != 13 {
		t.Error("length of address should be 13 ", len(addr))
	}
}
