package block

type Tx struct {
	Id     [32]byte //gets assigned when verified in a block
	Nonce  int
	Amount int
	//sender
	//receiver

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}
