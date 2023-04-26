package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"log"
	"path"
	"sfs-go/internal/encrypt/pre/recrypt"
	"sfs-go/internal/encrypt/pre/util"
	"testing"
	"time"
)

func TestRecrypt(t *testing.T) {
	// Alice Generate Alice key-pair
	aPriKey, err := getPriKey("ecc/a")
	aPubKey, err := getPubKey("ecc/a")
	// Bob Generate Bob key-pair

	bPriKey, err := getPriKey("ecc/b")
	bPubKey, err := getPubKey("ecc/b")
	if err != nil {
		log.Println(err)
		return
	}

	startTime := time.Now()
	// plain text
	m := "Hello, Proxy Re-Encryption Hello, Proxy Re-Encryption Hello, Proxy Re-Encryption Hello, Proxy Re-Encryption"
	fmt.Println("origin message:", m)
	// Alice encrypts to get cipherText and capsule

	// 本地加密后，上传密文和胶囊
	cipherText, capsule, err := recrypt.Encrypt(m, aPubKey)
	if err != nil {
		fmt.Println("1 ", err)
	}
	fmt.Println("Encrypt", time.Since(startTime))
	// 本地分享时，先产生重加密密钥，上传到rk,Pubx
	rk, pubX, err := recrypt.ReKeyGen(aPriKey, bPubKey)
	if err != nil {
		fmt.Println("2 ", err)
	}
	fmt.Println("ReKeyGen", time.Since(startTime))
	// fabric执行重加密过程，使用了了rk和胶囊，生成新的胶囊
	newCapsule, err := recrypt.ReEncryption(rk, capsule)
	if err != nil {
		fmt.Println("3 ", err.Error())
	}

	fmt.Println("ReEncryption", time.Since(startTime))

	// Bob decrypts the cipherText
	// 对方解密时
	plainText, err := recrypt.Decrypt(bPriKey, newCapsule, pubX, cipherText)
	if err != nil {
		fmt.Println("4 ", err)
	}
	fmt.Println("Decrypt", time.Since(startTime))
	// 自己解密
	//plainTextByMyPri, err := recrypt.DecryptOnMyPriKey(aPriKey, capsule, cipherText)
	//if err != nil {
	//	fmt.Println("5 ",err)
	//}
	//fmt.Println("PlainText by my own private key:", string(plainTextByMyPri))
	// get plainText
	fmt.Println("plainText:", string(plainText))
}

func getPubKey(dir string) (*ecdsa.PublicKey, error) {
	pubKeyPemPath := path.Join(dir, "224_p.pem")
	// get pem.Block
	block, err := util.GetKey(pubKeyPemPath)
	if err != nil {
		return nil, err
	}
	// X509
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {

		return nil, err
	}
	publicKey, _ := publicInterface.(*ecdsa.PublicKey)
	return publicKey, nil

}

func getPriKey(dir string) (*ecdsa.PrivateKey, error) {
	priKeyPemPath := path.Join(dir, "224.pem")
	// get pem.Block
	block, err := util.GetKey(priKeyPemPath)
	if err != nil {
		return nil, err
	}
	// x509
	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priKey, nil

}
