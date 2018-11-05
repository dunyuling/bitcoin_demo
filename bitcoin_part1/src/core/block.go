package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

//Block keeps block headers
type Block struct {
	Timestamp     int64  //区块创建时间戳
	Data          []byte //区块包含的数据
	PrevBlockHash []byte //上一个区块的Hash
	Hash          []byte //区块自身的Hash,用于验证区块数据的有效
}

//NewBlock creates and return Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Data: []byte(data), PrevBlockHash: prevBlockHash, Hash: []byte{}}
	block.SetHash()
	return block
}

//SetHash calculates and sets block hash
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewGenesisBlock() *Block  {
	return NewBlock("Genesis Block",[]byte{})
}
/*
func main() {
	prevBlockHash := sha256.Sum256([]byte("prevBlock"))
	fmt.Println("prevBlockHash:")
	fmt.Println(prevBlockHash[:])
	block := Block{Timestamp:time.Now().Unix(),Data:[]byte("123"),PrevBlockHash:prevBlockHash[:]}
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	fmt.Println("timestamp:")
	fmt.Println(timestamp)
	fmt.Println("data:")
	fmt.Println([]byte(block.Data))
	headers := bytes.Join([][]byte{block.PrevBlockHash,block.Data,timestamp},[]byte{})
	fmt.Println("headers:")
	fmt.Println(headers)
}*/