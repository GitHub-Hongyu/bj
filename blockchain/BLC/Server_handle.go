package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/labstack/gommon/log"
)

// 处理请求

// version
func handleVersion(request []byte, bc *BlockChain)  {
	fmt.Println("handleVersion")
	var buff bytes.Buffer
	var data Version
	// 获取request数据
	dataBytes := request[12:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the version cmd failed! %v\n", err)
	}
	// 获取区块高度
	height := bc.GetHeight()
	versionHeight := data.Height
	// 如果当前节点的区块高度大于versionHeight,将当前节点版本信息发送给请求节点
	if height > int64(versionHeight) {
		sendVersion(data.AddrFrom, bc)
	} else if height < int64(versionHeight){
		// 如果当前节点的区块高度小于versionHeight，向发送请求的节点发送同步请求
		sendGetBlocks(data.AddrFrom)
	}
}

// GetBlocks
func handleGetBlocks(request []byte, bc *BlockChain)  {
	fmt.Println("handleGetBlocks")
	var buff bytes.Buffer
	var data GetBlocks
	// 获取request数据
	dataBytes := request[12:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the getblocks cmd failed! %v\n", err)
	}
	blocks:= bc.GetBlockHashes()
	sendInv(data.AddrFrom, blocks)
}

// Inv
func handleInv(request []byte, bc *BlockChain)  {
	fmt.Println("handleInv")
	var buff bytes.Buffer
	var data Inv
	// 获取request数据
	dataBytes := request[12:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv cmd failed! %v\n", err)
	}
	blockHash := data.Hashes[0]
	sendGetData(data.AddrFrom, blockHash)
}

// getdata
func handleGetData(request []byte, bc *BlockChain)  {
	fmt.Println("handleGetData")
	var buff bytes.Buffer
	var data GetData
	// 获取request数据
	dataBytes := request[12:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv cmd failed! %v\n", err)
	}
	block := bc.GetBlock(data.ID)
	sendBlock(data.AddrFrom, block)
}

// block
func handleBlock(request []byte, bc *BlockChain)  {
	fmt.Println("handleBlock")
	var buff bytes.Buffer
	var data BlockData
	// 获取request数据
	dataBytes := request[12:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv cmd failed! %v\n", err)
	}
	// 添加到区块链中
	blockBytes := data.Block
	block := DeserializeBlock(blockBytes)
	bc.AddBlock(block)

	utxoSet := UTXOSet{bc}
	utxoSet.Update()
}