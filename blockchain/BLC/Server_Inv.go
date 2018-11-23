package BLC

// 向其它节点展示当前节点的区块
type Inv struct {
	AddrFrom 	string		// 当前节点的地址
	Hashes		[][]byte	// hash
}