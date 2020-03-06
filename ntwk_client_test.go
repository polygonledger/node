package main

import (
	"testing"
	"time"

	"github.com/polygonledger/node/client"
	"github.com/polygonledger/node/ntwk"
)

func TestClientServer_Run(t *testing.T) {

	testsrv := initserver()
	defer testsrv.Close()

	testclient := testclient()
	testclient := client.initClient()

	req_msg := ntwk.EncodeMsgString(ntwk.REQ, ntwk.CMD_PING, "")
	testclient.REQ_out <- req_msg

	time.Sleep(800 * time.Millisecond)

	//log.Println("TestServer_Run > ", testsrv, testsrv.Peers)

	if !testsrv.accepting {
		t.Error("not accepting")
	}

}
