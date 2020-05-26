package chain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/polygonledger/node/block"
	"github.com/polygonledger/node/crypto"
	"github.com/polygonledger/node/netio"
)

type AccountState struct {
	Accounts map[string]int
}
type ChainManager struct {
	Tx_pool []block.Tx
	Blocks  []block.Block
	State   AccountState
}

const ChainStorageFile = "data/chain.json"
const GenblockStorageFile = "data/genesis.json"

const (
	//Treasury_Address string = "P0614579c42f2"
	Treasury_Address string = "P2e2bfb58c9db"
	//Treasury_Address string = "PXXXXX"
)

func vlog(msg string) {
	//log.Println(string)
}

//TODO
func GenesisKeys() crypto.Keypair {
	keypair := crypto.PairFromSecret("genesis")
	return keypair

}

func CreateManager() ChainManager {
	//
	mgr := ChainManager{Tx_pool: []block.Tx{}, Blocks: []block.Block{}, State: AccountState{Accounts: make(map[string]int)}}
	return mgr
}

func (mgr *ChainManager) BlockHeight() int {
	return len(mgr.Blocks)
}

func (mgr *ChainManager) IsTreasury(account block.Account) bool {
	return account.AccountKey == Treasury_Address
}

//init genesis account
func (mgr *ChainManager) InitAccounts() {
	mgr.State.Accounts = make(map[string]int)

	vlog(fmt.Sprintf("init accounts %d", len(mgr.State.Accounts)))
	//Genesis_Account := block.AccountFromString(Treasury_Address)
	//set genesiss account, this is the amount that the genesis address receives
	genesisAmount := 400
	//tr := block.AccountFromString(Treasury_Address)
	mgr.SetAccount(Treasury_Address, genesisAmount)
	vlog(fmt.Sprintf("mgr.Accounts %v", mgr.State.Accounts))
}

//valid cash transaction
//instead of needing to evluate bytecode like Bitcoin or Ethereum this is hardcoded cash transaction, no multisig, no timelocks
//* sufficient balance of sender (the sender has the cash, no credit as of now)
//* the sender is who he says he is (authorized access to funds)
//speed of evaluation should be way less than 1 msec
//TODO check nonce
func TxValid(mgr *ChainManager, tx block.Tx) bool {

	//TODO check receiver has valid address format

	sufficientBalance := mgr.State.Accounts[tx.Sender] >= tx.Amount
	if !sufficientBalance {
		vlog("insufficientBalance")
	} else {
		vlog("suffcientBalance")
	}
	// vlog("sufficientBalance ", sufficientBalance, tx.Sender, Accounts[tx.Sender], tx.Amount)
	//TODO and signature

	//the transaction is signed by the sender
	//TODO fix this is only for testing

	verified := crypto.VerifyTxSig(tx)
	bTxValid := sufficientBalance && verified
	vlog(fmt.Sprintf("sigvalid %v", verified))
	//TODO check sig
	return bTxValid
}

func HandleTx(mgr *ChainManager, tx block.Tx) netio.Message {
	//hash of timestamp is same, check lenght of bytes used?
	//timestamp := time.Now().Unix()

	//TODO check timestamp
	//vlog("hash %x time %s sign %x", tx.Id, timestamp, tx.Signature)

	//TODO its own function
	if TxValid(mgr, tx) {
		//tx.Id = crypto.TxHash(tx)
		mgr.Tx_pool = append(mgr.Tx_pool, tx)
		vlog(fmt.Sprintf("append tx to pool %v", mgr.Tx_pool))
		//msg := Message{messageType: netio.REP, command: netio.CMD_TX}
		status := "ok"
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_TX, Data: []byte(status)}
		return msg
	} else {
		vlog("invalid tx")
		status := "error"
		msg := netio.Message{MessageType: netio.REP, Command: netio.CMD_TX, Data: []byte(status)}
		return msg
	}

	//log.Printf("tx_pool_size: \n%d", tx_pool_size)

	//log.Printf("tx_pool size: %d\n", len(Tx_pool))
	//return "ok"

}

//empty the tx pool
func EmptyPool(mgr *ChainManager) {
	mgr.Tx_pool = []block.Tx{}
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
func (mgr *ChainManager) moveCash(SenderAccount string, ReceiverAccount string, amount int) {
	//vlog(fmt.Srpintf("move cash %v %v %v %v %d", SenderAccount, ReceiverAccount, mgr.State.Accounts[SenderAccount], mgr.State.Accounts[ReceiverAccount], amount))

	mgr.State.Accounts[SenderAccount] -= amount
	mgr.State.Accounts[ReceiverAccount] += amount
}

func (mgr *ChainManager) applyTx(tx block.Tx) {
	//TODO check transaction type, not implemented yet
	valid := true //TxValid(tx)
	if valid {
		mgr.moveCash(tx.Sender, tx.Receiver, tx.Amount)
	} else {
		log.Printf("tx invalid, dont apply")
		//handle error
	}
}

func (mgr *ChainManager) SetAccount(account string, balance int) {
	mgr.State.Accounts[account] = balance
}

func ShowAccount(mgr *ChainManager, account string) {
	log.Printf("%s %d", account, mgr.State.Accounts[account])
}

func (mgr *ChainManager) RandomAccount() string {
	lenk := len(mgr.State.Accounts)
	vlog(fmt.Sprintf("lenk %d", lenk))

	//TODO
	keys := make([]string, 0, len(mgr.State.Accounts))
	for k := range mgr.State.Accounts {
		vlog(fmt.Sprintf("%v ", k))
		keys = append(keys, k)
	}

	rand.Seed(time.Now().UnixNano())
	ran := rand.Intn(lenk)
	vlog(fmt.Sprintf("%v %v", ran, keys))

	randomAccount := keys[ran]
	vlog(fmt.Sprintf("random account %v", randomAccount))
	return randomAccount
}

func GenesisTx() block.Tx {
	//Genesis_Account := block.AccountFromString(Treasury_Address)

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)
	//TODO fix
	r := crypto.RandomPublicKey()
	//kp := crypto.PairFromSecret("basic")
	address_r := crypto.Address(r)
	//r_account := block.AccountFromString(address_r)
	genesisAmount := 20 //just a number for now

	gTx := block.Tx{Nonce: randNonce, Sender: Treasury_Address, Receiver: address_r, Amount: genesisAmount}
	return gTx
}

func MakeGenesisBlock() block.Block {

	emptyhash := [32]byte{}
	timestamp := time.Now() //.Unix()
	b := []byte("banks on brink again")[:]
	genHash := sha256.Sum256(b)

	//add 10 genesis tx
	genesisTx := []block.Tx{}
	numseeder := 10
	for i := 0; i < numseeder; i++ {
		someTx := GenesisTx()
		genesisTx = append(genesisTx, someTx)
	}

	genesis_block := block.Block{Height: 0, Txs: genesisTx, Prev_Block_Hash: emptyhash, Hash: genHash, Timestamp: timestamp}
	return genesis_block
}

//append block to chain of blocks
func (mgr *ChainManager) AppendBlock(new_block block.Block) {
	mgr.Blocks = append(mgr.Blocks, new_block)
}

//reset blocks to 0
func (mgr *ChainManager) ResetBlocks() {
	mgr.Blocks = make([]block.Block, 0)
}

//apply block to the state
func (mgr *ChainManager) ApplyBlock(block block.Block) {
	vlog(fmt.Sprintf("ApplyBlock %v", block))
	//apply
	for j := 0; j < len(block.Txs); j++ {
		mgr.applyTx(block.Txs[j])
		//if success
		//assign id
		//block.Txs[j].Id = txHash(block.Txs[j])
		//og.Println("hash ", block.Txs[j].Id)
	}
}

func (mgr *ChainManager) ApplyBlocks(blocks []block.Block) {
	for _, block := range blocks {
		mgr.ApplyBlock(block)
	}
}

//trivial json storage
func (mgr *ChainManager) WriteChain() {
	dataJson, _ := json.Marshal(mgr.Blocks)
	ioutil.WriteFile(ChainStorageFile, []byte(dataJson), 0644)
}

//TODO error
func (mgr *ChainManager) ReadChain() bool {

	if _, err := os.Stat(ChainStorageFile); os.IsNotExist(err) {
		// path/to/whatever does not exist
		vlog("storage file does not exist")
		return false
	}

	dat, _ := ioutil.ReadFile(ChainStorageFile)

	var newblocks []block.Block
	if err := json.Unmarshal(dat, &newblocks); err != nil {
		panic(err)
	}

	//apply blocks
	mgr.Blocks = newblocks
	mgr.ApplyBlocks(newblocks)

	vlog(fmt.Sprintf("read chain success from %s. block height %d", ChainStorageFile, len(mgr.Blocks)))
	return true

}

func WriteState() {
	//TODO
}

func ReadState() {
	//TODO
}

func WriteGenBlock(block block.Block) {
	//TODO
	dataJson, _ := json.Marshal(block)
	//dataJson, _ := json.MarshalIndent(block, "", "    ")
	ioutil.WriteFile(GenblockStorageFile, []byte(dataJson), 0644)
}

func ReadGenBlock() block.Block {
	//TODO

	if _, err := os.Stat(GenblockStorageFile); os.IsNotExist(err) {
		// path/to/whatever does not exist
		vlog("storage file does not exist")
		//return nil
	}

	dat, _ := ioutil.ReadFile(GenblockStorageFile)

	var newgenblock block.Block

	if err := json.Unmarshal(dat, &newgenblock); err != nil {
		panic(err)
	}

	log.Printf("read gen block success from %s", GenblockStorageFile)

	return newgenblock
}

func (mgr *ChainManager) LastBlock() block.Block {
	n := len(mgr.Blocks)
	return mgr.Blocks[n-1]
}

// function to create blocks, called periodically
// currently assumes we can create blocks at will and we don't sync
func MakeBlock(mgr *ChainManager) {

	//log.Printf("make block? ")
	//start := time.Now()
	//elapsed := time.Since(start)
	//log.Printf("%s", start)

	//create new block if there is tx in the pool
	if len(mgr.Tx_pool) > 0 {

		timestamp := time.Now() //.Unix()
		Latest_block := mgr.LastBlock()
		new_block := block.Block{Height: len(mgr.Blocks), Txs: mgr.Tx_pool, Prev_Block_Hash: Latest_block.Hash, Timestamp: timestamp}
		new_block = blockHash(new_block)

		mgr.ApplyBlock(new_block)
		//TODO! fix
		mgr.AppendBlock(new_block)

		vlog(fmt.Sprintf("new block %v", new_block))
		EmptyPool(mgr)

		mgr.WriteChain()

	} else {
		//log.Printf("no block to make")
		//handle special case of no tx
		//now we don't add blocks, which means there are empty periods and blocks are not evenly spaced in time
	}

}

func MakeBlockLoop(mgr *ChainManager, blocktime time.Duration) {

	go func() {
		for {
			MakeBlock(mgr)
			time.Sleep(blocktime)
		}
	}()
}
