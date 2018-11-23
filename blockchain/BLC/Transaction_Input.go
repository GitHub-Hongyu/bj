package BLC

import "bytes"

type TxInput struct {
	// 交易哈希(不是当前交易的哈希，而是引入的上一笔交易的哈希)
	TxHash 		[]byte

	// 引用的上一笔交易的输出的索引
	Vout		int

	// 数字签名
	Signature 	[]byte

	// 公钥
	PublicKey 	[]byte
}

// 权限判断
// address:要查找余额的地址
func (in *TxInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	// 获取input的ripemd160哈希
	inputRipemd160 := Ripemd160Hash(in.PublicKey)
	return bytes.Compare(inputRipemd160, ripemd160Hash) == 0
}