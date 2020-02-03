package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//1.定义结构
type Block struct {
	//版本号
	Version uint64
	//前区块哈希
	PrevHash []byte
	//Merkel根（梅克尔根，一个哈希值）
	MerkelRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Difficulty uint64
	//随机数，挖矿要找的数据
	Nonce uint64
	//a.当前区块哈希
	Hash []byte
	//b.数据
	//Data []byte
	//真实的交易数组
	Transactions []*Transaction

}


//实现一个辅助函数，将uint64转成byte
func UintToByte(num uint64) []byte {
	//todo
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
//序列化
func (block *Block) Serialize() []byte  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err!=nil{
		log.Panic("编码出错")
	}

	return buffer.Bytes()
}
//反序列化
func Deserialize(data []byte) Block  {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err := decoder.Decode(&block)
	if err!=nil {
		log.Panic("解码出错")
	}

	return block
}




//2.创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{}, //先填空，后面计算
		//Data:       []byte(data),
		Transactions:txs,
	}
	block.MerkelRoot=block.MakeMerkelRoot()
	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//查找随机数，不停的进行哈希运算
	hash, nonce := pow.Run()
	block.Hash=hash
	block.Nonce=nonce
	return &block
}

//3.生成哈希 block内部函数
func (block *Block) SetHash() {
	//var blockInfo []byte
	//1.拼装数据
	//blockInfo = append(block.PrevHash, UintToByte(block.Version)...)
	//blockInfo = append(block.PrevHash, block.PrevHash...)
	//blockInfo = append(block.PrevHash, block.MerkelRoot...)
	//blockInfo = append(block.PrevHash, UintToByte(block.TimeStamp)...)
	//blockInfo = append(block.PrevHash, UintToByte(block.Difficulty)...)
	//blockInfo = append(block.PrevHash, UintToByte(block.Nonce)...)
	//blockInfo = append(block.PrevHash, block.Data...)
	tmp:=[][]byte{
		UintToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		UintToByte(block.TimeStamp),
		UintToByte(block.Difficulty),
		UintToByte(block.Nonce),
		//block.Transactions,
	}
	//将二维切片数组连接起来，返回一个以为切片
	blockInfo := bytes.Join(tmp, []byte{})

	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
//模拟梅克尔根，只对交易的数据做简单拼接，不做二叉树处理
func (block *Block) MakeMerkelRoot() []byte {
	return []byte{}
}
