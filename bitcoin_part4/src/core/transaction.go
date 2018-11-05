package core

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const subsidy = 10

//Transaction represents a Bitcoin transaction
type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

//TXInput represents a transaction input
type TXInput struct {
	Txid []byte
	Vout int
	ScriptSig string
}

//TXOutput represents a transaction output
type TXOutput struct {
	Value int
	ScriptPubKey string
}

//IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) SetID() {
	var encodedBuf bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encodedBuf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encodedBuf.Bytes())
	tx.ID = hash[:]
}

//CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

//CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

//NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to ,data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'",to)
	}

	txin := TXInput{[]byte{},-1,data}
	txout := TXOutput{subsidy,to}
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}
	tx.SetID()

	return &tx
}

//NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction  {
	var inputs []TXInput
	var outputs []TXOutput

	acc,validOutputs := bc.FindSpendableOutputs(from,amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	//Builds a list of inputs
	for txid,outs := range validOutputs {
		txID,err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _,out := range outs {
			input := TXInput{txID,out,from}
			inputs = append(inputs,input)
		}
	}

	//Builds a list of outputs
	outputs = append(outputs,TXOutput{amount,to})
	if acc > amount {
		outputs = append(outputs,TXOutput{acc-amount,from})
	}

	tx := Transaction{nil,inputs,outputs}
	tx.SetID()

	return &tx
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       ScriptSig: %s", input.ScriptSig))
	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %s", output.ScriptPubKey))
	}
	return strings.Join(lines, "\n")
}

//Serialize serializes the block
func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(*tx)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}