package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义工作量证明结构proofofwork
type ProofOfWork struct {
	//a.block
	block *Block
	//b.目标值,一个非常大的数
	target *big.Int
}

//提供创建pow函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//指定的难度值，目前string类型，需转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，目的是将上面难度值转化成big.Int
	temInt := big.Int{}
	//将难度值赋值给big.int 16进制
	temInt.SetString(targetStr, 16)
	pow.target = &temInt
	return &pow
}

//提供不断计算的hash函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	block := pow.block
	var hash [32]byte
	fmt.Println("开始挖矿。。。")
	for {
		//1.拼装数据（区块数据，不断变化的随机数）
		tmp := [][]byte{
			UintToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			UintToByte(block.TimeStamp),
			UintToByte(block.Difficulty),
			UintToByte(nonce),
			//只对区块头做哈希值，区块体通过MerkelRoot影响
			//block.Data,
		}
		//将二维切片数组连接起来，返回一个以为切片
		blockInfo := bytes.Join(tmp, []byte{})
		//2.做哈希运算
		hash = sha256.Sum256(blockInfo)
		//3.与pow中的target进行比较
		tmpInt := big.Int{}
		//讲的到的hash数组转化成一个big.Int
		tmpInt.SetBytes(hash[:])
		//比较当前哈希与目标哈希值，如果当前哈希值小于目标哈希值说明找到了
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//
		//func (x *Int) Cmp(y *Int) (r int) {

		if tmpInt.Cmp(pow.target) == -1 {
			//a.找到了，退出返回
			fmt.Printf("找到了，挖矿成功，hash：%x,nonce :%d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			//b.没找到，继续找随机数+1
			nonce++
			//fmt.Println("还没找到 000:",nonce)
		}
	}

}
