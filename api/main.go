package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"sfs-go/internal/encrypt/pre/recrypt"
	"sfs-go/internal/tools"
)

type RekeySerialize struct {
	Fdenc           string `json:"fdenc"`
	Rk              string `json:"rk"`
	RkSign          string `json:"rk_sign"`
	XA              string `json:"xa"`
	Capsule         string `json:"capsule"`
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

type ReKey struct {
	Fdenc   []byte           `json:"fdenc"`
	Rk      *big.Int         `json:"rk"`
	XA      *ecdsa.PublicKey `json:"xa"`
	Capsule *recrypt.Capsule `json:"capsule"`
}

func main() {
	r := gin.Default()
	r.POST("/JSON", func(c *gin.Context) {
		fmt.Println("收到请求")
		var rekeySerialize RekeySerialize
		if err := c.ShouldBindJSON(&rekeySerialize); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, "参数错误, "+err.Error())
			return
		}

		// 执行重加密过程

		// key pair parse
		xa, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(rekeySerialize.XA))
		if err1 != nil {
			c.JSON(http.StatusOK, "ParsePKIXPublicKey1 err: "+err1.Error())
			return
		}

		ce, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(rekeySerialize.CapsuleE))
		if err1 != nil {
			c.JSON(http.StatusOK, "ParsePKIXPublicKey2 err: "+err1.Error())
			return
		}

		cv, err1 := x509.ParsePKIXPublicKey(tools.StringToByte(rekeySerialize.CapsuleV))
		if err1 != nil {
			c.JSON(http.StatusOK, "ParsePKIXPublicKey2 err:"+err1.Error())
			return
		}

		// deal rk sign
		rkbig := big.NewInt(1)
		rk := rkbig.SetBytes(tools.StringToByte(rekeySerialize.Rk))

		//capsule, err := recrypt.DecodeCapsule(tools.StringToByte(rekeySerialize.Capsule))
		//if err != nil {
		//	fmt.Println("DecodeCapsule err", err)
		//	return
		//}
		capsuleBig := big.NewInt(1)
		rekey := ReKey{
			Fdenc: tools.StringToByte(rekeySerialize.Fdenc),
			Rk:    rk,
			XA:    xa.(*ecdsa.PublicKey),
			Capsule: &recrypt.Capsule{
				E: ce.(*ecdsa.PublicKey),
				V: cv.(*ecdsa.PublicKey),
				S: capsuleBig.SetBytes(tools.StringToByte(rekeySerialize.CapsuleBint)),
			},
		}

		newCapsule, err := recrypt.ReEncryption(rekey.Rk, rekey.Capsule)
		if err != nil {
			c.JSON(http.StatusOK, "ReEncryption failure: "+err.Error())
			return
		}

		newXa, err2 := x509.MarshalPKIXPublicKey(rekey.XA)
		newCe, err3 := x509.MarshalPKIXPublicKey(newCapsule.E)
		NewCv, err4 := x509.MarshalPKIXPublicKey(newCapsule.V)
		if err2 != nil || err3 != nil || err4 != nil {
			fmt.Println(err1, err2, err3)
			return
		}

		encryptEntity := EncryptEntity{XA: tools.ByteToString(newXa), CapsuleE: tools.ByteToString(newCe), CapsuleV: tools.ByteToString(NewCv), CapsuleBint: tools.ByteToString(newCapsule.S.Bytes()), CapsuleBintSign: fmt.Sprintf("%d", newCapsule.S.Sign()), Fdenc: tools.ByteToString(rekey.Fdenc)}

		c.JSON(http.StatusOK, encryptEntity)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
