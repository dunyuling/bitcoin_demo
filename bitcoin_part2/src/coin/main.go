package main

import (
	"core"
	"fmt"
	"strconv"
	"time"
)

func main() {
	/*var a *big.Int
	var b *big.Int
	a = big.NewInt(1)
	b = big.NewInt(1)
	if a.Cmp(b) == 1{
		fmt.Println("1111")
		fmt.Printf("%x\n",a)
		fmt.Printf("%x\n",b)
		fmt.Println()
	} else {
		fmt.Println("2222")
		fmt.Printf("%x\n",a)
		fmt.Printf("%x\n",b)
		fmt.Println()
	}*/

	/*target := big.NewInt(1)
	target.Lsh(target, uint(256 - 20))

	fmt.Println(target)*/

	fmt.Println(time.Now())
	bc := core.NewBlockChain()
	defer fmt.Println(time.Now())

	//bc.AddBlock("Send 1 BTC to Ivan")
	//bc.AddBlock("Send 2 more BTC to lhg")

	for _,block := range bc.Blocks {
		fmt.Printf("Prev.hash:%x\n",block.PrevBlockHash)
		fmt.Printf("Data:%s\n",block.Data)
		fmt.Printf("Hash:%x\n",block.Hash)

		pow := core.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n",strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
