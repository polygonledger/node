package main

import (
	"testing"

	"github.com/polygonledger/node/block"
	chain "github.com/polygonledger/node/chain"
)

func TestLoad(t *testing.T) {

	mgr := chain.CreateManager()

	if mgr.BlockHeight() != 0 {
		t.Error("BlockHeight")
	}

	b := block.Block{}
	mgr.AppendBlock(b)

	//after appending block height increases by 1
	if mgr.BlockHeight() != 1 {
		t.Error("BlockHeight")
	}

	mgr.WriteChain()

	mgr.ReadChain()

	if mgr.BlockHeight() != 1 {
		t.Error("BlockHeight")
	}

}
