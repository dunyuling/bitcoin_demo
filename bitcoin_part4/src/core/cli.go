package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

//CLI responsible for processing command line arguments
type CLI struct {}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tgetbalance -address ADDRESS  ---Get balance of ADDRESS")
	fmt.Println("\tcreateblockchain -address ADDRESS ---create a blockchain and send genesis block reward to ADDRESS ")
	fmt.Println("\tprintchain ---print all the blocks of the blockchain")
	fmt.Println("\tsend -from From -to To -amount AMOUNT ---Send AMOUNT of coins from FROM ADDRESS to To ADDRESS")
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain()
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _,out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s':%d\n",address,balance)
}

func (cli *CLI) createBlockChain(address string) {
	bc := CreateBlockChain(address)
	defer bc.db.Close()
	fmt.Println("create blockchain " + address + " done!")
}

func (cli *CLI) send(from,to string,amount int) {
	bc := NewBlockChain()
	defer bc.db.Close()

	tx := NewUTXOTransaction(from,to,amount,bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("send from " + from + " to " + to + " success!")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	bc := NewBlockChain()
	defer bc.db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

//Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd :=flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd :=flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	sendCmd :=flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address","","The address to get balance ")
	createBlockChainAddress := createBlockChainCmd.String("address","","The address to")
	sendFrom := sendCmd.String("from","","Source wallet address")
	sendTo := sendCmd.String("to","","Destination wallet address")
	sendAmount := sendCmd.Int("amount",0,"Amount to Send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockChainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom,*sendTo,*sendAmount)
	}

}

func (cli *CLI) GetBalanceTest(address string) {
	cli.getBalance(address)
}

func (cli *CLI) CreateBlockChainTest(address string)  {
	cli.createBlockChain(address)
}

func (cli *CLI) SendTest(from,to string, amount int)  {
	cli.send(from,to,amount)

}

func (cli *CLI) PrintChainTest()  {
	cli.printChain()
}