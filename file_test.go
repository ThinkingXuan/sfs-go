package main

import (
	"fmt"
	"log"
	"os"
	"sfs-go/internal/encrypt"
	"sfs-go/internal/file"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	filePath := "CIQNFPFKXWIL6JXF4UBNVKV7HZ3BCF34RSD4QV3X4EAKEWKMV6YY6TI.hex"
	// get file info detail
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Println("file no exist")
		return
	}
	t.Log(fileInfo.Size())

	fileBytes, _ := file.ReadFileBytes("CIQNFPFKXWIL6JXF4UBNVKV7HZ3BCF34RSD4QV3X4EAKEWKMV6YY6TI.hex")

	fmt.Printf("source file size: %d\n", len(fileBytes))
	// create a new AES struct
	aes := encrypt.NewAes()

	// generate a secret key
	key := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}




	//startEncrpt := time.Now()
	//for i := 0; i < 1000; i++ {
	//
	//	// encrypt
	//	_, err := aes.AESEncrypt(fileBytes, key)
	//	if err != nil {
	//		t.Fatal("Aes encrypt failure")
	//	}
	//	//t.Log(time.Since(startEncrpt))
	//}
	//t.Log(time.Since(startEncrpt))

	// encrypt
	cipherText, err := aes.AESEncrypt(fileBytes, key)
	if err != nil {
		t.Fatal("Aes encrypt failure")
	}

	fmt.Printf("encrypt file size: %d\n", len(cipherText))



	startDecrpt := time.Now()
	for i := 0; i < 1000; i++ {
		// decrypt
		_, err := aes.AESDecrypt(cipherText, key)
		if err != nil {
			t.Fatal("Aes decrypt failure")
		}
	}
	t.Log(time.Since(startDecrpt))
}

func Test11111(t *testing.T) {

	fileBytes, _ := file.ReadFileBytes("CIQNFPFKXWIL6JXF4UBNVKV7HZ3BCF34RSD4QV3X4EAKEWKMV6YY6TI.hex")
	//fileBytes, _ := file.ReadFileBytes("CIQP43Z76C75PGHHW2BTQIHG6SPTUYVFF6SCIC7XKCXMTLKRVAIV2EA.hex")


	// create a new AES struct
	aes := encrypt.NewAes()

	// generate a secret key
	key := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	//fmt.Println(string(fileBytes))

	// encrypt
	cipherText, err := aes.AESEncrypt(fileBytes, key)
	if err != nil {
		t.Fatal("Aes encrypt failure")
	}

	test, err := aes.AESDecrypt(cipherText,key)
	fmt.Println(string(test))
}