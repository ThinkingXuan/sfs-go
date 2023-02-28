package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sfs-go/internal/encrypt/pre/recrypt"
	"sfs-go/internal/encrypt/util"
	"sfs-go/internal/tools"
)

type ReKey struct {
	Fdenc   []byte           `json:"fdenc"`
	Rk      *big.Int         `json:"rk"`
	XA      *ecdsa.PublicKey `json:"xa"`
	Capsule *recrypt.Capsule `json:"capsule"`
}

type RekeySerialize struct {
	Fdenc  string `json:"fdenc"`
	Rk     string `json:"rk"`
	RkSign string `json:"rk_sign"`
	XA     string `json:"xa"`
	//Capsule         string `json:"capsule"`
	CapsuleE        string `json:"capsule_e"`
	CapsuleV        string `json:"capsule_v"`
	CapsuleBint     string `json:"capsule_bint"`
	CapsuleBintSign string `json:"capsule_bint_sign"`
}

type EncryptEntity struct {
	FileID            string `json:"file_id"`
	FileEncryptCipher string `json:"file_encrypt_cipher"`

	XA              string `json:"xa"`
	CapsuleE        string `json:"capsule_e"`
	CapsuleV        string `json:"capsule_v"`
	CapsuleBint     string `json:"capsule_bint"`
	CapsuleBintSign string `json:"capsule_bint_sign"`
	Fdenc           string `json:"fdenc"`
}

func main() {

	// alice (sender)
	myPriKey, err := util.GetPriKey("api/client/Alice")
	myPubKey, err := util.GetPubKey("api/client/Alice")

	// bob （receiver）
	bobPubKey, err := util.GetPubKey("config/")
	fmt.Println(bobPubKey)

	// 本地加密后，上传密文和胶囊
	fdenc, capsule, err := recrypt.Encrypt("Hello world!", myPubKey)
	if err != nil {
		fmt.Println(err)
	}

	// 本地分享时，先产生重加密密钥，上传到rk,Pubx
	rk, pubX, err := recrypt.ReKeyGen(myPriKey, bobPubKey)
	if err != nil {
		fmt.Println(err)
	}

	// 原始rekey
	rekey := ReKey{
		Fdenc:   fdenc,
		Rk:      rk,
		XA:      pubX,
		Capsule: capsule,
	}

	xa, err1 := x509.MarshalPKIXPublicKey(rekey.XA)
	ce, err2 := x509.MarshalPKIXPublicKey(rekey.Capsule.E)
	cv, err3 := x509.MarshalPKIXPublicKey(rekey.Capsule.V)

	fmt.Println("xa", xa, " ", tools.ByteToString(xa))
	fmt.Println("ce", ce, " ", tools.ByteToString(ce))
	fmt.Println("cv", cv, " ", tools.ByteToString(cv))

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println(err1, err2, err3)
	}

	//capsuleEncode, err := recrypt.EncodeCapsule(*capsule)
	//if err != nil {
	//	fmt.Println("EncodeCapsule error", err)
	//	return
	//}

	// 需要序列化传输的rekey
	rekeySerialize := RekeySerialize{
		Fdenc:  tools.ByteToString(rekey.Fdenc),
		Rk:     tools.ByteToString(rekey.Rk.Bytes()),
		RkSign: fmt.Sprintf("%d", rekey.Rk.Sign()),
		XA:     tools.ByteToString(xa),
		//Capsule:         tools.ByteToString(capsuleEncode),
		CapsuleE:        tools.ByteToString(ce),
		CapsuleV:        tools.ByteToString(cv),
		CapsuleBint:     tools.ByteToString(rekey.Capsule.S.Bytes()),
		CapsuleBintSign: fmt.Sprintf("%d", rekey.Capsule.S.Sign()),
	}

	rekeyByte, err := json.Marshal(rekeySerialize)
	if err != nil {
		fmt.Println("rekeybyte Marshal: " + err.Error())
	}

	request, error := http.NewRequest("POST", "http://localhost:8080/JSON", bytes.NewBuffer(rekeyByte))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	res, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	var fileEncryptEntity EncryptEntity
	err = json.Unmarshal(content, &fileEncryptEntity)
	if err != nil {
		fmt.Println("json error", err)
		return
	}

	// is share's file
	if len(fileEncryptEntity.CapsuleE) > 0 && len(fileEncryptEntity.CapsuleV) > 0 {
		//ReCreateKey

		// get myprikey
		myPriKey, _ := util.GetPriKey("config/")
		// get newCapule

		// key pair parse
		xa, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(fileEncryptEntity.XA))
		if err1 != nil {
			fmt.Println("ParsePKIXPublicKey1 err: " + err1.Error())
		}

		ce, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(fileEncryptEntity.CapsuleE))
		if err1 != nil {
			fmt.Println("ParsePKIXPublicKey2 err: " + err1.Error())
		}
		cv, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(fileEncryptEntity.CapsuleV))
		if err1 != nil {
			fmt.Println("ParsePKIXPublicKey2 err:" + err1.Error())
		}

		sInt := big.NewInt(1)
		sIntSign := sInt.SetBytes(tools.StringToByte(fileEncryptEntity.CapsuleBint))

		newCapsule := &recrypt.Capsule{
			E: ce.(*ecdsa.PublicKey),
			V: cv.(*ecdsa.PublicKey),
			S: sIntSign,
		}
		//
		//fmt.Println(myPriKey)
		//fmt.Println(newCapsule)
		//fmt.Println(*xa.(*ecdsa.PublicKey))
		//fmt.Println("fenc: ", fileEncryptEntity.Fdenc)

		fd, err := recrypt.Decrypt(myPriKey, newCapsule, xa.(*ecdsa.PublicKey), tools.StringToByte(fileEncryptEntity.Fdenc))
		if err != nil {
			log.Println("pre failure: ", err)
			return
		}
		fmt.Println("fd1: ", string(fd))

	}

}
