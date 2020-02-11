package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cryptoutil "github.com/polygonledger/node/crypto"
)

func writeKeys(kp cryptoutil.Keypair, keysfile string) {

	pubkeyHex := cryptoutil.PubKeyToHex(kp.PubKey)
	log.Println("pub ", pubkeyHex)

	privHex := cryptoutil.PrivKeyToHex(kp.PrivKey)
	log.Println("privHex ", privHex)

	t := pubkeyHex + "\n" + privHex
	ioutil.WriteFile(keysfile, []byte(t), 0644)
}

func main() {

	log.Println("create keypair")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	pw, _ := reader.ReadString('\n')
	fmt.Println(pw)

	//check if exists
	//dat, _ := ioutil.ReadFile("keys.txt")
	//check(err)
	//fmt.Print("keys ", string(dat))

	kp := cryptoutil.PairFromSecret(pw)
	log.Println("keypair ", kp)

	writeKeys(kp, "keys.txt")

}
