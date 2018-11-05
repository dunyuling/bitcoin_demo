package core

import (
	"log"
	"fmt"
)

func (cli *CLI) createBlockChain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockChain(address)
	defer bc.db.Close()
	fmt.Println("Done!")
}
