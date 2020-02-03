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
	db *bolt.DB
	tail []byte //存储最后一个区块哈希

}

const blockChainDB  = "blockChain.db"
const blockBucket  = "blockBucket"
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
	if err!=nil {
		log.Panic("打开数据库失败")
	}
	//操作数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket(如果没有就创建)
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil {
			bucket,err=tx.CreateBucket([]byte(blockBucket))
			if err!=nil {
				log.Panic("创建bucket(blockBucket)失败")
			}
			//将创世块作为第一个区块添加到区块链中
			genesisBlock:=GenesisBlock(address)
			//写数据
			//hash作为key,block的字节流作为value
			bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"),genesisBlock.Hash)
			lastHash=genesisBlock.Hash

			///**
			//	读数据测试用************
			//*/
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block:=Deserialize(blockBytes)
			//fmt.Printf("block Info:%s\n",block)

		}else {
			lastHash=bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db,lastHash}

}
//定义一个创世块
func GenesisBlock(address string) *Block{
	coinbase :=NewCoinbaseTX(address,"第一个创世块")
	return NewBlock([]*Transaction{coinbase},[]byte{})
}


//6.添加区块
/*func (bc *BlockChain)AddBlock(data string,prevHash []byte)  {
	//创建新的区块
	block := NewBlock(data, prevHash)
	//添加到区块链数组中
	bc.blocks=append(bc.blocks,block)

}*/
func (bc *BlockChain)AddBlock(txs []*Transaction)  {
	/*//获取前区块哈希
	//获取最后一个区块
	lastBLock:=bc.blocks[len(bc.blocks)-1]
	prevHash:=lastBLock.Hash
	//创建新的区块
	block := NewBlock(data, prevHash)
	//添加到区块链数组中
	bc.blocks=append(bc.blocks,block)*/
	//获取前区块哈希
	db:=bc.db	//区块链数据库
	lastHash:=bc.tail	
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			log.Panic("bucket不应该为空，请检查")
		}
		block:=NewBlock(txs,lastHash)
		//hash 作为key block的字节作为value
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"),block.Hash)
		lastHash=block.Hash
		//更新一下内存中的区块链，把最后的小尾巴tail更新
		bc.tail=block.Hash
		return nil
	})
		//获取最后一个区块
		//创建新的区块
	//添加到区块链数组中
}
func (bc *BlockChain) PrintChain()  {
	blockHeight:=0
	bc.db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte("blockBucket"))
		//从第一个key->value 进行遍历，遍历到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k,[]byte("LastHashKey")){
				return nil
			}
			block:=Deserialize(v)
			fmt.Printf("====================区块高度：%d==========\n",blockHeight)
			blockHeight++
			fmt.Printf("版本号：%d\n",block.Version)
			fmt.Printf("前区块哈希值：%x\n",block.PrevHash)
			fmt.Printf("梅克尔根：%x\n",block.MerkelRoot)
			fmt.Printf("时间戳：%d\n",block.TimeStamp)
			fmt.Printf("难度值（随便写的）：%d\n",block.Difficulty)
			fmt.Printf("随机数：%d\n",block.Nonce)
			fmt.Printf("当前区块哈希值：%x\n",block.Hash)
			fmt.Printf("区块数据：%s\n",block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}
//找到指定地址的所有utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	return UTXO
}


