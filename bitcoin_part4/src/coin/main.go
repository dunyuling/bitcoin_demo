package main

import (
	"core"
)

func main() {
	cli := core.CLI{}
	//cli.Run()

	//cli.CreateBlockChainTest("Ivan")
	//cli.GetBalanceTest("Ivan")
	//cli.SendTest("Ivan","Pedro",6)
	cli.GetBalanceTest("Ivan")

}