package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/labstack/gommon/log"
	"math/big"
	"time"
)

// 交易管理相关
type Transaction struct {
	// 交易的唯一标识符
	TxHash 	[]byte
	// 输入列表
	Vins	[]*TxInput
	// 输出列表
	Vouts	[]*TxOutput
}

// 生成交易哈希
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if nil != err {
		log.Panicf("tx hash encoded failed! %v\n", err)
	}
	// 添加时间戳标识
	// 如果没有时间标识会导致所有coinbadse的交易哈希完全一样
	tm := time.Now().Unix()
	txHashBytes := bytes.Join([][]byte{result.Bytes(),IntToHex(tm)}, []byte{})
	// 生成交易哈希
	hash := sha256.Sum256(txHashBytes)
	tx.TxHash = hash[:]
}
// 生成coinbase交易
/*
	address : 地址
*/
func NewCoinbaseTransaction(address string) *Transaction {

	// 输入
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	// 输出
	txOutput := NewTxOutput(10, address)

	txCoinbase := Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}

	// hash
	txCoinbase.HashTransaction()
	fmt.Printf("txCoinbase-txHash : %v\n", txCoinbase.TxHash)
	return &txCoinbase
}

// 生成普通转账交易
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction, nodeID string) *Transaction {
	var txInputs []*TxInput 		// 输入
	var txOutputs []*TxOutput 		// 输出
	// 查找指定地址from的UTXO
	money, spendableUTXODic := bc.FindSpendableUTXO(from, int64(amount), txs)
	fmt.Printf("money : %d\n", money)
	// 获取钱包集合
	wallets := NewWallets(nodeID)
	// 查找到对应的钱包结构
	wallet := wallets.Wallets[from]
	for txHash, indexArray := range spendableUTXODic {
		// input(消费源)
		txHashBytes, err :=  hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string %s failed! %v\n", err)
		}
		for _, index := range indexArray {
			txInput := &TxInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}

	// 输出(转账源)

	txOutput := NewTxOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)

	// 输出(找零)
	txOutput = NewTxOutput(money - int64(amount), from)
	txOutputs = append(txOutputs, txOutput)
	tx := &Transaction{nil, txInputs, txOutputs}
	// 生成新的交易哈希
	tx.HashTransaction()
	// 对交易进行签名
	bc.SignTransaction(tx, wallet.PrivateKey)
	return tx
}

// 判断指定交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

// 交易签名
func (tx * Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction)  {
	// 处理输入
	for _, vin := range tx.Vins {

		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("ERROR:Prev transaction is not correct\n")
		}
	}
	// 提取需要签名的属性
	// 获取copy tx
	txCopy := tx.TrimmedCopy()
	for vin_id, vin := range txCopy.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 发送者
		txCopy.Vins[vin_id].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		// 签名核心函数：
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign to tx [%x] failed! %v\n", tx.TxHash, err)
		}
		// 组成签名
		signature := append(r.Bytes(), s.Bytes()...)
		// 保存签名数据
		tx.Vins[vin_id].Signature = signature
	}

}

// 交易拷贝，生成一个交易的副本，用于交易签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}

// 设置用于签名交易(交易拷贝)的哈希
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	// 置空txCopy的哈希(原来是txHash)
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// 交易序列化
func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	// 新建encoder对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to byte failed! %v\n", err)
	}
	return result.Bytes()
}
// 交易的验证
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if !tx.IsCoinbaseTransaction() {
		return true
	}
	// 检查能否找到交易
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("VERIFY ERROR : Tx is InCorrect\n")
		}
	}
	// 获取相同的交易副本
	txCopy := tx.TrimmedCopy()

	// 使用相同的椭圆获取密钥对
	curve := elliptic.P256()

	// 遍历tx输入，对每笔输入所引用的输出进行验证
	for vinId, vin := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[vinId].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 1. 需要验证数据，该数据要与签名时的数据完全一致
		txCopy.TxHash = txCopy.Hash()
		// 2. 获取r, s
			// 签名是一个数字对，r,s就代表签名
			// 将signature中的r,s值抽取
			// r,s长度相等(根据椭圆加密计算结果)
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen/2)])
		s.SetBytes(vin.Signature[(sigLen/2):])
		// 3. 获取公钥
		// 公钥是由X,Y坐标组合的生成
		// 在验证原始公钥的时候，将X,Y拆开
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen/2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen/2):])
		// 4. 组装成原始结构的公钥
		rawPublicKey := ecdsa.PublicKey{curve,&x,&y}

		// 5. 验证签名
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r,&s) {
			return false
		}
	}
	return true
}