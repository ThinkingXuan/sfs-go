package encrypt

import (
	"os"
)

//  security attributes
const (
	general = iota
	confidential
)

const (
	AES128 = iota
	AES192
	AES256
	DES
	ThreeDES
)

func GenerateEncryptionMethod(filePath string, securityAttr int) int {
	return dynamicCalculate(filePath, securityAttr)
}

func dynamicCalculate(filePath string, securityAttr int) int {
	fileInfo, _ := os.Stat(filePath)
	// 读取文件大小
	filesize := readFileSize(fileInfo)

	switch securityAttr {
	case general:
		if filesize >= 0 && filesize <= 100 {
			return ThreeDES
		} else {
			return DES
		}
	case confidential:
		if filesize >= 0 && filesize <= 100 {
			return AES128
		} else if filesize <= 500 {
			return AES192
		} else {
			return AES256
		}
	}
	return AES256
}

// readFileSize 获取文件大小，返回单位M
func readFileSize(info os.FileInfo) int {
	fileSize := float32(info.Size() / (1024.0 * 1024.0))
	return int(fileSize)
}
