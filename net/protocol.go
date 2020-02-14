package net

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	block "github.com/polygonledger/node/block"
	cryptoutil "github.com/polygonledger/node/crypto"
)

const (
	// Port is the port number that the server listens to.
	Server_address string = "127.0.0.1"
	Port                  = ":8888"
	//TODO move to message type
	CMD_PING                  = "PING"
	CMD_TX                    = "TX"
	CMD_RANDOM_ACCOUNT        = "RANACC"
	CMD_BALANCE               = "BALANCE" //get balance of account
	Genesis_Address    string = "P0614579c42f2"
	DELIM              byte   = '|'
	DELIM_HEAD         byte   = '#'
	EMPTY_MSG                 = ""
	ERROR_READ                = "error_read"
)

//---- protocol layer ------
//given a sream read from it
//TODO proper error handling
func ReadStream(rw *bufio.ReadWriter) string {
	msg, err := rw.ReadString(DELIM)
	if err != nil {
		//issue
		//special case is empty message if client disconnects?
		if len(msg) == 0 {
			log.Println("empty message")
			return EMPTY_MSG
		} else {
			log.Println("Failed ", err)
			//log.Println(err.)
			return ERROR_READ
		}
	}
	return msg
}

//we have a message string and will parse it to a message struct
//delimiters of two kind:
//* DELIM for delimiting the entire message
//* DELIM_
//currently we employ delimiters instead of byte encoding, so the size of messages is unlimited
//can however easily fix by adding size to header and reject messages larger than maximum size
func ParseMessage(msgString string) Message {
	msgString = strings.Trim(msgString, string(DELIM))
	s := strings.Split(msgString, string(DELIM_HEAD))

	//ERROR handling of malformed messages

	var msg Message
	msg.MessageType = s[0]
	msg.Command = s[1]
	dataJson := s[2] //data can empty but still we expect the delim to be there

	msg.Data = []byte(dataJson)
	fmt.Println(msg)
	return msg
}

func ReadMsg(rw *bufio.ReadWriter) string {
	//TODO handle err
	msg, _ := rw.ReadString(DELIM)
	msg = strings.Trim(msg, string(DELIM))
	return msg
}

func EncodeReply(resp string) string {
	//TODO! header missing
	response := resp + string(DELIM)
	return response
}

func EncodeMessageTx(txJson []byte) string {
	emptyData := ""
	msgCmd := "TX"
	msg := REQ + string(DELIM_HEAD) + msgCmd + string(DELIM_HEAD) + string(txJson) + string(DELIM_HEAD) + emptyData + string(DELIM)
	return msg
}

//Faucet => send fixed number of coins to specified address

//getBlocks
//registerPeer

//pickRandomAccount

//storeBalance

//handlers TODO this is higher level and should be somewhere else

func RandomTx(account_s block.Account) block.Tx {
	// s := cryptoutil.RandomPublicKey()
	// address_s := cryptoutil.Address(s)
	// account_s := block.AccountFromString(address_s)
	// log.Printf("%s", s)

	//FIX
	//doesn't work on client side
	//account_r := chain.RandomAccount()

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)

	kp := cryptoutil.PairFromSecret("test111??")
	log.Println("PUBKEY ", kp.PubKey)

	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	account_r := block.AccountFromString(address_r)

	//TODO make sure the amount is covered by sender
	rand.Seed(time.Now().UnixNano())
	randomAmount := rand.Intn(20)

	log.Printf("randomAmount ", randomAmount)
	log.Printf("randNonce ", randNonce)
	testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	testTx.Signature = cryptoutil.SignTx(testTx, kp)
	log.Println(">> ran tx", testTx.Signature)
	return testTx
}

// func ReadPipe(rw *bufio.ReadWriter) {

// }

func RequestReply(rw *bufio.ReadWriter, req_msg string) string {
	//REQUEST
	WritePipe(rw, req_msg)

	//REPLY
	resp_msg := ReadMsg(rw)

	return resp_msg
}

func WritePipe(rw *bufio.ReadWriter, message string) error {
	n, err := rw.WriteString(message)
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	} else {
		log.Println(strconv.Itoa(n) + " bytes written")
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

func ConstructMessage(cmd string) string {
	//delim := "\n"
	msg := cmd + string(DELIM)
	return msg
}

/*
request account address
*/
func RequestTest(rw *bufio.ReadWriter) error {
	log.Println("RequestTest")
	a := Message{Command: "test"}
	log.Printf("msg: \n%#v\n", a)
	enc := gob.NewEncoder(rw)
	err := enc.Encode(a)
	rw.WriteString(string(DELIM))
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed")
	}
	return err
}

/*
request account address
*/
func RequestAccount(rw *bufio.ReadWriter) error {
	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)
	msg := ConstructMessage(CMD_RANDOM_ACCOUNT)
	error := WritePipe(rw, msg)
	return error
}

/*
request account address
*/
func ReceiveAccount(rw *bufio.ReadWriter) error {
	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)
	return nil
}

/*
send message
*/
func SendMessage(rw *bufio.ReadWriter) error {

	//Command
	msg := "data" + string(DELIM)
	//log.Println("?? ", msg)
	n, err := rw.WriteString(msg)
	if err != nil {
		return errors.Wrap(err, "Could not write data ("+strconv.Itoa(n)+" bytes written)")
	}

	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed")
	}
	return nil
}

/*
sends account address
*/
func SendAccount(rw *bufio.ReadWriter) error {

	a := block.Account{AccountKey: "test"}

	log.Printf("account: \n%#v\n", a)

	enc := gob.NewEncoder(rw)
	//Command

	msg := CMD_RANDOM_ACCOUNT + string(DELIM)
	//log.Println("?? ", msg)
	n, err := rw.WriteString(msg)
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(a)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", a)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed")
	}
	return nil
}
