package chain

//TODO

import (
	"crypto/sha256"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/polygon/block"
	cryptoutil "github.com/polygon/crypto"
)

//TODO fix circular import
const (
	Genesis_Address string = "P0614579c42f2"
)

var Tx_pool []block.Tx
var Blocks []block.Block
var Latest_block block.Block
var Accounts map[block.Account]int

//testing
func InitAccounts() {
	Accounts = make(map[block.Account]int)
	log.Println("init accounts %i", len(Accounts))
	//Genesis_Account := block.AccountFromString(Genesis_Address)
	//set genesiss account
	SetAccount(block.AccountFromString(Genesis_Address), 200)
}

//handlers

func txHash(tx block.Tx) [32]byte {
	b := []byte(string(tx.Nonce)[:])
	hash := sha256.Sum256(b)
	return hash
}

func HandleTx(tx block.Tx) {
	//TODO its own function txhash
	//hash of timestamp is same, check lenght of bytes used??
	timestamp := time.Now().Unix()
	//b := []byte(append(string(timestamp)[:], string(tx.Nonce)[:]))
	tx.Id = txHash(tx)

	log.Println("tx id %x time %s", tx.Id, timestamp)

	//hash := sha256.Sum256(xdata)

	//tx.Id = hash

	//TODO its own function
	if txValid(tx) {
		Tx_pool = append(Tx_pool, tx)
	} else {
		log.Println("invalid tx")
	}
	//blockheight += 1

	//log.Printf("tx_pool_size: \n%d", tx_pool_size)

	log.Printf("tx_pool size: %d\n", len(Tx_pool))
}

//#### blockchain functions

//empty the tx pool
func EmptyPool() {
	Tx_pool = []block.Tx{}
}

func blockHash(block block.Block) block.Block {
	//FIX hash of proper data, merkle and such
	timeFormat := "2020-02-02 16:06:06"
	new_hash := []byte(string(block.Timestamp.Format(timeFormat))[:])
	blockheight_string := []byte(strconv.FormatInt(int64(block.Height), 10))
	//new_hash = append(new_hash, blockheight_string)
	log.Printf("newhash %x", new_hash)
	block.Hash = sha256.Sum256(blockheight_string)
	return block
}

func txValid(tx block.Tx) bool {
	sufficientBalance := Accounts[tx.Sender] >= tx.Amount
	log.Println("sufficientBalance ", sufficientBalance, tx.Sender, Accounts[tx.Sender], tx.Amount)
	valid := sufficientBalance //and signature
	return valid
}

func moveCash(SenderAccount block.Account, ReceiverAccount block.Account, amount int) {
	log.Printf("move cash %v %v %v %v %d", SenderAccount, ReceiverAccount, Accounts[SenderAccount], Accounts[ReceiverAccount], amount)
	Accounts[SenderAccount] -= amount
	Accounts[ReceiverAccount] += amount

}

func applyTx(tx block.Tx) {
	//TODO check transaction type, not implemented yet

	moveCash(tx.Sender, tx.Receiver, tx.Amount)
}

func SetAccount(account block.Account, balance int) {
	Accounts[account] = balance
}

func ShowAccount(account block.Account) {
	log.Printf("%s %d", account, Accounts[account])
}

// pick a random account in the ledger
func RandomAccount() block.Account {
	lenk := len(Accounts)

	keys := make([]block.Account, 0, len(Accounts))
	for k := range Accounts {
		keys = append(keys, k)
	}

	rand.Seed(time.Now().UnixNano())
	ran := rand.Intn(lenk)

	randomAccount := keys[ran]
	log.Println("random account ", randomAccount)
	return randomAccount
}

func GenesisTx() block.Tx {
	Genesis_Account := block.AccountFromString(Genesis_Address)
	//block.AccountFromString("") //sender is empty

	//genesisSender := "" //genesisSender is the bootstrap account

	//log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	r_account := block.AccountFromString(address_r)
	genesisAmount := 20 //just a number for now
	//TODO id

	gTx := block.Tx{Nonce: randNonce, Sender: Genesis_Account, Receiver: r_account, Amount: genesisAmount}
	return gTx
}

func MakeGenesisBlock() block.Block {
	emptyhash := [32]byte{}
	timestamp := time.Now() //.Unix()
	b := []byte("banks on brink again")[:]
	genHash := sha256.Sum256(b)

	//add 10 genesis tx
	genesisTx := []block.Tx{}
	for i := 0; i < 10; i++ {
		someTx := GenesisTx()
		genesisTx = append(genesisTx, someTx)
	}

	genesis_block := block.Block{Height: 0, Txs: genesisTx, Prev_Block_Hash: emptyhash, Hash: genHash, Timestamp: timestamp}
	return genesis_block
}

//append block to chain of blocks
func AppendBlock(new_block block.Block) {
	Latest_block = new_block
	Blocks = append(Blocks, new_block)
}

//apply block to the state
func ApplyBlock(block block.Block) {
	//apply
	for j := 0; j < len(block.Txs); j++ {
		applyTx(block.Txs[j])
		//apply id
		block.Txs[j].Id = txHash(block.Txs[j])
	}
	//TODO add block hash now
}

// function to create blocks, called periodically
// currently assumes we can create blocks at will and we don't sync
func MakeBlock(t time.Time) {

	log.Printf("make block?")
	start := time.Now()
	//elapsed := time.Since(start)
	log.Printf("%s", start)

	//create new block if there is tx in the pool
	if len(Tx_pool) > 0 {

		timestamp := time.Now() //.Unix()
		new_block := block.Block{Height: len(Blocks), Txs: Tx_pool, Prev_Block_Hash: Latest_block.Hash, Timestamp: timestamp}
		new_block = blockHash(new_block)
		ApplyBlock(new_block)
		AppendBlock(new_block)

		log.Printf("new block %v", new_block)
		EmptyPool()

		Latest_block = new_block
	} else {
		log.Printf("no block to make")
		//handle special case of no tx
		//now we don't add blocks, which means there are empty periods and blocks are not evenly spaced in time
	}

}

/*func randomAccount() block.Account {
	/*rand.Seed(time.Now().UnixNano())
	keys := reflect.ValueOf(accounts).MapKeys()
	rkey := rand.Intn(len(keys))
	raccount := keys[rkey]
	return raccount
}*/
