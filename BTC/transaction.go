package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

//1.定义交易结构
type Transaction struct {
	TXID      []byte     //交易id
	TXInputs  []TXInput  //交易输入的数组
	TXOutputs []TXOutput //交易输出的数组
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
	Value float64
	//所定脚本，用地址来模拟
	PubKeyHash string
}

//设置交易id
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//实现一个函数，判断当前交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//1.交易input只有一个
	//2.交易id 为空
	//3.交易index为-1。
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid)==0 && tx.TXInputs[0].Index==1{
		return true
	}
	return false
}

const reward = 12.5

//2.创建交易方法(挖矿交易)
func NewCoinbaseTX(address string, data string) *Transaction {
	//挖矿交易特点
	//1.只有一个input
	//2.无需引用交易id
	//3.无需引用index
	//旷工由于挖矿时无需指定签名，所以sig字段可以有矿工自由填写数据，一般是矿池名
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
}

// 创建普通转账交易
//1.找到最合理的UTXO集合，map[string][]uint64
//2.将这些utxo逐一转成inputs
//3.创建outputs
//4.如果有零钱要找零。
func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction  {
	//1.找到最合理的UTXO集合，map[string][]uint64
	utxos,resValue:=bc.FindNeedUTXOs(from,amount)
	if resValue<amount {
		fmt.Printf("余额不足，交易失败\n")
		return nil
	}


	var inputs []TXInput
	var outputs []TXOutput

	//2.将这些utxo逐一转成inputs
	for id,indexArray:=range utxos{
		for _,i:=range indexArray {
			  input:=TXInput{[]byte(id),int64(i),from}
			  inputs = append(inputs, input)
		}
	}
	//创建交易输出
	output:=TXOutput{amount,to}
	outputs = append(outputs, output)
	if resValue>amount {
		//找零
		outputs=append(outputs,TXOutput{resValue-amount,from})
	}
	tx:=Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	return &tx
}
