package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
)

// Aes 加密算法
type Aes struct {
}

func NewAes() *Aes {
	return &Aes{}
}

// AESEncrypt AES加密
func (a *Aes) AESEncrypt(plainText []byte, key []byte) ([]byte, error) {
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(plainText))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, plainText)
	return encrypted, nil
}

// AESDecrypt AES解密
func (a *Aes) AESDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(cipherText))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, cipherText)
	return decrypted, nil
}

// AESCRTCrypt AEC加密和解密（CRT模式）
func (a *Aes) AESCRTCrypt(text []byte, key []byte) ([]byte, error) {
	//指定加密、解密算法为AES，返回一个AES的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//指定计数器,长度必须等于block的块尺寸
	count := []byte("12345678abcdefgh")
	//指定分组模式
	blockMode := cipher.NewCTR(block, count)
	//执行加密、解密操作
	message := make([]byte, len(text))
	blockMode.XORKeyStream(message, text)
	//返回明文或密文
	return message, nil
}
