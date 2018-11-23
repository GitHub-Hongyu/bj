package BLC

import (
	"bytes"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net"
)

// 区块版本验证
func sendVersion(toAddress string, bc *BlockChain)  {
	// 在比特币中，消息是底层的比特序列
	// 前12个字节指定命令名称(version)
	// 后面的字节包含gob编码的相关消息结构
	heigth := bc.GetHeight()
	// 得到数据
	data := gobEncode(Version{NODE_VERSION, int(heigth), nodeAddress})
	// 发送version命令
	request := append(commandToBytes(CMD_VERSION), data...)
	sendMessage(toAddress, request)
}


// 节点发送请求
func sendMessage(to string, msg []byte)  {
	fmt.Println("向服务器发送请求")
	// 1. 连接服务器
	conn, err := net.Dial(PROTOCOL, to)
	if nil != err {
		log.Panicf("connect to server [%s] failed! %v\n", to, err)
	}
	defer conn.Close()
	// 要发送的数据添加到请求中
	_, err = io.Copy(conn, bytes.NewReader(msg))
	if nil != err {
		log.Panicf("add the data failed! %v\n", err)
	}
}

// 从指定节点同步数据
func sendGetBlocks(toAddress string)  {
	data := gobEncode(GetBlocks{AddrFrom:nodeAddress})
	request := append(commandToBytes(CMD_GETBLOCKS), data...)
	sendMessage(toAddress, request)
}

// 向其它节点展示区块信息
func sendInv(toAddress string, hashes [][]byte)  {
	data := gobEncode(Inv{AddrFrom:nodeAddress, Hashes:hashes})
	request := append(commandToBytes(CMD_INV), data...)
	sendMessage(toAddress, request)
}

// 发送指定区块请求
func sendGetData(toAddress string, hash []byte)  {
	data := gobEncode(GetData{AddrFrom:nodeAddress, ID:hash})
	request := append(commandToBytes(CMD_GETDATA), data...)
	sendMessage(toAddress, request)
}

// 发送区块信息
func sendBlock(toAddress string, block []byte)  {
	data := gobEncode(BlockData{AddrFrom:nodeAddress, Block:block})
	request := append(commandToBytes(CMD_BLOCK), data...)
	sendMessage(toAddress, request)
}