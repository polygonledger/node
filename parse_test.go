package main

import (
	"testing"

	"github.com/polygonledger/node/crypto"
)

//basic block functions
func TestBasicParse(t *testing.T) {

	txmap := CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)

	if txmap != `{:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42}` {
		t.Error("create txmap fail ", txmap)
	}

	pubk := "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31"
	sig := "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"
	sigmap := CreateSigmap(pubk, sig)

	if sigmap != `{:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}` {
		t.Error("sigmap fail ", sigmap)
	}

}

func TestTxassemble(t *testing.T) {
	simpletx := CreateSimpleTxContent("Pa033f6528cc1", "P7ba453f23337", 42)

	keypair := crypto.PairFromSecret("test")
	sigmap := SignMap(keypair, simpletx)
	v := txVector(simpletx, sigmap)

	if v != `[:STX {:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42} {:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}]` {
		t.Error("tx vector not proper ", v)
	}
	valid := VerifyTxScriptSig(v)

	if !valid {
		t.Error("tx not valid")
	}

}
