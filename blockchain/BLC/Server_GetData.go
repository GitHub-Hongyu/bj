package BLC

// 请求指定的区块
type GetData struct {
	AddrFrom 	string	 // 从哪个地址请求
	ID 			[]byte	 // 哈希
}
