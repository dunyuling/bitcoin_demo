package core

import (
	"math"
	"math/big"
	"bytes"
	"fmt"
	"crypto/sha256"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 36

//ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block *Block
	target *big.Int
}

//NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte  {
	data := bytes.Join(
			[][]byte{
				pow.block.PrevBlockHash,
				pow.block.Data,
				IntToHex(pow.block.Timestamp),
				IntToHex(int64(targetBits)),
				IntToHex(int64(nonce))},
			[]byte{})
	return data
}

//Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int,[]byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n",pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		//fmt.Printf("\r%v\n,%x\n,%v\n,%v\n,%v\n,%v\n,%v\n\n\n",data,hash,&hashInt, hashInt,pow.target,*pow.target,nonce)
		fmt.Printf("\r%x",hash)

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")

	return nonce,hash[:]
}

//Validate validate block's PoW
func(pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}