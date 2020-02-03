package main

import (
	bolt "bolt-master"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	//游标，用于不断索引
	currentHashPainter []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator  {
	return &BlockChainIterator{
		db:                 bc.db,
		//最初指向最后一个区块，随next变化不断调用
		currentHashPainter: bc.tail,
	}
}
/*
迭代器属于区块链
Next方法属于迭代器
1.返回当前的区块
2.指针前移
*/
func (it *BlockChainIterator) Next() *Block  {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			log.Panic("迭代器遍历是bucket不应该为空")
		}
		blockTmp := bucket.Get(it.currentHashPainter)
		//解码
		block=Deserialize(blockTmp)
		//游标哈希左移
		it.currentHashPainter=block.PrevHash
	return nil
	})
	return &block
}