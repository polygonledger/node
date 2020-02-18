package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

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
	return PubKeyToHex(someKey)
}

func PubFromSecret(secret string) btcec.PublicKey {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedsecret := hex.EncodeToString(hasher.Sum(nil))

	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), []byte(hashedsecret))
	return *pubKey
}

func MsgHash(message string) []byte {
	messageHash := chainhash.DoubleHashB([]byte(message))
	return messageHash
}

func Sign(keypair Keypair, message string) btcec.Signature {
	messageHash := MsgHash(message)
	signature, err := keypair.PrivKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		//return nil
	}
	//fmt.Println("signature ", signature)
	return *signature

}
