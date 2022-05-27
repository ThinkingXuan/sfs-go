package tools

import "encoding/hex"

// ByteToString 把字节数组转换为十六进制字符串
func ByteToString(b []byte) (s string) {
	return hex.EncodeToString(b)
}

// StringToByte 十六进制字符串转换为字节数组
func StringToByte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
