package main

import (
	"fmt"
	"testing"
)

//basic block functions
func TestBasicParse(t *testing.T) {
	txmap := `{:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42}`

	sigmap := `{:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}`

	verified := verifySigmap(sigmap, txmap)

	if !verified {
		t.Error("not verify")
	}

	v := txVector(txmap, sigmap)

	if v != `[:STX {:Sender "Pa033f6528cc1" :Receiver "P7ba453f23337" :amount 42} {:SenderPubkey "03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31" :Signature "304502210086d04e9613514174e75558ea4e7fd96e691e87b5deed39b4da3d6774e1ffe81b02202e63019ad59b7cd42dbeacfe9b1a7b05a421f72705d4659aea6b0450db638b96"}]` {
		fmt.Println(v)
		//t.Error(v)
	}

	sigmap1, txmap1 := ScanScript(v)
	if sigmap != sigmap1 {
		t.Error("fail sigmap ", sigmap, sigmap1)
	}
	if txmap != txmap1 {
		t.Error("fail txmap ", sigmap1)
	}
}
