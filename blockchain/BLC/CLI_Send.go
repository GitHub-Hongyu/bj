package BLC

import (
	"fmt"
	"os"
)

// 发送交易
func (cli *CLI) send(from, to, amount []string, nodeID string)  {
	// 检测数据库
	if !dbExists(nodeID) {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchain := BlockchainObject(nodeID)
	defer blockchain.DB.Close()
	// 发起转账，产生挖矿
	blockchain.MineNewBlock(from, to, amount, nodeID)

	// 更新utxo table
	utxoSet := &UTXOSet{blockchain}
	utxoSet.Update()
}