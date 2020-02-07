package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/btcd/btcec"
	"github.com/btcd/chaincfg/chainhash"
	"github.com/polygonledger/node/block"
)

//var hash = sha256("secret")
//var keypair = MakeKeypair(hash)

//type Keypair struct {
//public, private

type Keypair struct {
	PrivKey btcec.PrivateKey
	PubKey  btcec.PublicKey
}

func Address(pubkey string) string {
	/*var private = {};
	private.getAddressByPublicKey = function (publicKey) {
		var publicKeyHash = crypto.createHash('sha256').update(publicKey, 'hex').digest();
		var temp = new Buffer(8);
		for (var i = 0; i < 8; i++) {
			temp[i] = publicKeyHash[7 - i];
		}

		var address = bignum.fromBuffer(temp).toString() + "C";
		return address;
	}*/
	return "P" + GetSHAHash(pubkey)[:12]
}

func SignTx(tx block.Tx, keypair Keypair) btcec.Signature {
	//message := strconv.Itoa(tx.Id)
	message := fmt.Sprintf("%d", tx.Id)

	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := keypair.PrivKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println("signature ", signature)
	return *signature

	//verified := signature.Verify(messageHash, &keypair.PubKey)
	//fmt.Printf("Signature Verified? %v\n", verified)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetSHAHash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RandomPublicKey() string { //type key
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(10000)
	somePubKey := PubFromSecret("secret" + strconv.Itoa(randNonce))
	//!needs checking/fixing
	return GetMD5Hash(string(hex.EncodeToString(somePubKey.SerializeCompressed())))
}

func PubHexFromSecret() string {
	someKey := PubFromSecret("secret")
	//!needs checking
	return string(hex.EncodeToString(someKey.SerializeCompressed()))
}

func PubFromSecret(secret string) btcec.PublicKey {
	//secret := "secret"
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

func SomeKeypair() Keypair {
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		//return nil
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	kp := Keypair{PrivKey: *privKey, PubKey: *pubKey}
	return kp
}

func SignExample(keypair Keypair) {

	//privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := keypair.PrivKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("signature ", signature)

	verified := signature.Verify(messageHash, &keypair.PubKey)
	fmt.Printf("Signature Verified? %v\n", verified)
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
