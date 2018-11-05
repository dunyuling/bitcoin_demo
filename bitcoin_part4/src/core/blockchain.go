package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"encoding/hex"
)

const dbFie = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The time 03/Jan/2009 Chancellor on brink of second bailout for banks"

//BlockChain keeps a sequence of  Blocks
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

//BlockChainIterator is used to iterate over blockchain blocks
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

//MineBlock mines a new block with the provided transactions
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions,lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"),newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//Iterator
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}
	return bci
}

//Next
func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodeBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodeBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}

func dbExists() bool {
	if _,err := os.Stat(dbFie);os.IsNotExist(err) {
		return false
	}
	return true
}

//NewBlockChain creates a new BlockChain with genesis Block
func NewBlockChain() *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFie, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}
	return &bc
}

//CreateBlockChain creates a new blockchain DB
func CreateBlockChain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("BlockChain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFie, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		//b := tx.Bucket([]byte(blocksBucket))
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}
	return &bc
}

//FindUnSpentTransactions returns a list of transactions containing unspent outputs
func (bc *BlockChain) FindUnSpentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _,tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx,out := range tx.Vout {
				//Was the output spent
				if spentTXOs[txID] != nil {
					for _,spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs,*tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _,in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

//FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int,map[string][]int) {
	unspentOutputs :=  make(map[string][]int)
	unspentTXs := bc.FindUnSpentTransactions(address)
	accumulated := 0

Work:
	for _,tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx,out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated,unspentOutputs
}

//FindUTXO finds and returns all unspent transaction outputs
func (bc *BlockChain) FindUTXO(address string) []TXOutput{
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnSpentTransactions(address)

	for _,tx := range unspentTransactions {
		for _,out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs,out)
			}
		}
	}
	return UTXOs
}