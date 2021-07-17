package encrypt

import (
	"bytes"
	"crypto/des"
	"errors"
)

// Des Des加密算法
type Des struct {
}

func NewDes() *Des {
	return &Des{}
}

// DESEncrypt DES加密
func (d *Des) DESEncrypt(plainText []byte, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	plainText = zeroPadding(plainText, bs)
	if len(plainText)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(plainText))
	dst := out
	for len(plainText) > 0 {
		block.Encrypt(dst, plainText[:bs])
		plainText = plainText[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// DESDecrypt DES解密
func (d *Des) DESDecrypt(cipherText []byte, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(cipherText))
	dst := out
	bs := block.BlockSize()
	if len(cipherText)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(cipherText) > 0 {
		block.Decrypt(dst, cipherText[:bs])
		cipherText = cipherText[bs:]
		dst = dst[bs:]
	}
	out = zeroUnPadding(out)
	return out, nil
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padText...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
