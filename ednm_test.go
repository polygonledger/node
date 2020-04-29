package main

import (
	"testing"
	"fmt"

	"github.com/polygonledger/node/parser"
)

//basic block functions
func TestMap(t *testing.T) {

	m := map[string]string{"test": "value"}

	fmt.Println("????? ")
	mstr := parser.MakeMap(m)

	if mstr != "{:test value}" {
		t.Error("error creating map")
	}

}
