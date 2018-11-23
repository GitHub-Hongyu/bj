package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 工作量证明管理相关

// 目标难度位数
// 代表生成的哈希值需要前targetBit位为0， 才能满足条件
const targetBit = 16

// 工作量证明
type ProofOfWork struct {
	Block *Block 		// 对指定的区块进行验证
	target *big.Int		// 大数据存储
}

// 创建一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	// 数据总长度为8位
	// 前2位都为0
	// 8-2 = 6
	// 0000 0001
	// 左移代表乘2
	// a << n a * 2^n
	// 0000 0001->
	// 0100 0000
	// 0011 1111
	// 1 * 2 ^ 6 = 64
	target := big.NewInt(1)
	// 一个字节代表8位，而且sha256生成的哈希值是一个32位的字节数组
	// 所以此处是32*8=256
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}

// 比较哈希值，开始工作量证明
func (proofOfWork *ProofOfWork)Run() ([]byte, int64)  {
	var nonce = 0; // 碰撞次数
	var hash [32]byte // 生成的哈希值
	var hashInt big.Int //存储哈希转换之后生成的数据，最终和targe数据进行比较
	for {
		dataBytes := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\rhash : %x", hash)
		hashInt.SetBytes(hash[:])
		// 难度比较
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			// 找到了符合需求的哈希值，跳出循环
			break
		}
		nonce++
	}
	fmt.Printf("\n碰撞次数: %d\n",nonce)
	fmt.Printf("last hash : %x\n", hash)
	fmt.Println("--------------------------------------")
	return hash[:], int64(nonce)
}

// 准备数据，将区块相关属性拼接到一起，返回一个字节数组
func (pow *ProofOfWork)prepareData(nonce int) []byte  {
	var data []byte
	data = bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.HashTransactions(),
		IntToHex(pow.Block.TimeStamp),
		IntToHex(pow.Block.Height),
		IntToHex(int64(nonce)),
		IntToHex(targetBit),
	}, []byte{})
	return data
}