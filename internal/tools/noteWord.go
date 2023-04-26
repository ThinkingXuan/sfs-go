package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// GetWordList 根据randkey获取注记词
func GetWordList(randkey string, fillSize int) []string {
	entropy := calEntropy(randkey, fillSize)
	wordList := calWordList(entropy)
	return wordList
}

func GetRandKey(words []string, fillSize int) string {
	return calRandKey(words, fillSize)
}

//  calRandKey 通过注记词生成randkey
func calRandKey(words []string, fillSize int) string {
	// 拼接二级制
	binayStr := ""
	randKey := ""
	for i := 0; i < len(words); i++ {
		wordIndex := MapEnglish[words[i]]
		binayStr += fmt.Sprintf("%0*b", 11, wordIndex)
	}

	// 去除校验码
	binayStr = binayStr[:len(binayStr)-fillSize]

	for i := 0; i < len(binayStr); i = i + 8 {
		res, err := strconv.ParseInt(binayStr[i:i+8], 2, 64)
		if err != nil {
			fmt.Println("parse int err :", err)
			return ""
		}
		randKey += fmt.Sprintf("%c", byte(res))
	}
	return randKey

}

// calWordList 计算注记词
func calWordList(entropy string) []string {
	var wordList []string
	for i := 0; i < len(entropy); i = i + 11 {
		wIndex, err := strconv.ParseInt(entropy[i:i+11], 2, 64)
		if err != nil {
			fmt.Println("parse int err :", err)
			return nil
		}
		wordList = append(wordList, EnglishList[wIndex])
	}
	return wordList
}

// calEntropy 计算二进制的熵
func calEntropy(key string, fillSize int) string {

	// 获取hash
	hash := getSHA256HashCode([]byte(key))
	// 取前2个字节作为检验码
	checkCode := hash[:2]
	// 转换为二进制

	var binaryString string
	for i := 0; i < len(key); i++ {
		// 不足8位向前补0
		binaryString += fmt.Sprintf("%0*b", 8, key[i])

	}
	for i := 0; i < len(checkCode); i++ {
		binaryString += fmt.Sprintf("%0*b", 8, checkCode[i])
	}
	// checkCode取七位，所有去掉最后一位
	binaryString = binaryString[:len(binaryString)-(16-fillSize)]

	return binaryString
}

// ZeroFillByStr
// @Description: 字符串补零
// @param str :需要操作的字符串
// @param resultLen 结果字符串的长度
// @param reverse true 为前置补零，false 为后置补零
// @return string
func ZeroFillByStr(str string, resultLen int, reverse bool) string {
	if len(str) > resultLen || resultLen <= 0 {
		return str
	}
	if reverse {
		return fmt.Sprintf("%0*s", resultLen, str) //不足前置补零
	}
	result := str
	for i := 0; i < resultLen-len(str); i++ {
		result += "0"
	}
	return result
}

// FillRandKey FillKey 填充随机密钥
func FillRandKey(len int) string {
	return RandAllString(len)
}

// GetSHA256HashCode SHA256生成哈希值
func getSHA256HashCode(message []byte) string {
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

	//方法二：
	//bytes2:=sha256.Sum256(message)//计算哈希值，返回一个长度为32的数组
	//hashCode2:=hex.EncodeToString(bytes2[:])//将数组转换成切片，转换成16进制，返回字符串
	//return hashCode2
}
