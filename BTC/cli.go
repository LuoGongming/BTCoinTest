package main

import (
	"fmt"
	"os"
)

//这是一个用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA	"add data to blockchain"
	printChain				"print all blockchain data"
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
			data:=args[3]
			cli.bc.AddBlock(data)
		}else {
			fmt.Printf("添加数据有误请检查\n")
			fmt.Printf(Usage)
		}

		//a.获取数据
		//b.使用bc添加区块链AddBlock

	case "printChain":
		fmt.Printf("打印区块")
		cli.PrintBlockChain()
	default:
		fmt.Printf("无效的命令，请检查")
		fmt.Printf(Usage)

	}
}


