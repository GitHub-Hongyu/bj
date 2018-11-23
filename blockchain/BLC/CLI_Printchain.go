package BLC

import (
	"fmt"
	"os"
)

// 输出区块链信息
func (cli *CLI) printchain(nodeID string) {
	if !dbExists(nodeID) {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject(nodeID) // 获取区块链对象
	defer blockchain.DB.Close()
	blockchain.PrintChain()
}