package core

type BlockChain struct {
	Blocks []*Block
}

//AddBlock saves provided data as a block in the blockchain
func (bc *BlockChain) AddBlock(data string) {
	if len(bc.Blocks) == 0 {
		bc.Blocks = append(bc.Blocks,NewGenesisBlock())
		return
	}
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data,prevBlock.Hash)
	bc.Blocks = append(bc.Blocks,newBlock)
}

//NewBlockChain creates a new BlockChain with genesis Block
func NewBlockChain() *BlockChain  {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}