package esdsa_bak

import (
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/magiconair/properties/assert"
	"testing"
)

// TestECCSign 测试ECC签名
func TestECCSign(t *testing.T) {
	ecc := NewECC()
	privKey := ecc.GKey.GetPrivKey()
	t.Log("My privateKey is :", ByteToString(privKey))
	pubKey := ecc.GKey.GetPubKey()
	t.Log("My publickKey is :", ByteToString(pubKey))
	address := ecc.GetAddress()
	t.Log("My address is:", address)
	text := []byte("hahahaha~!")
	signature, _ := ecc.Sign(text)
	t.Log("Signature is :", signature)
	isSuccess, _ := ecc.Verify(text, signature, &ecc.GKey.PublicKey)

	assert.Equal(t, isSuccess, true)
}

// TestECCEncrypt 测试ECC加密
func TestECCEncrypt(t *testing.T) {

	ecc := NewECC()
	// create a message to be encrypted
	message := "this is a message"

	// encrypt
	psk := ecies.ImportECDSA(ecc.GKey.privateKey)
	cipherText, err := ecc.ECCEncrypt([]byte(message), psk.PublicKey)
	if err != nil {
		t.Fatal("ecc encrypt failure", err)
	}

	// decrypt
	plainText, err := ecc.ECCDecrypt(cipherText, psk)
	if err != nil {
		t.Fatal("ecc decrypt failure", err)
	}
	// compare
	assert.Equal(t, string(plainText), message)
}
