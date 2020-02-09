package main

import (
	bolt "bolt-master"
	"bytes"
	"fmt"
	"log"
)

//4.引入区块链
//blockChain代码重构，数据库代替
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	db   *bolt.DB
	tail []byte //存储最后一个区块哈希

}

const blockChainDB = "blockChain.db"
const blockBucket = "blockBucket"

//5.定义一个区块链
func NewBlockChain(address string) *BlockChain {

	/*return &BlockChain{
		blocks:[]*Block{genesisBlock},
	}*/
	//最后一个区块的哈希
	var lastHash []byte

	//打开数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)
	//defer db.Close()
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//操作数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket(如果没有就创建)
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(blockBucket)失败")
			}
			//将创世块作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address)
			//写数据
			//hash作为key,block的字节流作为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash

			///**
			//	读数据测试用************
			//*/
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block:=Deserialize(blockBytes)
			//fmt.Printf("block Info:%s\n",block)

		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}

}

//定义一个创世块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "第一个创世块")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//6.添加区块
/*func (bc *BlockChain)AddBlock(data string,prevHash []byte)  {
	//创建新的区块
	block := NewBlock(data, prevHash)
	//添加到区块链数组中
	bc.blocks=append(bc.blocks,block)

}*/
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	/*//获取前区块哈希
	//获取最后一个区块
	lastBLock:=bc.blocks[len(bc.blocks)-1]
	prevHash:=lastBLock.Hash
	//创建新的区块
	block := NewBlock(data, prevHash)
	//添加到区块链数组中
	bc.blocks=append(bc.blocks,block)*/
	//获取前区块哈希
	db := bc.db //区块链数据库
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应该为空，请检查")
		}
		block := NewBlock(txs, lastHash)
		//hash 作为key block的字节作为value
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)
		lastHash = block.Hash
		//更新一下内存中的区块链，把最后的小尾巴tail更新
		bc.tail = block.Hash
		return nil
	})
	//获取最后一个区块
	//创建新的区块
	//添加到区块链数组中
}
func (bc *BlockChain) PrintChain() {
	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))
		//从第一个key->value 进行遍历，遍历到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("====================区块高度：%d==========\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号：%d\n", block.Version)
			fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
			fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
			fmt.Printf("时间戳：%d\n", block.TimeStamp)
			fmt.Printf("难度值（随便写的）：%d\n", block.Difficulty)
			fmt.Printf("随机数：%d\n", block.Nonce)
			fmt.Printf("当前区块哈希值：%x\n", block.Hash)
			fmt.Printf("区块数据：%s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}

//找到指定地址的所有utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	txs:=bc.FindTUXOTransactions(address)
	for _,tx:=range txs{
		for _,output:=range tx.TXOutputs {
			if address==output.PubKeyHash{
				UTXO=append(UTXO,output)
			}

		}
	}
	return UTXO

}
//根据需求找到合理的utxo
func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64,float64) {
	//找到的合理的utxo集合
	 utxos:= make(map[string][]uint64)
	//找到的utxo里面包含的总数
	var calc float64

	txs:=bc.FindTUXOTransactions(from)
	for _,tx:=range txs{
		for i,output:=range tx.TXOutputs {
			//if from==output.PubKeyHash{
			//	UTXO=append(UTXO,output)
			//}
			if output.PubKeyHash == from {
				//UTXO = append(UTXO, output)

				//需要实现的逻辑，找到最少utxo
				//比较一下是否满足转账需求，满足的话知己返回，不满足继续统计
				if calc<amount{
					//把utxo加进来
					utxos[string(tx.TXID)]=append(utxos[string(tx.TXID)],uint64(i))
					//2统计一下当前utxo总额
					calc+=output.Value
					//加完之后满足条件了
					if calc>=amount{
						return utxos,calc
					}

				}else {
					fmt.Printf("不满足转账金额，当前总额%f,目标金额：%f\n",calc,amount)
				}
			}
		}
	}
	return utxos,calc
}

func (bc *BlockChain) FindTUXOTransactions(address string) []*Transaction {

	var txs []*Transaction //存储所有包含UTXO交易集合

	//定义一个map来保存消费过的output,key是这个output的交易id
	//map[交易id][]int64
	spendOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()
	for {
		//1.遍历区块
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//fmt.Printf("current txid %x\n", tx.TXID)
			//3.遍历output，找到和自己相关的utxo
		OUTPUT:
			for i, output := range tx.TXOutputs {
				fmt.Printf("current index %x\n", i)
				//在这里做一个过滤，将所有消耗过的outputs个当前所即将添加的output对比
				//如果相同就跳过，否则添加，
				//如果当前的交易id 存在与我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spendOutputs[string(tx.TXID)] != nil {
					for _, j := range spendOutputs[string(tx.TXID)] {
						//
						if int64(i) == j {
							//当前准备添加的output已经消耗过了，不要再加了
							continue OUTPUT
						}
					}
				}
				//这个output和我们目标地址相同，满足条件，加到返回utxo数组中
				if output.PubKeyHash == address {
					//UTXO = append(UTXO, output)
					//返回所有包含我的utxo交易集合
					txs=append(txs, tx)

				}

			}

			//如果当前交易是挖矿交易，那么不遍历直接跳过
			if !tx.IsCoinbase() {

				//4.遍历input找到自己花费过utxo的集合（把自己消耗过的表示出来）
				for _, input := range tx.TXInputs {
					//判断一下当前这个input和目标（李四）是否一致
					//如果一致说明是李四消耗过的output就加进来，
					if input.Sig == address {
						indexArray := spendOutputs[string(input.TXid)]
						indexArray = append(indexArray, input.Index)
					}
				}
			}}

		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块链遍历完成退出")
		}
	}

	return txs

}