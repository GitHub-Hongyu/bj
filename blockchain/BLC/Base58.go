package BLC

import (
	"bytes"
	"fmt"
	"math/big"
)

// 实现base58编码相关

// base58字符表
var b58Alphabet = []byte("123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")

// 编码函数
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input) // 将bytes转换为bigint
	// 设置base58求余的基数
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	fmt.Printf("")
	// 余数
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		// 以余数为下标，查找base58字母表中对应的字符
		result = append(result, b58Alphabet[mod.Int64()])
	}
	// 反转切片
	Reverse(result)
	// 添加前缀
	for b := range input{
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	fmt.Printf("result : %s\n", result)
	return result
}

// 解码函数
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}
	// 去掉前缀1
	data := input[zeroBytes:]

	for _, b := range data {
		// 得到bytes数组中指定数字/字符第一次出现的索引
		charIndex := bytes.IndexByte(b58Alphabet, b)
		// 结果乘以58
		result.Mul(result, big.NewInt(58))
		// 加上余数
		result.Add(result,big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}
