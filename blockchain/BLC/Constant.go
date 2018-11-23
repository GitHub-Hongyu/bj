package BLC

// 保存网络服务常量
const PROTOCOL  = "tcp"



// 版本号
const NODE_VERSION = 1

/*
	version:验证当前节点的末端区块是否是最新区块
	getBlocks:从最长的链上面获取区块
	Inv:向其它节点展示当前节点有哪些区块
	getData:请求指定区块
	block:接收到新区块的时候，进行处理
*/

// version
const CMD_VERSION = "version"
const CMD_GETBLOCKS = "getblocks"
const CMD_INV = "inv"
const CMD_GETDATA = "getdata"
const CMD_BLOCK = "block"