package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/polygonledger/edn"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/netio"
	"github.com/polygonledger/node/parser"
)

//basic block functions
func TestBasicParse(t *testing.T) {

	txmap := parser.CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)
	sm := `{:amount 42 :receiver "P7ba453f23337" :sender "Pa033f6528cc1"}`
	if txmap != sm {
		t.Error("create txmap fail ", txmap, sm)
	}

	pubk := "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31"
	sig := "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"
	sigmap := parser.CreateSigmap(pubk, sig)

	if sigmap != `{:senderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}` {
		t.Error("sigmap fail ", sigmap)
	}

}

func TestTxassemble(t *testing.T) {
	simpletx := parser.CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)

	keypair := crypto.PairFromSecret("test")
	sigmap := parser.SignMap(keypair, simpletx)
	v := parser.TxVector(simpletx, sigmap)

	s := `[:STX {:amount 42 :receiver "P7ba453f23337" :sender "Pa033f6528cc1"} {:senderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :signature "30450221008ef704458815e7318ba5e161e2b11dcfa446e00146ebe7d4beecd4c3f812105002201438f2d12aae0cf391f0e6893243b48efa27b77c551603fc956954faa742a923"}]`
	if v != s {
		t.Error("tx vector not proper ", v)
		t.Error(s)
	}
	valid := parser.VerifyTxScriptSig(v)

	if !valid {
		t.Error("tx not valid")
	}

}

func TestReadMap(t *testing.T) {

	ms := "{:mykey bla :second bar :third abc}"
	vs, ks := parser.ReadMapP(ms)

	h := []string{"bla bar abc"}
	if reflect.DeepEqual(vs, h) {
		t.Error("scan map")
	}

	h2 := []string{"mykey second third"}
	if reflect.DeepEqual(ks, h2) {
		t.Error("scan map")
	}

}

func TestTxStore(t *testing.T) {
	simpletx := parser.CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)

	keypair := crypto.PairFromSecret("test")
	sigmap := parser.SignMap(keypair, simpletx)
	v := parser.TxVector(simpletx, sigmap)

	testtx := "test.tx"

	ioutil.WriteFile(testtx, []byte(v), 0644)

	// dat, _ := ioutil.ReadFile(testtx)

	// txtype, sigmap, txmap := parser.ScanScript(string(dat))
	// if txtype == nil {

	// }
	// fmt.Println(sigmap)
	// fmt.Println(txmap)
	// fmt.Println(txtype)

	os.Remove(testtx)

}

// func TestMsgBasic(t *testing.T) {
// 	msgstring := "{:REP X :data 0}"
// 	msg := netio.ParseMessageMapData(msgstring)
// 	if msg.MessageType != "REP" {
// 		t.Error("parse msg MessageType", msg)
// 	}
// 	if msg.Command != "X" {
// 		t.Error("parse msg ", msg.Command)
// 	}
// 	arr := []byte("0")
// 	if !reflect.DeepEqual(msg.Data, arr) {
// 		t.Error("parse msg data ", msg.Data)
// 	}

// }

func TestAccmsg(t *testing.T) {

	msg := netio.Message{MessageType: "REQ", Command: "Test", Data: []byte("000")}
	b, _ := edn.Marshal(msg)
	arr := []byte("{:messageType\"REQ\":command\"Test\":data #base64\"MDAw\"}")

	if !reflect.DeepEqual(b, arr) {
		t.Error("??", string(b))
	}

	//FIX
	// msgstring := "{:REP ACCOUNTS :data {\"P06d058383389\"20}}"
	// msg := netio.ParseMessageMapData(msgstring)
	// if msg.MessageType != "REP" {
	// 	t.Error("parse msg ", msg.MessageType)
	// }
	// if msg.Command != "ACCOUNTS" {
	// 	t.Error("parse msg ", msg.Command)
	// }
	// tb := []byte("{\"P06d058383389\"20}")
	// if !reflect.DeepEqual(msg.Data, tb) {
	// 	t.Error("data ", msg.Data, tb)
	// }

}

// func TestTxSimpleStruct(t *testing.T) {

// 	data := `{:amount 10 :sender "abc" :receiver "xyz"}`
// 	var tx block.SimpleTx
// 	edn.Unmarshal([]byte(data), &tx)

// 	if tx.Amount != 10 {
// 		t.Error("wrong amount")
// 	}

// 	d2 := `{:senderPubkey "abc" :signature "xyz"}`
// 	var txs block.TxSigmap
// 	edn.Unmarshal([]byte(d2), &txs)
// 	if txs.SenderPubkey != "abc" || txs.Signature != "xyz" {
// 		t.Error("senderpubkey")
// 	}

// 	d3 := `{:TxType "STX" :TxTransfer {:amount 10 :sender "abc" :receiver "xyz"} :Sigmap {:senderPubkey "abc" :signature "xyz"}}`
// 	var txe block.TxExpr
// 	edn.Unmarshal([]byte(d3), &txe)
// 	if txe.Transfer.Amount != 10 {
// 		t.Error("amount")
// 	}

// 	var txe2 block.TxExpr
// 	txe2.TxType = "STX"
// 	txe2.Transfer.Amount = 5
// 	txe2.Transfer.Sender = "abc"
// 	txe2.Transfer.Receiver = "xyz"
// 	txe2.Sigmap.SenderPubkey = "zzz"
// 	txe2.Sigmap.Signature = "aaa"
// 	b, _ := edn.Marshal(txe2)
// 	fmt.Println(string(b))

// }
