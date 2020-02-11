package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cryptoutil "github.com/polygonledger/node/crypto"
)

func main() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	fmt.Println(pw)

	//if exists
	//dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)
	//fmt.Print("keys ", string(dat))

	kp := cryptoutil.PairFromSecret(pw)
	log.Println("keypair ", kp)

	pubkeyHex := cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println("pub ", pubkeyHex)

	privHex := cryptoutil.PrivKeyToHex(kp.PrivKey)
	log.Println("privHex ", privHex)

	t := pubkeyHex + "\n" + privHex
	ioutil.WriteFile("keys.txt", []byte(t), 0644)
}
