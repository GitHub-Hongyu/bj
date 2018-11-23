package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"testing"
)

func TestUTXOSet_ResetUTXOSet(t *testing.T) {
	blockchain := BlockchainObject("")
	fmt.Printf("blockchain : %v\n", blockchain)
	if nil == blockchain {
		return
	}
	utxoSet := &UTXOSet{blockchain}
	utxoSet.ResetUTXOSet()

	utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			// 游标
			c := b.Cursor()
			// 游标迭代
			// k -> 交易哈希
			// v -> 输出结构的字节数组
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("v : %v\n", DeserializeTXOutputs(v))
			}
		}
		return nil
	})
}