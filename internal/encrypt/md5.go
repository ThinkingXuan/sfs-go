package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5Check MD5检查
func MD5Check(content, encrypted string) bool {
	return strings.EqualFold(MD5Encode(content), encrypted)
}

// MD5Encode MD5编码
func MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

