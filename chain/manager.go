package chain

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/polygonledger/node/block"
	crypto "github.com/polygonledger/node/crypto"
)

var Tx_pool []block.Tx
var Blocks []block.Block
var Latest_block block.Block
var Accounts map[block.Account]int

const storageFile = "chain.json"

//TODO fix circular import
const (
	//Genesis_Address string = "P0614579c42f2"
	Genesis_Address string = "P2e2bfb58c9db"
	//Treasury_Address string = "PXXXXX"
)

//TODO
func GenesisKeys() crypto.Keypair {
	keypair := crypto.PairFromSecret("genesis")
	return keypair

}

//testing
func InitAccounts() {
	Accounts = make(map[block.Account]int)

	log.Println("init accounts %i", len(Accounts))
	//Genesis_Account := block.AccountFromString(Genesis_Address)
	//set genesiss account, this is the amount that the genesis address receives
	genesisAmount := 400
	SetAccount(block.AccountFromString(Genesis_Address), genesisAmount)
}

//valid cash transaction
//instead of needing to evluate bytecode like Bitcoin or Ethereum this is hardcoded cash transaction, no multisig, no timelocks
//* sufficient balance of sender (the sender has the cash, no credit as of now)
//* the sender is who he says he is (authorized access to funds)
//speed of evaluation should be way less than 1 msec
//TODO check nonce
func txValid(tx block.Tx) bool {

	//TODO check receiver has valid address format

	sufficientBalance := Accounts[tx.Sender] >= tx.Amount
	if !sufficientBalance {
		log.Println("insufficientBalance")
	} else {
		log.Println("suffcientBalance")
	}
	// log.Println("sufficientBalance ", sufficientBalance, tx.Sender, Accounts[tx.Sender], tx.Amount)
	//TODO and signature

	//the transaction is signed by the sender
	//TODO fix this is only for testing

	verified := crypto.VerifyTxSig(tx)
	btxValid := sufficientBalance && verified
	log.Println("sigvalid ", verified)
	//TODO check sig
	return btxValid
	//return true
}

//handlers

func HandleTx(tx block.Tx) string {
	//hash of timestamp is same, check lenght of bytes used??
	//timestamp := time.Now().Unix()

	//TODO check timestamp
	//log.Println("hash %x time %s sign %x", tx.Id, timestamp, tx.Signature)

	//TODO its own function
	if txValid(tx) {
		tx.Id = crypto.TxHash(tx)
		Tx_pool = append(Tx_pool, tx)
		return "ok"
	} else {
		log.Println("invalid tx")
		return "error"
	}

	//log.Printf("tx_pool_size: \n%d", tx_pool_size)

	//log.Printf("tx_pool size: %d\n", len(Tx_pool))
	//return "ok"

}

//#### blockchain functions

//empty the tx pool
func EmptyPool() {
	Tx_pool = []block.Tx{}
}

func blockHash(block block.Block) block.Block {
	//FIX hash of data, merkle tree
	timeFormat := "2020-02-02 16:06:06"
	new_hash := []byte(string(block.Timestamp.Format(timeFormat))[:])
	blockheight_string := []byte(strconv.FormatInt(int64(block.Height), 10))
	//new_hash = append(new_hash, blockheight_string)
	log.Printf("newhash %x", new_hash)
	block.Hash = sha256.Sum256(blockheight_string)
	return block
}

//move cash in the chain, we should know tx is checked to be valid by now
func moveCash(SenderAccount block.Account, ReceiverAccount block.Account, amount int) {
	log.Printf("move cash %v %v %v %v %d", SenderAccount, ReceiverAccount, Accounts[SenderAccount], Accounts[ReceiverAccount], amount)

	Accounts[SenderAccount] -= amount
	Accounts[ReceiverAccount] += amount
}

func applyTx(tx block.Tx) {
	//TODO check transaction type, not implemented yet
	valid := true //txValid(tx)
	if valid {
		moveCash(tx.Sender, tx.Receiver, tx.Amount)
	} else {
		log.Printf("tx invalid, dont apply")
		//handle error
	}
}

func SetAccount(account block.Account, balance int) {
	Accounts[account] = balance
}

func ShowAccount(account block.Account) {
	log.Printf("%s %d", account, Accounts[account])
}

func RandomAccount() block.Account {
	lenk := len(Accounts)

	keys := make([]block.Account, 0, len(Accounts))
	for k := range Accounts {
		keys = append(keys, k)
	}

	rand.Seed(time.Now().UnixNano())
	log.Println(lenk)
	ran := rand.Intn(lenk)

	randomAccount := keys[ran]
	log.Println("random account ", randomAccount)
	return randomAccount
}

func GenesisTx() block.Tx {
	Genesis_Account := block.AccountFromString(Genesis_Address)

	//log.Printf("%s", s)
	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	//TODO fix
	r := crypto.RandomPublicKey()
	//kp := crypto.PairFromSecret("basic")
	address_r := crypto.Address(r)
	r_account := block.AccountFromString(address_r)
	genesisAmount := 20 //just a number for now

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
		//if success
		//assign id
		//block.Txs[j].Id = txHash(block.Txs[j])
		//og.Println("hash ", block.Txs[j].Id)
	}
}

//trivial json storage
func writeChain() {
	dataJson, _ := json.Marshal(Blocks)
	ioutil.WriteFile(storageFile, []byte(dataJson), 0644)
}

//TODO error
func ReadChain() bool {

	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Println("storage file does not exist")
		return false
	}

	dat, _ := ioutil.ReadFile(storageFile)
	//var tx block.Tx

	if err := json.Unmarshal(dat, &Blocks); err != nil {
		panic(err)
	}

	log.Printf("read chain success. block height %d", len(Blocks))
	return true

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

		writeChain()

	} else {
		log.Printf("no block to make")
		//handle special case of no tx
		//now we don't add blocks, which means there are empty periods and blocks are not evenly spaced in time
	}

}
