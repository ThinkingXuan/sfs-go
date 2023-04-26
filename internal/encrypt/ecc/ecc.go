package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"os"
	"path"
	"runtime"
	"sfs-go/internal/encrypt/util"
	"sfs-go/internal/file"
	"sfs-go/internal/tools"
	"strings"
)

var DesKeyError = "key size should be 8"
var DesIvError = "IV size should be 8"
var TripleDesKeyError = "key size should be 24"
var AesKeyError = "key size should be 16, 24 or 32"
var AesIvError = "IV size should be 16, 24 or 32"
var RsatransError = "error occur when trans to *rsa.Publickey"

//var RsaNilError = "error occur when decrypt"
var EcckeyError = "key size should be 224 256 384 521"

// ECC ECC结构
type ECC struct {
	dir string
}

// NewECC 创建一个ECC, dir公私钥所在路径
func NewECC(dir string) *ECC {
	return &ECC{dir: dir}
}

// GenerateECC 通过随机密钥的来产生ECC的公私钥并保存，随机密钥的长度来决定曲线的复杂性。
// randkey >=36   P224曲线
// randkey >=40   P256曲线
// randkey >=56   P384曲线
// randkey >=73   P521曲线
// 这样做的原因事保存randkey来方法注记词的构建
//
func (ecc ECC) GenerateECC(keySize int, wordList []string) error {

	randkey := ""
	randKeySize := 0
	fillSize := 0 // 需要填充的二进制位数
	// randKey不为空，取恢复随机密钥
	if len(wordList) > 0 {
		// 恢复随机密钥
		if len(wordList) == 27 {
			randkey = tools.GetRandKey(wordList, 9)
		} else if len(wordList) == 30 {
			randkey = tools.GetRandKey(wordList, 10)
		} else if len(wordList) == 41 {
			randkey = tools.GetRandKey(wordList, 3)
		} else if len(wordList) == 54 {
			randkey = tools.GetRandKey(wordList, 10)
		} else {
			return errors.New("word list error")
		}
		randKeySize = len(randkey)
	} else {
		//自己生成
		if !(keySize == 224 || keySize == 256 || keySize == 384 || keySize == 521) {
			return errors.New("key not null")
		}

		switch keySize {
		case 224:
			randKeySize = 224/8 + 8 //36 288 9
			fillSize = 9
		case 256:
			randKeySize = 256/8 + 8 //40 320 10
			fillSize = 10
		case 384:
			randKeySize = 384/8 + 8 //56 448 3
			fillSize = 3
		case 521:
			randKeySize = 521/8 + 8 //73 584 10
			fillSize = 10
		}
		randkey = tools.FillRandKey(randKeySize)
	}

	// 根据随机密匙的长度创建私匙
	var curve elliptic.Curve
	if randKeySize >= 73 {
		curve = elliptic.P521()
	} else if randKeySize >= 56 {
		curve = elliptic.P384()
	} else if randKeySize >= 40 {
		curve = elliptic.P256()
	} else if randKeySize >= 36 {
		curve = elliptic.P224()
	}
	// 生成私匙
	priKey, err := ecdsa.GenerateKey(curve, strings.NewReader(randkey+"0"))
	if err != nil {
		return errors.New("generate key error:" + err.Error())
	}

	if priKey == nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, EcckeyError)
	}
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	// x509
	derText, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	// pem block
	block := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}
	f, err := os.Create(ecc.dir + "eccPrivate.pem")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	err = pem.Encode(f, block)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	f.Close()
	// public key
	pubKey := priKey.PublicKey
	derText, err = x509.MarshalPKIXPublicKey(&pubKey)
	block = &pem.Block{
		Type:  "ecdsa public key",
		Bytes: derText,
	}
	f, err = os.Create(ecc.dir + "eccPublic.pem")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	err = pem.Encode(f, block)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	f.Close()

	// 生成并保存注记词
	if len(wordList) <= 0 {
		wordList := tools.GetWordList(randkey, fillSize)
		wordListStr := strings.Join(wordList, " ")
		file.WriteWithFile(ecc.dir+"WordList", wordListStr)
	}
	return nil
}

// GenerateECCKey
// 生成ECC私钥对
// keySize 密钥大小, 224 256 384 521
// dirPath 密钥文件生成后保存的目录
// 返回 错误
func (ecc ECC) GenerateECCKey(keySize int) error {
	// generate private key
	var priKey *ecdsa.PrivateKey
	var err error
	switch keySize {
	case 224:
		priKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case 256:
		priKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case 384:
		priKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case 521:
		priKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		priKey, err = nil, nil
	}
	if priKey == nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, EcckeyError)
	}
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	// x509
	derText, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	// pem block
	block := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}
	file, err := os.Create(ecc.dir + "eccPrivate.pem")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	err = pem.Encode(file, block)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	file.Close()
	// public key
	pubKey := priKey.PublicKey
	derText, err = x509.MarshalPKIXPublicKey(&pubKey)
	block = &pem.Block{
		Type:  "ecdsa public key",
		Bytes: derText,
	}
	file, err = os.Create(ecc.dir + "eccPublic.pem")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	err = pem.Encode(file, block)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return util.Error(file, line+1, err.Error())
	}
	file.Close()
	return nil
}

// EccEncrypt
// Ecc 加密
// plainText 明文
// filePath 公钥文件路径
// 返回 密文 错误
func (ecc ECC) EccEncrypt(plainText []byte) ([]byte, error) {

	pubKeyPemPath := path.Join(ecc.dir, "eccPublic.pem")
	// get pem.Block
	block, err := util.GetKey(pubKeyPemPath)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, err.Error())
	}
	// X509
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, err.Error())
	}
	publicKey, flag := publicInterface.(*ecdsa.PublicKey)
	if flag == false {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, RsatransError)
	}
	cipherText, err := ecies.Encrypt(rand.Reader, util.PubEcdsaToEcies(publicKey), plainText, nil, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, err.Error())
	}
	return cipherText, err
}

// EccDecrypt 解密
// cipherText 密文
// filePath 私钥文件路径
// 返回 明文 错误
func (ecc ECC) EccDecrypt(cipherText []byte) (plainText []byte, err error) {
	priKeyPemPath := path.Join(ecc.dir, "eccPrivate.pem")
	// get pem.Block
	block, err := util.GetKey(priKeyPemPath)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, err.Error())
	}
	// get privateKey
	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)
	priKey := util.PriEcdsaToEcies(privateKey)
	plainText, err = priKey.Decrypt(cipherText, nil, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, util.Error(file, line+1, err.Error())
	}
	return plainText, nil
}

// ECCSign 签名
// plainText 明文
// priPath 私钥路径
// 返回 签名结果
func (ecc ECC) ECCSign(plainText []byte) ([]byte, []byte, error) {
	priKeyPemPath := path.Join(ecc.dir, "eccPrivate.pem")
	// get pem.Block
	block, err := util.GetKey(priKeyPemPath)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, nil, util.Error(file, line+1, err.Error())
	}
	// x509
	priKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, nil, util.Error(file, line+1, err.Error())
	}
	hashText := sha256.Sum256(plainText)
	// sign
	r, s, err := ecdsa.Sign(rand.Reader, priKey, hashText[:])
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, nil, util.Error(file, line+1, err.Error())
	}
	// marshal
	rText, err := r.MarshalText()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, nil, util.Error(file, line+1, err.Error())
	}
	sText, err := s.MarshalText()
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, nil, util.Error(file, line+1, err.Error())
	}
	return rText, sText, nil
}

// ECCVerify 签名验证
// plainText 明文
// rText,sText 签名
// pubPath公钥文件路径
// 返回 验签结果 错误
func (ecc ECC) ECCVerify(plainText, rText, sText []byte) (bool, error) {
	pubKeyPemPath := path.Join(ecc.dir, "eccPublic.pem")
	// get pem.Block
	block, err := util.GetKey(pubKeyPemPath)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return false, util.Error(file, line+1, err.Error())
	}
	// x509
	pubInter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return false, util.Error(file, line+1, err.Error())
	}
	// assert
	pubKey := pubInter.(*ecdsa.PublicKey)
	hashText := sha256.Sum256(plainText)
	var r, s big.Int
	// unmarshal
	err = r.UnmarshalText(rText)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return false, util.Error(file, line+1, err.Error())
	}
	err = s.UnmarshalText(sText)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		return false, util.Error(file, line+1, err.Error())
	}
	// verify
	ok := ecdsa.Verify(pubKey, hashText[:], &r, &s)
	return ok, nil
}

func (ecc ECC) GetECCPublicKey() ([]byte, error) {
	bytes, err := file.ReadFileBytes("config/eccPublic.pem")
	return bytes, err
}

// GetAddress 获取地址
func (ecc ECC) GetAddress() (address string) {

	pubKeyPemPath := path.Join(ecc.dir, "eccPublic.pem")
	// get pem.Block
	block, err := util.GetKey(pubKeyPemPath)
	if err != nil {
		return ""
	}
	// X509
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {

		return ""
	}
	publicKey, flag := publicInterface.(*ecdsa.PublicKey)
	if flag == false {
		return ""
	}
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */
	//pub_bytes := e.GKey.GetPubKey()

	/* SHA256 HASH */
	//fmt.Println("1 - Perform SHA-256 hashing on the public key")
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(publicKey.Y.Bytes())
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

// b58chechencode encodes version ver and byte slice b into a base-58 chech encoded string.
func b58checkencode(ver uint8, b []byte) (s string) {
	/* Prepend version */
	//fmt.Println("3 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)")
	bcpy := append([]byte{ver}, b...)
	//fmt.Println(ByteToString(bcpy))
	//fmt.Println("================")

	/* Create a new SHA256 context */
	sha256H := sha256.New()

	/* SHA256 HASH #1 */
	//fmt.Println("4 - Perform SHA-256 hash on the extended PIPEMD-160 result")
	sha256H.Reset()
	sha256H.Write(bcpy)
	hash1 := sha256H.Sum(nil)
	//fmt.Println(ByteToString(hash1))
	//fmt.Println("================")

	/* SHA256 HASH #2 */
	//fmt.Println("5 - Perform SHA-256 hash on the result of the previous SHA-256 hash")
	sha256H.Reset()
	sha256H.Write(hash1)
	hash2 := sha256H.Sum(nil)
	//fmt.Println(ByteToString(hash2))
	//fmt.Println("================")

	/* Append first four bytes of hash */
	//fmt.Println("6 - Take the first 4 bytes of the second SHA-256 hash. This is the address chechsum")
	//fmt.Println(ByteToString(hash2[0:4]))
	//fmt.Println("================")

	//fmt.Println("7 - Add the 4 checksum bytes from stage 7 at the end of extended PIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.")
	bcpy = append(bcpy, hash2[0:4]...)
	//fmt.Println(ByteToString(bcpy))
	//fmt.Println("================")

	/* Encode base58 string */
	s = b58encode(bcpy)

	/* For number  of leading 0's in bytes, prepend 1 */
	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}
	//fmt.Println("8 - Convet the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format")
	//fmt.Println(s)
	//fmt.Println("================")

	return s
}

// b58encode encodea a byte slice b into a base-58 encoded string.
func b58encode(b []byte) (s string) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */
	const BITCOIN_BASE58_TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	x := new(big.Int).SetBytes(b)
	// Initialize
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x /58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(BITCOIN_BASE58_TABLE[r.Int64()]) + s
	}
	return s
}
