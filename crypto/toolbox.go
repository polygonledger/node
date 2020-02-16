package crypto

//general crypto toolbox
//higher level functions building on btcec

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcd/btcec"
	"github.com/btcd/chaincfg/chainhash"
	"github.com/polygonledger/node/block"
)

type Keypair struct {
	PrivKey btcec.PrivateKey
	PubKey  btcec.PublicKey
}

//TODO only from pubkey type
func Address(pubkey string) string {
	return "P" + GetSHAHash(pubkey)[:12]
}

func PairFromHex(hexString string) Keypair {
	pkBytes, _ := hex.DecodeString(hexString)
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	kp := Keypair{PrivKey: *privKey, PubKey: *pubKey}
	return kp
}

func PairFromSecret(secret string) Keypair {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedsecret := hex.EncodeToString(hasher.Sum(nil))

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), []byte(hashedsecret))
	kp := Keypair{PrivKey: *privKey, PubKey: *pubKey}
	return kp
}

func PrivKeyToHex(privkey btcec.PrivateKey) string {
	return hex.EncodeToString(privkey.Serialize())
}

func PrivKeyFromHex(hexString string) btcec.PrivateKey {
	//TODO handle errors
	pkBytes, _ := hex.DecodeString(hexString)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return *privKey
}

func PubKeyToHex(pubkey btcec.PublicKey) string {
	return string(hex.EncodeToString(pubkey.SerializeCompressed()))
}

func PubKeyFromHex(hexString string) btcec.PublicKey {
	pubKeyBytes, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println(err)
		//return
	}
	pubKey, _ := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	return *pubKey
}

// Decode hex-encoded serialized signature
func SignatureFromHex(hexString string) btcec.Signature {
	//TODO handle errors
	sigBytes, err := hex.DecodeString(hexString)

	if err != nil {
		fmt.Println(err)
		//return
	}
	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		//return
	}
	return *signature
}

func VerifyMessageSignPub(signature btcec.Signature, pubkey btcec.PublicKey, message string) bool {

	messageHash := MsgHash(message)
	verified := signature.Verify(messageHash, &pubkey)
	return verified
}

func VerifyMessageSign(signature btcec.Signature, keypair Keypair, message string) bool {

	messageHash := MsgHash(message)
	verified := signature.Verify(messageHash, &keypair.PubKey)
	//log.Println("?? ", message, verified)
	return verified
}

func SignMsgHash(keypair Keypair, message string) btcec.Signature {
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := keypair.PrivKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		//return
	}
	return *signature
}

func SignTx(tx block.Tx, keypair Keypair) btcec.Signature {
	//TODO sign tx not just id
	txJson, _ := json.Marshal(tx)
	//log.Println(string(txJson))
	//message := fmt.Sprintf("%d", tx.Id)

	messageHash := chainhash.DoubleHashB([]byte(txJson))
	signature, err := keypair.PrivKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		//return
	}
	//fmt.Println("signature ", signature)
	return *signature

	//verified := signature.Verify(messageHash, &keypair.PubKey)
	//fmt.Printf("Signature Verified? %v\n", verified)
}

func RemoveSigTx(tx block.Tx) block.Tx {
	tx.Signature = ""
	return tx
}

func RemovePubTx(tx block.Tx) block.Tx {
	tx.SenderPubkey = ""
	return tx
}

//TODO
func VerifyTxSig(tx block.Tx) bool {
	getpubkey := PubKeyFromHex(tx.SenderPubkey)
	gotsighex := tx.Signature
	sign := SignatureFromHex(gotsighex)
	//need to remove sig and pubkey for validation
	tx = RemoveSigTx(tx)
	tx = RemovePubTx(tx)

	txJson, _ := json.Marshal(tx)
	verified := VerifyMessageSignPub(sign, getpubkey, string(txJson))
	return verified

}
