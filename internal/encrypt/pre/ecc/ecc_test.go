package ecc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"log"
	"sfs-go/internal/encrypt/pre/file"
	"sfs-go/internal/encrypt/pre/util"
	"testing"
)

func TestGenerateECCKey(t *testing.T) {
	ecc := NewECC("b/")
	//err := ecc.GenerateECCKey(224)
	//if err != nil {
	//	fmt.Println(err)
	//}
	////os.Exit(0)
	err := ecc.GenerateECCKey(256)
	if err != nil {
		fmt.Println(err)
	}
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
	////os.Exit(0)
	//err = ecc.GenerateECCKey(1024)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func TestEcc(t *testing.T) {
	ecc := NewECC("../../../config/")

	plainText, _ := file.ReadFileBytes("D:\\workspace\\go_workspace\\src\\sfs-go\\chaincode\\sfscc.go")
	cipherText, err := ecc.EccEncrypt(plainText)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("加密后：%s\n", cipherText)
	plainText, err = ecc.EccDecrypt(cipherText)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("解密后：%s\n", plainText)
}

func TestSignVerify(t *testing.T) {
	ecc := NewECC("./")
	plainText := []byte("张华考上了北京大学；李萍进了中等技术学校；我在百货公司当售货员：我们都有美好的未来")
	rText, sText, _ := ecc.ECCSign(plainText)
	ok, err := ecc.ECCVerify(plainText, rText, sText)
	fmt.Println(err)
	fmt.Printf("验证成功？ %t", ok)
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
