//package main

//client network

// func requestreply(ntchan ntwk.Ntchan, req_msg string) {

// 	//TODO! use readloop and REQ/REP chans

// 	log.Println("requestreply >> ", req_msg)
// 	//REQUEST
// 	ntchan.REQ_out <- req_msg
// 	//REPLY

// 	resp_string := <-ntchan.REP_in
// 	log.Println("REP_in >> ", resp_string)

// 	// msg := ntwk.ParseMessage(resp_string)
// 	// log.Println("response ", msg.MessageType)
// 	// if msg.MessageType == ntwk.REP {
// 	// 	//need to match to know this is the same request ID?
// 	// 	log.Println("REPLY ", msg)
// 	// }
// }
