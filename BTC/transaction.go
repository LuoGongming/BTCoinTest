package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//1.定义交易结构
type Transaction struct {
	TXID []byte //交易id
	TXInputs []TXInput	//交易输入的数组
	TXOutputs []TXOutput	//交易输出的数组
}
//定义教育输入
type TXInput struct {
	//引用交易ID
	TXid []byte
	//引用output索引值
	Index int64
	//解锁脚本，用地址来模拟
	Sig string
}
//定义交易输出
type TXOutput struct {
	//转账金额
	value float64
	//所定脚本，用地址来模拟
	PubKeyHash string
}
//设置交易id
func (tx *Transaction) SetHash()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	data:=buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID=hash[:]
}

const reward  =12.5
//2.创建交易方法(挖矿交易)
func NewCoinbaseTX(address string,data string) *Transaction  {
	//挖矿交易特点
	//1.只有一个input
	//2.无需引用交易id
	//3.无需引用index
	//旷工由于挖矿时无需指定签名，所以sig字段可以有矿工自由填写数据，一般是矿池名
	input:=TXInput{[]byte{},-1,data}
	output:=TXOutput{reward,address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
}
//3.创建挖矿交易
//4.根据交易调整程序
