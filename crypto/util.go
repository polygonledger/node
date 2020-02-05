package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcd/btcec"
)

//var hash = sha256("secret")
//var keypair = MakeKeypair(hash)

//type Keypair struct {
//public, private

func PubHexFromSecret() string {
	someKey := PubFromSecret()
	//!needs checking
	return string(hex.EncodeToString(someKey.SerializeCompressed()))
}

func PubFromSecret() btcec.PublicKey {
	secret := "secret"
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedsecret := hex.EncodeToString(hasher.Sum(nil))

	//hashedsecret := sha256.Sum256([32]byte("secret"))
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), []byte(hashedsecret))
	return *pubKey
}

func RanAddress() *btcec.PublicKey {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f8720ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return pubKey

}

func KeyExample() {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	log.Println("pubkey example %v", pubKey)
	log.Println(privKey)

	//hash := sha256.Sum256(pubKey.Serialize())

	// Sign a message using the private key.
	// message := "test message"
	// messageHash := chainhash.DoubleHashB([]byte(message))
	// signature, err := privKey.Sign(messageHash)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Serialize and display the signature.
	// fmt.Printf("Serialized Signature: %x\n", signature.Serialize())
	// // Verify the signature for the message using the public key.
	// verified := signature.Verify(messageHash, pubKey)
	// fmt.Printf("Signature Verified? %v\n", verified)

	// data := []byte("hello")
	// hash := sha256.Sum256(data)
	// fmt.Printf("%x", hash[:])

	// timestamp := time.Now().Unix()
	// b := []byte(string(timestamp))
	// hash = sha256.Sum256(b)

	// fmt.Printf("\n%x", hash[:])
}
