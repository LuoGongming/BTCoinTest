package main

import (
	"fmt"
	"os"
	"strconv"
)

//这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA	"add data to blockchain"
	printChain				"print all blockchain data"
	printChainR				"反向打印区块链"
	getBalance --address ADDRESS "获取指定地址余额"
	send FROM TO AMOUNT MINER DATA "由from转amount给to,由miner挖矿，写入data"
`

//接收参数的动作放到函数中
func (cli *CLI) Run() {
	//1.得到所有命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Printf("添加区块")
		//确保命令有效
		if len(args) == 4 && args[2] == "--data" {
			//data:=args[3]
			//cli.bc.AddBlock(data)
		} else {
			fmt.Printf("添加数据有误请检查\n")
			fmt.Printf(Usage)
		}

		//a.获取数据
		//b.使用bc添加区块链AddBlock

	case "printChain":
		fmt.Printf("打印区块\n")
		cli.PrintBlockChain()
	case "printChainR":
		fmt.Printf("反向打印区块\n")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Printf("获取余额\n")
		if len(args) == 4 && args[2] == "--address" {
			address:=args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始.....\n")
	if len(args)!=7{
		fmt.Printf("参数个数错误，请检查\n")
		fmt.Printf(Usage)
	}
	from:=args[2]
	to:=args[3]
	amount,_:=strconv.ParseFloat(os.Args[4],64)
	miner:=args[5]
	data:=args[6]
	cli.Send(from,to,amount,miner,data)
	default:
		fmt.Printf("无效的命令，请检查")
		fmt.Printf(Usage)

	}
}



