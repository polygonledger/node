package main

import (
	"testing"
	"time"

	"github.com/polygonledger/node/ntwk"
)

func TestBasicNtwk(t *testing.T) {

	var ntchan ntwk.Ntchan
	go ntwk.ReadProcessor(ntchan, 1*time.Millisecond)

	if len(ntchan.Reader_queue) != 0 {
		t.Error("reader queue not empty")
	}

}
