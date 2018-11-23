package BLC


// 代表当前区块版本信息(决定是否需要同步)
type Version struct {
	Version		int		// 版本
	Height		int		// 当前节点的区块高度
	AddrFrom	string	// 当前节点地址
}
