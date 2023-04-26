package tools

import (
	"encoding/hex"
	"fmt"
	"sfs-go/internal/file"
)

// ByteToString 把字节数组转换为十六进制字符串
func ByteToString(b []byte) (s string) {
	return hex.EncodeToString(b)
}

// StringToByte 十六进制字符串转换为字节数组
func StringToByte(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func GetMyAddress() string {
	return file.ReadWithFile("config/my.address")
}
