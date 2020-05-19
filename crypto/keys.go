package crypto

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// --- utils ---

func ReadKeys(keysfile string) KeypairAll {

	log.Println("read keys from ", keysfile)
	dat, _ := ioutil.ReadFile(keysfile)
	s := string(dat)
	// vs, _ := parser.ReadMapP(s)
	// privHex := parser.StringUnWrap(vs[0])
	// pubkeyHex := parser.StringUnWrap(vs[1])
	//log.Println("pub ", pubkeyHex)

	var kpa KeypairAll
	json.Unmarshal([]byte(s), &kpa)
	return kpa
}

// func CreateKeypairFormat(privkey string, pubkey_string string, address string) string {
// 	mp := map[string]string{"privkey": parser.StringWrap(privkey), "pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
// 	m := parser.MakeMap(mp)
// 	return m
// }

// func CreatePubKeypairFormat(pubkey_string string, address string) string {
// 	mp := map[string]string{"pubkey": parser.StringWrap(pubkey_string), "address": parser.StringWrap(address)}
// 	m := parser.MakeMap(mp)
// 	return m
// }

//write keys to file with format
func WriteKeys(kp Keypair, keysfile string) {
	pubkeyHex := PubKeyToHex(kp.PubKey)
	privHex := PrivKeyToHex(kp.PrivKey)
	address := Address(pubkeyHex)
	kpa := KeypairAll{PubKey: PubKeyFromHex(pubkeyHex), PrivKey: PrivKeyFromHex(privHex), Address: address}
	//s := CreateKeypairFormat(privHex, pubkeyHex, address)
	s, _ := json.Marshal(kpa)
	ioutil.WriteFile(keysfile, []byte(s), 0644)
}

func WritePubKeys(kp Keypair, keysfile string) {
	pubkeyHex := PubKeyToHex(kp.PubKey)
	address := Address(pubkeyHex)
	kpp := KeypairPub{PubKey: PubKeyFromHex(pubkeyHex), Address: address}
	//s := CreatePubKeypairFormat(pubkeyHex, address)
	s, _ := json.Marshal(kpp)
	ioutil.WriteFile(keysfile, []byte(s), 0644)
}
