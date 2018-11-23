package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
)

// int64转换成字节数组
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer,binary.BigEndian,data)
	if nil != err {
		log.Panicf("int to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// 标准JSON格式转切片
// 在windows下，JSON转账成slice的标准的输入格式：
// bc.exe send -from "[\"Alice\"]" -to "[\"Bob\"]" -amount "[\"2\"]"
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	// 通过json包进行转换
	if err := json.Unmarshal([]byte(jsonString), &strSlice); err != nil {
		log.Panicf("json to []string failed! %v\n", err)
	}
	return strSlice
}

// 反转切片
func Reverse(data []byte)  {
	for i, j := 0, len(data) - 1; i < j; i, j = i + 1, j - 1 {
		data[i], data[j] = data[j], data[i]
	}
	fmt.Println(data)
}

// gob编码
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if nil != err {
		log.Panicf("encode the data failed! %v\n", err)
	}
	return buff.Bytes()
}

// 将命令转换为字节数组(长度最长为12位)
func commandToBytes(command string) []byte {
	var bytes[12]byte // 命令长度
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

// 将字节数组转换成cmd
func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}