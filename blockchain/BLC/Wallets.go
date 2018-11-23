package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
)

// 钱包集合，对钱包进行维护、管理的文件

// 钱包集合的文件,存储钱包集合
const walletFile = "Wallets_%s.dat"

// 钱包集合的结构
type Wallets struct {
	// key : 地址
	// value : 钱包结构
	Wallets map[string] *Wallet
}

// 初始化钱包集合
func NewWallets(nodeID string) *Wallets {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	// 1. 判断文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string] *Wallet)
		return wallets
	}
	// 2. 文件存在，读取内容
	fileContent, err := ioutil.ReadFile(walletFile)
	if nil != err {
		log.Panicf("get the file content failed! %v\n", err)
	}

	var wallets Wallets
	// register主要用于需要解析的参数中包含interface的情况
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if nil != err {
		log.Panicf("decode filecontent failed! %v\n", err)
	}
	return &wallets
}

// 创建新的钱包,并且添加到钱包集合中
func (wallets *Wallets) CreateWallet(nodeID string) {
	// 新建钱包对象
	wallet := NewWallet()
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	// 保存到文件中
	wallets.SaveWallets(nodeID)
}

// 持久化钱包信息(写入文件)
func (w *Wallets) SaveWallets(nodeID string) {
	var content bytes.Buffer
	// 注册
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	// 序列化钱包数据
	err := encoder.Encode(&w)
	if nil != err {
		log.Panicf("encode the struct of wallets failed! %v\n", err)
	}
	walletFile := fmt.Sprintf(walletFile, nodeID)
	// 保存到文件中
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if nil != err {
		log.Panicf("write the content of wallets to file failed! %v\n", err)
	}
}