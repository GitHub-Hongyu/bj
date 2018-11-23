package BLC

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net"
)

// 服务管理
// 3000 作为主节点地址
var knowNodes = []string{"localhost:3000"}

// 节点地址
var nodeAddress string
// 启动服务器
func startServer(nodeID string)  {

	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	fmt.Printf("启动服务[%s]...\n", nodeAddress)
	// 1. 监听节点
	listen, err := net.Listen(PROTOCOL,nodeAddress)
	if nil != err {
		log.Panicf("listen address of %s failed! %v\n", nodeAddress, err)
	}
	defer listen.Close()
	bc := BlockchainObject(nodeID)
	// 两个节点，主节点负责保存所有数据，钱包节点负责发送请求同步数据
	if nodeAddress != knowNodes[0] {
		// 非主节点，向主节点发送请求，同步数据
		// sendMessage(knowNodes[0], nodeAddress)
		sendVersion(knowNodes[0], bc)
	}
	// 2. 接收请求
	for {
		conn, err := listen.Accept()
		if nil != err {
			log.Panicf("accept connect failed! %v\n", err)
		}
		//request, err := ioutil.ReadAll(conn)
		//if nil != err {
		//	log.Panicf("Receive Message failed! %\n", err)
		//}
		//fmt.Printf("Receive a Message : %v\n", request)
		// 3. 处理请求
		// 单独启动一个goroutine进行请求处理
		go handleConnection(conn, bc)
	}
}

// 用于请求处理的函数
func handleConnection(conn net.Conn, bc *BlockChain)  {
	request, err := ioutil.ReadAll(conn)
	if nil != err {
		log.Panicf("Receive a Message failed! %v\n", err)
	}
	// 提取命令
	cmd := bytesToCommand(request[:12])
	fmt.Printf("Receive a Command : %s\n", cmd)

	switch cmd {
	case CMD_VERSION:
		handleVersion(request,bc)
	case CMD_GETDATA:
		handleGetData(request,bc)
	case CMD_BLOCK:
		handleBlock(request,bc)
	case CMD_GETBLOCKS:
		handleGetBlocks(request,bc)
	case CMD_INV:
		handleInv(request,bc)
	default:
		fmt.Println("Unknown command")
	}
}
