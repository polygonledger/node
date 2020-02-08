package net

//generic message
type Message struct {
	MessageType string //Request <--> Reply
	Command     string //Specific message command
	//Data        []byte
	//Signature       btcec.Signature
}
