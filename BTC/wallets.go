package main

import (
	"io/ioutil"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/elliptic"
	"os"
	"./lib/base58"
)

const walletFile = "wallet.dat"

//定一个 Wallets结构，它保存所有的wallet以及它的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

//创建方法，返回当前所有钱包的实例
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet

	ws.saveToFile()
	return address
}

//保存方法，把新建的wallet添加进去
func (ws *Wallets) saveToFile() {

	var buffer bytes.Buffer

	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	//一定要注意校验！！！
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
}

//读取文件方法，把所有的wallet读出来
func (ws *Wallets) loadFile() {
	//在读取之前，要先确认文件是否在，如果不存在，直接退出
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		//ws.WalletsMap = make(map[string]*Wallet)
		return
	}

	//读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	//解码
	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(content))

	var wsLocal Wallets

	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}

	//ws = &wsLocal
	//对于结构来说，里面有map的，要指定赋值，不要再最外层直接赋值
	ws.WalletsMap = wsLocal.WalletsMap
}

func (ws *Wallets) ListAllAddresses() []string {
	var addresses []string
	//遍历钱包，将所有的key取出来返回
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}

//通过地址返回公钥的哈希值
func GetPubKeyFromAddress(address string) []byte {
	//1. 解码
	//2. 截取出公钥哈希：去除version（1字节），去除校验码（4字节）
	addressByte := base58.Decode(address) //25字节
	len := len(addressByte)

	pubKeyHash := addressByte[1:len-4]

	return pubKeyHash
}
