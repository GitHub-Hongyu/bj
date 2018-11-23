package BLC

import "crypto/sha256"

// 实现Merkle树相关功能

// merkle树结构
type  MerkleTree struct {
	// 根节点
	RootNode 	*MerkleNode
}

// merkle节点结构
type MerkleNode struct {
	// 左子节点
	Left 	*MerkleNode
	// 右子节点
	Right 	*MerkleNode
	// 数据(存储该节点哈希值)
	Data 	[]byte
}

// 创建Merkle树
// 包含的就是一个根节点,里面保存了当前区块中所有的交易
func NewMerkleTree(datas [][]byte) *MerkleTree  {
	var nodes []MerkleNode // 保存节点
	// 判断交易数据条数,如果是奇数条,则把最后一条拷贝一份
	if len(datas) % 2 != 0 {
		datas = append(datas, datas[len(datas) - 1])
	}
	// 遍历所有交易数据, 创建叶子节点
	for _, data := range datas {
		node := NewMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}
	// 创建上级节点(非叶子节点)
	for i := 0; i < len(datas) / 2; i++ {
		var parentNodes []MerkleNode // 父节点列表
		for j := 0; j < len(nodes); j+=2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			parentNodes = append(parentNodes, *node)
		}
		if len(parentNodes) % 2 != 0 {
			parentNodes = append(parentNodes, parentNodes[len(parentNodes) - 1])
		}
		// 最终,nodes列表只保存根节点哈希值
		nodes = parentNodes
	}
	mtree := MerkleTree{&nodes[0]}
	return &mtree
}


// 创建Merkle节点
func NewMerkleNode(left, rigth *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}
	// 叶子节点
	if left == nil && rigth == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		// 非叶子节点,保存左子节点哈希,右子节点哈希合到一起之后,再哈希
		prevHashes := append(left.Data, rigth.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = rigth
	return node
}