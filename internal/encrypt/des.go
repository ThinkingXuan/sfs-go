package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// Des Des加密算法
type Des struct {
}

func NewDes() *Des {
	return &Des{}
}

//// DESEncrypt DES加密
//func (d *Des) DESEncrypt(plainText []byte, key []byte) ([]byte, error) {
//
//	block, err := des.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	bs := block.BlockSize()
//	plainText = zeroPadding(plainText, bs)
//	if len(plainText)%bs != 0 {
//		return nil, errors.New("Need a multiple of the blocksize")
//	}
//	out := make([]byte, len(plainText))
//	dst := out
//	for len(plainText) > 0 {
//		block.Encrypt(dst, plainText[:bs])
//		plainText = plainText[bs:]
//		dst = dst[bs:]
//	}
//	return out, nil
//}

//// DESDecrypt DES解密
//func (d *Des) DESDecrypt(cipherText []byte, key []byte) ([]byte, error) {
//
//	block, err := des.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	out := make([]byte, len(cipherText))
//	dst := out
//	bs := block.BlockSize()
//	if len(cipherText)%bs != 0 {
//		return nil, errors.New("crypto/cipher: input not full blocks")
//	}
//	for len(cipherText) > 0 {
//		block.Decrypt(dst, cipherText[:bs])
//		cipherText = cipherText[bs:]
//		dst = dst[bs:]
//	}
//	out = zeroUnPadding(out)
//	return out, nil
//}
//
//func zeroPadding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padText := bytes.Repeat([]byte{0}, padding)
//	return append(ciphertext, padText...)
//}
//
//func zeroUnPadding(origData []byte) []byte {
//	return bytes.TrimFunc(origData,
//		func(r rune) bool {
//			return r == rune(0)
//		})
//}

// PaddingText 填充字符串（末尾）
func PaddingText(str []byte, blockSize int) []byte {
	//需要填充的数据长度
	paddingCount := blockSize - len(str)%blockSize
	//填充数据为：paddingCount ,填充的值为：paddingCount
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	newPaddingStr := append(str, paddingStr...)
	//fmt.Println(newPaddingStr)
	return newPaddingStr
}

// UnPaddingText 去掉字符（末尾）
func UnPaddingText(str []byte) []byte {
	n := len(str)
	count := int(str[n-1])
	newPaddingText := str[:n-count]
	return newPaddingText
}

// DESEncrypt DES加密
func (d *Des) DESEncrypt(plainText, key []byte) []byte {
	//1、创建并返回一个使用DES算法的cipher.Block接口
	block, _ := des.NewCipher(key)
	//2、对数据进行填充
	src1 := PaddingText(plainText, block.BlockSize())

	//3.创建一个密码分组为链接模式，底层使用des加密的blockmode接口
	iv := []byte("aaaabbbb")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//4加密连续的数据块
	desc := make([]byte, len(src1))
	blockMode.CryptBlocks(desc, src1)
	return desc
}

// DESDecrypt DES解密
func (d *Des) DESDecrypt(cipherText, key []byte) []byte {
	//创建一个block的接口
	block, _ := des.NewCipher(key)
	iv := []byte("aaaabbbb")
	//链接模式，创建blockMode接口
	blockeMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	blockeMode.CryptBlocks(cipherText, cipherText)
	//去掉填充
	newText := UnPaddingText(cipherText)
	return newText
}

// DES3Encrypt 3DES加密
func (d *Des) DES3Encrypt(src, key []byte) []byte {
	//des包下的三次加密接口
	block, _ := des.NewTripleDESCipher(key)
	src = PaddingText(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(src, src)
	return src
}

// DES3Decrypt 3DES解密
func (d *Des) DES3Decrypt(src, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(src, src)
	src = UnPaddingText(src)
	return src
}
