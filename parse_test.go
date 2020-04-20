package main

import (
	"reflect"
	"testing"

	"github.com/polygonledger/node/crypto"
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

	if v != `[:STX {:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42} {:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}]` {
		t.Error("tx vector not proper ", v)
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
