package main

//keys.txt layout
//pubkey\nprivkey

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	cryptoutil "github.com/polygonledger/node/crypto"
)

func readKeys(keysfile string) cryptoutil.Keypair {

	dat, _ := ioutil.ReadFile(keysfile)
	s := strings.Split(string(dat), string("\n"))

	pubkeyHex := s[0]
	log.Println("pub ", pubkeyHex)

	privHex := s[1]
	log.Println("privHex ", privHex)

	return cryptoutil.Keypair{PubKey: cryptoutil.PubKeyFromHex(pubkeyHex), PrivKey: cryptoutil.PrivKeyFromHex(privHex)}
}

func writeKeys(kp cryptoutil.Keypair, keysfile string) {

	pubkeyHex := cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println("pub ", pubkeyHex)

	privHex := cryptoutil.PrivKeyToHex(kp.PrivKey)
	log.Println("privHex ", privHex)

	t := pubkeyHex + "\n" + privHex
	ioutil.WriteFile(keysfile, []byte(t), 0644)
}

func createKeys() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.Trim(pw, string('\n'))
	fmt.Println(pw)

	//check if exists
	//dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)

	kp := cryptoutil.PairFromSecret(pw)
	log.Println("keypair ", kp)

	writeKeys(kp, "keys.txt")

}

func main() {

	optionPtr := flag.String("option", "createkeys", "createkeys files")
	flag.Parse()
	fmt.Println("option:", *optionPtr)

	if *optionPtr == "createkeys" {
		createKeys()
	} else if *optionPtr == "readkeys" {
		kp := readKeys("keys.txt")
		log.Println(kp)
	} else if *optionPtr == "sign" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to sign: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		fmt.Println(msg)
		kp := readKeys("keys.txt")
		signature := cryptoutil.SignMsgHash(kp, msg)
		log.Println("signature ", signature)

		sighex := hex.EncodeToString(signature.Serialize())
		log.Println("sighex ", sighex)
	} else if *optionPtr == "verify" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message to verify: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.Trim(msg, string('\n'))
		fmt.Println(msg)

		fmt.Print("Enter signature to verify: ")
		msgsig, _ := reader.ReadString('\n')
		msgsig = strings.Trim(msgsig, string('\n'))
		//x := "3045022100dd2781cc37edb84c5ed21b3d8fc03d49ebddf5647d23a9132eeea8bd2b951bd1022041519c47b77803d528d1b428ccb4d84a90ce3b67a22662d5feaa84c4521e5759"
		//fmt.Println("??", msgsig[0], len(msgsig), len(x))

		sign := cryptoutil.SignatureFromHex(msgsig)

		fmt.Print("Enter pubkey to verify: ")
		msgpub, _ := reader.ReadString('\n')
		fmt.Println(msgpub)
		msgpub = strings.Trim(msgpub, string('\n'))

		//msgpub = "039f6095ba1afa34c437a88fceb444bf177326eb9222d4938336387ecb4cbe7234"

		pubkey := cryptoutil.PubKeyFromHex(msgpub)

		verified := cryptoutil.VerifyMessageSignPub(sign, pubkey, msg)
		log.Println("verified ", verified)

	}

}
