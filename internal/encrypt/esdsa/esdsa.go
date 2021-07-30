package esdsa

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"log"
	"math/big"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addreddChechsumLen = 4
	privKeyBytesLen    = 32
)

// ECC ECC结构
type ECC struct {
	GKey *GKey
}

// GKey ECC密钥
type GKey struct {
	privateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

// NewECC 创建一个ECC
func NewECC() *ECC {
	gkey, err := MakeNewKey()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &ECC{gkey}
}

// GetPrivKey 获取私钥
func (k GKey) GetPrivKey() []byte {
	d := k.privateKey.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKey := paddedAppend(privKeyBytesLen, b, d) // []bytes type
	// s := byteToString(priKey)
	return priKey
}

// GetPubKey 获取公钥
func (k GKey) GetPubKey() []byte {
	pubKey := append(k.PublicKey.X.Bytes(), k.privateKey.Y.Bytes()...) // []bytes type
	// s := byteToString(pubKey)
	return pubKey
}

// Sign 对text签名,返回加密结果，结果为数字证书r、s的序列化后拼接，然后用hex转换为string
func (e ECC) Sign(text []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, e.GKey.privateKey, text)
	if err != nil {
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, err = w.Write([]byte(string(rt) + "+" + string(st)))
	if err != nil {
		return "", err
	}
	w.Flush()
	return hex.EncodeToString(b.Bytes()), nil
}

// GetAddress 得到地址
func (e ECC) GetAddress() (address string) {
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */
	pub_bytes := e.GKey.GetPubKey()

	/* SHA256 HASH */
	//fmt.Println("1 - Perform SHA-256 hashing on the public key")
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(pub_bytes)
	pub_hash_1 := sha256_h.Sum(nil) // 对公钥进行hash256运算
	//fmt.Println(ByteToString(pub_hash_1))
	//fmt.Println("================")

	/* RIPEMD-160 HASH */
	//fmt.Println("2 - Perform RIPEMD-160 hashing on the result of SHA-256")
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil) // 对公钥hash进行ripemd160运算
	//fmt.Println(ByteToString(pub_hash_2))
	//fmt.Println("================")
	/* Convert hash bytes to base58 chech encoded sequence */
	address = b58checkencode(0x00, pub_hash_2)

	return address
}

// getSign 证书分解,通过hex解码，分割成数字证书r，s
func (e ECC) getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		err = errors.New("decrypt error," + err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		err = errors.New("decode error," + err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

// Verify 校验文本内容是否与签名一致,使用公钥校验签名和文本内容
func (e ECC) Verify(text []byte, signature string, pubKey *ecdsa.PublicKey) (bool, error) {
	rint, sint, err := e.getSign(signature)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(pubKey, text, &rint, &sint)
	return result, nil
}

// MakeNewKey 生成密钥对
func MakeNewKey() (*GKey, error) {
	var gkey GKey

	// 初始化椭圆曲线
	pubkeyCurve := elliptic.P256()
	private, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	gkey = GKey{private, private.PublicKey}
	return &gkey, nil
}

// ECCEncrypt 椭圆曲线加密
func (e ECC) ECCEncrypt(pt []byte, puk ecies.PublicKey) ([]byte, error) {
	ct, err := ecies.Encrypt(rand.Reader, &puk, pt, nil, nil)
	return ct, err
}

// ECCDecrypt 椭圆曲线解密
func (e ECC) ECCDecrypt(ct []byte, prk *ecies.PrivateKey) ([]byte, error) {
	pt, err := prk.Decrypt(ct, nil, nil)
	return pt, err
}
