package BLC

import (
	"bytes"
	"encoding/gob"
	"github.com/labstack/gommon/log"
	"time"
)

// 区块管理相关

// 实现一个最基本的区块结构
type Block struct {
	TimeStamp 			int64		// 区块时间戳，代表区块产生时间
	Height				int64		// 区块高度(索引、号码)，代表当前区块数量
	PrevBlockHash		[]byte		// 前区块哈希
	Hash 				[]byte		// 当前区块哈希
	Txs 				[] *Transaction	// 交易数据
	Nonce 				int64		// 随机数，也就是POW运行时的动态数据
}

//创建新的区块
func NewBlock(height int64,prevBlockHash []byte,txs []*Transaction) *Block {
	var block Block
	block=Block{
		TimeStamp:time.Now().Unix(),
		Height:height,
		PrevBlockHash:prevBlockHash,
		Txs:txs,
	}
	// 通过计算生成当前区块的哈希
	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run() //钥匙(执行工作量证明算法)
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// 计算区块哈希
//func (b *Block)SetHash() {
//	timeStampBytes := IntToHex(b.TimeStamp)
//	heightBytes := IntToHex(b.Height)
//	// 拼接区块的所有属性，进行哈希
//	blockBytes := bytes.Join([][]byte{
//		heightBytes,
//		timeStampBytes,
//		b.PrevBlockHash,
//		b.Data,
//	},[]byte{})
//	hash := sha256.Sum256(blockBytes)
//	b.Hash=hash[:]
//}

// 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

// 区块结构序列化，将区块结构序列化为[]byte(字节数组)
func (block *Block)Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result) // 新建encode对象
	if err := encoder.Encode(block); nil != err {
		log.Panicf("serialize the block to byte failed! %v\n", err)
	}
	return result.Bytes()
}

// 反序列化，将字节数组结构化为区块
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	// 新建decode对象
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); nil != err {
		log.Panicf("deserialize the []byte to block failed! %v\n", err)
	}
	return &block
}

// 把区块中的所有交易结构序列化
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.Serialize())
		//txHashes = append(txHashes, tx.TxHash)
	}
	// 将区块中所有的交易拼接之后进行哈希
	//txsHash := sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	mTree := NewMerkleTree(txHashes)
	return mTree.RootNode.Data
}