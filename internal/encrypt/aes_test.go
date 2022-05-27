package encrypt

import (
	"fmt"
	"sfs-go/internal/file"
	"testing"
	"time"
)

func TestAes(t *testing.T) {
	// create a new AES struct
	aes := NewAes()

	// generate a secret key
	//key128 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
	//	0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	//}
	//key192 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
	//	0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	//	0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
	//}
	key256 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	// create a message to be encrypted
	//message := "this is a message"
	message, err := file.ReadFileBytes("D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip")
	if err != nil {
		t.Fatal(err)
	}
	nowTime := time.Now()
	// encrypt
	cipherText, err := aes.AESEncrypt([]byte(message), key256)
	if err != nil {
		t.Fatal("Aes encrypt failure:", err)
	}

	fmt.Println(time.Since(nowTime))

	// decrypt
	_, err = aes.AESDecrypt(cipherText, key256)
	if err != nil {
		t.Fatal("Aes decrypt failure:", err)
	}

	fmt.Println(time.Since(nowTime))
	// compare
	//assert.Equal(t, string(plainText), message)

}

func TestAESCTR(t *testing.T) {
	// create a new AES struct
	aes := NewAes()
	// generate a secret key
	key128 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	message, err := file.ReadFileBytes("D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip")
	if err != nil {
		t.Fatal(err)
	}
	//message := "this is a message"
	nowTime := time.Now()

	// encrypt
	cipherText, err := aes.AESCRTCrypt([]byte(message), key128)
	if err != nil {
		t.Fatal("Aes encrypt failure:", err)
	}

	fmt.Println(time.Since(nowTime))

	// decrypt
	_, err = aes.AESCRTCrypt(cipherText, key128)
	if err != nil {
		t.Fatal("Aes decrypt failure:", err)
	}

	fmt.Println(time.Since(nowTime))

}
func TestAes128Experiment(t *testing.T) {
	// create a new AES struct
	aes := NewAes()
	// generate a secret key
	key128 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}

	message, err := file.ReadFileBytes("D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip")
	if err != nil {
		t.Fatal(err)
	}
	//message := "this is a message"

	var enTime time.Duration
	var deTime time.Duration
	for i := 0; i < 3; i++ {

		nowTime := time.Now()
		// encrypt
		cipherText, err := aes.AESEncrypt(message, key128)
		if err != nil {
			t.Fatal("Aes encrypt failure:", err)
		}
		enTime += time.Since(nowTime)

		time.Sleep(time.Second * 5)
		nowTime = time.Now()
		// decrypt
		_, err = aes.AESDecrypt(cipherText, key128)
		if err != nil {
			t.Fatal("Aes decrypt failure:", err)
		}
		deTime += time.Since(nowTime)
		time.Sleep(time.Second * 5)
		t.Log("加密消耗时间：", enTime)
		t.Log("解密消耗时间：", deTime)
	}

	t.Log("平均加密时间：", enTime/3)
	t.Log("平均解密时间：", deTime/3)
}

func TestAes192Experiment(t *testing.T) {
	// create a new AES struct
	aes := NewAes()
	// generate a secret key
	key192 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
	}
	message, err := file.ReadFileBytes("D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip")
	if err != nil {
		t.Fatal(err)
	}
	//message := "this is a message"

	var enTime time.Duration
	var deTime time.Duration
	for i := 0; i < 3; i++ {

		nowTime := time.Now()
		// encrypt
		cipherText, err := aes.AESEncrypt(message, key192)
		if err != nil {
			t.Fatal("Aes encrypt failure:", err)
		}
		enTime += time.Since(nowTime)

		time.Sleep(time.Second * 5)
		nowTime = time.Now()
		// decrypt
		_, err = aes.AESDecrypt(cipherText, key192)
		if err != nil {
			t.Fatal("Aes decrypt failure:", err)
		}
		deTime += time.Since(nowTime)
		time.Sleep(time.Second * 5)
		t.Log("加密消耗时间：", enTime)
		t.Log("解密消耗时间：", deTime)
	}

	t.Log("平均加密时间：", enTime/3)
	t.Log("平均解密时间：", deTime/3)
}
func TestAes256Experiment(t *testing.T) {
	// create a new AES struct
	aes := NewAes()
	// generate a secret key
	key256 := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}

	message, err := file.ReadFileBytes("D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip")
	if err != nil {
		t.Fatal(err)
	}
	//message := "this is a message"

	var enTime time.Duration
	var deTime time.Duration
	for i := 0; i < 2; i++ {

		nowTime := time.Now()
		// encrypt
		cipherText, err := aes.AESEncrypt(message, key256)
		if err != nil {
			t.Fatal("Aes encrypt failure:", err)
		}
		enTime += time.Since(nowTime)

		time.Sleep(time.Second * 5)
		nowTime = time.Now()
		// decrypt
		_, err = aes.AESDecrypt(cipherText, key256)
		if err != nil {
			t.Fatal("Aes decrypt failure:", err)
		}
		deTime += time.Since(nowTime)
		time.Sleep(time.Second * 5)
		t.Log("加密消耗时间：", enTime)
		t.Log("解密消耗时间：", deTime)
	}

	t.Log("平均加密时间：", enTime/2)
	t.Log("平均解密时间：", deTime/2)

}
