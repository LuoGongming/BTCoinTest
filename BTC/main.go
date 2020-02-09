package main

func main()  {
	bc:=NewBlockChain("邝")
	cli:=CLI{bc}
	cli.Run()


	/*bc.AddBlock("罗公明向邝靖月转了10个比特币")
	bc.AddBlock("罗公明又向邝靖月转了10个比特币")
	//创建迭代器
	it:=bc.NewIterator()
	for{
		//返回区块，左移
		block := it.Next()
		fmt.Printf("\n\n")
		fmt.Printf("前区块哈希值：%x\n",block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n",block.Hash)
		fmt.Printf("区块数据：%s\n",block.Data)
		if len(block.PrevHash)==0{
			break
		}

	}*/

	/*for i,block:=range bc.blocks{


	}*/
}
