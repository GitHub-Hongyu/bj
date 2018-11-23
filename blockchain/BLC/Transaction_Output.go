package BLC

import "bytes"

// 交易输出
type TxOutput struct {
	// 金额
	Value 			int64
	// 用户名(该UTXO的拥有者)
	Ripemd160Hash 	[]byte

}

// output身份验证
func (txOutput *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	hash160 := TransLock(address)
	return bytes.Compare(hash160, txOutput.Ripemd160Hash) == 0
}

// string转hash160
func TransLock(address string) []byte {
	pubkeyHash := Base58Decode([]byte(address))
	hash160 := pubkeyHash[1:len(pubkeyHash)-addressCheckSumLen]
	return hash160
}

// 新建output对象
func NewTxOutput(value int64, address string) *TxOutput {
	// 新建对象
	txOutput := &TxOutput{}
	hash160 := TransLock(address)
	txOutput.Value = value
	txOutput.Ripemd160Hash = hash160
	return txOutput
}