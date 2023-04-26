package ecc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"log"
	"sfs-go/internal/encrypt/util"
	"sfs-go/internal/file"
	"strings"
	"testing"
)

func TestGenerateECCKey(t *testing.T) {
	ecc := NewECC("./")
	err := ecc.GenerateECCKey(224)
	if err != nil {
		fmt.Println(err)
	}
	//os.Exit(0)
	err = ecc.GenerateECCKey(256)
	if err != nil {
		fmt.Println(err)
	}
	//os.Exit(0)
	err = ecc.GenerateECCKey(384)
	if err != nil {
		fmt.Println(err)
	}
	//os.Exit(0)
	err = ecc.GenerateECCKey(512)
	if err != nil {
		fmt.Println(err)
	}
	//os.Exit(0)
	err = ecc.GenerateECCKey(1024)
	if err != nil {
		fmt.Println(err)
	}
}

func TestECC(t *testing.T) {
	ecc := NewECC("../../../config/")

	//plainText, _ := file.ReadFileBytes("D:\\workspace\\go_workspace\\src\\sfs-go\\chaincode\\sfscc.go"
	plainText := []byte("this is a message, this is a message, this is a message, this is a message, this is a message\nthis is a message, this is a message, this is a message, this is a message, this is a message\nthis is a message, this is a message, this is a message, this is a message, this is a message")
	t.Log("明文：\n", string(plainText))
	cipherText, err := ecc.EccEncrypt(plainText)
	if err != nil {
		fmt.Println(err)
	}
	t.Log("加密后密文: \n", cipherText)
	plainText, err = ecc.EccDecrypt(cipherText)
	if err != nil {
		fmt.Println(err)
	}
	t.Log("解密后：\n", string(plainText))

}

func TestSignVerify(t *testing.T) {
	ecc := NewECC("../../../config/")
	plainText := []byte("this is a sign, this is a sign, this is a sign, this is a sign, this is a sign\nthis is a sign, this is a sign, this is a sign, this is a sign, this is a sign\n this is a sign, this is a sign, this is a sign, this is a sign, this is a sign")
	//t.Log("签名的信息: \n", string(plainText))
	rText, sText, _ := ecc.ECCSign(plainText)
	//t.Log("签名：\n", rText, sText)
	ok, _ := ecc.ECCVerify(plainText, rText, sText)
	t.Log("验证成功？", ok)
}

func TestGetAddress(t *testing.T) {
	ecc := NewECC("./")
	address := ecc.GetAddress()
	fmt.Println(address, len(address))
}

func TestGetPukey(t *testing.T) {
	// get pem.Block
	block, err := util.GetKey("./eccPublic.pem")
	if err != nil {
		log.Println(err)
	}
	// x509
	pubInter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
	}
	// assert
	pubKey := pubInter.(*ecdsa.PublicKey)
	log.Println(pubKey)
}

func TestECC_GenerateECC(t *testing.T) {
	ecc := NewECC("./")
	wordListStr := file.ReadWithFile("WordList")
	wordList := strings.Split(wordListStr, " ")
	err := ecc.GenerateECC(521, wordList)
	if err != nil {
		fmt.Println(err)
	}
	////os.Exit(0)
	//err = ecc.GenerateECCKey(256)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////os.Exit(0)
	//err = ecc.GenerateECCKey(384)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////os.Exit(0)
	//err = ecc.GenerateECCKey(512)
	//if err != nil {
	//	fmt.Println(err)
	//}

}
