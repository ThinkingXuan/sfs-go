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
