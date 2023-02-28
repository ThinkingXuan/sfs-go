package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"goRecrypt/recrypt"
	"math/big"
)

/**
hyperledger fabric chaincode
*/

// PublicKey address: publicKey
type PublicKey struct {
	PK string `json:"pk"`
}

// AddressFile address receive files
type AddressFile struct {
	FileEncrypt []EncryptEntity `json:"file_encrypt,omitempty"`
	Files       []File          `json:"files,omitempty"`
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

// File file info
type File struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
	FileDate string `json:"file_date"`
	FileHash string `json:"file_hash"`
}

type ReKey struct {
	Fdenc   []byte           `json:"fdenc"`
	Rk      *big.Int         `json:"rk"`
	XA      *ecdsa.PublicKey `json:"xa"`
	Capsule *recrypt.Capsule `json:"capsule"`
}

type RekeySerialize struct {
	Fdenc           string `json:"fdenc"`
	Rk              string `json:"rk"`
	RkSign          string `json:"rk_sign"`
	XA              string `json:"xa"`
	CapsuleE        string `json:"capsule_e"`
	CapsuleV        string `json:"capsule_v"`
	CapsuleBint     string `json:"capsule_bint"`
	CapsuleBintSign string `json:"capsule_bint_sign"`
}

type SfsCC struct {
}

func main() {
	err := shim.Start(new(SfsCC))
	if err != nil {
		fmt.Printf("chaincode start failure: %v", err)
	}
}

func (t *SfsCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("instantiate chaincode success!")
	return shim.Success(nil)
}

// Invoke invoke chaincode
func (t *SfsCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取调用链码时传递的参数内容(包括要调用的函数名及参数)
	fun, args := stub.GetFunctionAndParameters()

	fmt.Println("fun:", fun)
	fmt.Println("args:", args)

	if fun == "insert_pk" {
		return insertPkAndAddress(stub, args)
	} else if fun == "query_pk" {
		return queryPkAndAddress(stub, args)
	} else if fun == "insert_file" {
		return insertFile(stub, args)
	} else if fun == "query_file" {
		return queryFile(stub, args)
	} else if fun == "insert_addr_file" {
		return insertAddressFile(stub, args)
	} else if fun == "insert_share_address_file" {
		return insertShareAddressFile(stub, args)
	} else if fun == "query_addr_file" {
		return queryAddressFile(stub, args)
	} else {
		shim.Error("failure")
	}

	return shim.Error("failure")
}

// insertPkAndAddress insert user's public key and address to fabric
// k:v  key = address1, value = PublicKey's json
func insertPkAndAddress(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var pk = PublicKey{args[1]}

	pkBytes, _ := json.Marshal(pk)
	err := stub.PutState(args[0], pkBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// queryPkAndAddress according user's address1 to query public key
func queryPkAndAddress(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	pkBytes, _ := stub.GetState(args[0])
	var pk PublicKey
	err := json.Unmarshal(pkBytes, &pk)
	if err != nil {
		return shim.Error("PublicKey Unmarshal err: " + err.Error())
	}

	return shim.Success([]byte(pk.PK))
}

// insertFile insert file
// k:v  key = fileID, value = File's json
func insertFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var file = File{args[1], args[2], args[3], args[4], args[5]}

	// save file info
	fileBytes, err := json.Marshal(file)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], fileBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// save file encrypt entity

	return shim.Success(nil)
}

// queryFile according file id to query file info
func queryFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fileBytes, _ := stub.GetState(args[0])
	return shim.Success(fileBytes)
}

// insertAddressFile insert file to address's received file
// k:v  key = address2, value = addressFile's json
func insertAddressFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	addrFileBytes, _ := stub.GetState(args[0])

	//// add file to address files
	addrFile := AddressFile{}
	if len(addrFileBytes) > 0 {
		err := json.Unmarshal(addrFileBytes, &addrFile)
		if err != nil {
			return shim.Error("Unmarshal err:" + err.Error())
		}
	}

	encryptEntity := EncryptEntity{FileID: args[1], FileEncryptCipher: args[2]}
	addrFile.FileEncrypt = append(addrFile.FileEncrypt, encryptEntity)

	newAddrFileBytes, _ := json.Marshal(addrFile)
	// put state
	err := stub.PutState(args[0], newAddrFileBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// insertShareAddressFile
func insertShareAddressFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	addrFileBytes, _ := stub.GetState(args[0])

	//// add file to address files
	addrFile := AddressFile{}
	if len(addrFileBytes) > 0 {
		err := json.Unmarshal(addrFileBytes, &addrFile)
		if err != nil {
			return shim.Error("AddressFile Unmarshal err")
		}
	}

	// pre ReEncryption
	var rekeySerialize RekeySerialize
	err := json.Unmarshal([]byte(args[2]), &rekeySerialize)
	if err != nil {
		return shim.Error("ReKey Unmarshal failure: " + err.Error())
	}

	// key pair parse
	xa, err1 := x509.ParsePKIXPublicKey(StringToByte(rekeySerialize.XA))
	if err1 != nil {
		return shim.Error("ParsePKIXPublicKey1 err: " + err1.Error())
	}

	ce, err1 := x509.ParsePKIXPublicKey(StringToByte(rekeySerialize.CapsuleE))
	if err1 != nil {
		return shim.Error("ParsePKIXPublicKey2 err: " + err1.Error())
	}

	cv, err1 := x509.ParsePKIXPublicKey(StringToByte(rekeySerialize.CapsuleV))
	if err1 != nil {
		return shim.Error("ParsePKIXPublicKey2 err:" + err1.Error())
	}

	// deal rk sign
	rkbig := big.NewInt(1)
	rk := rkbig.SetBytes(StringToByte(rekeySerialize.Rk))

	//capsule, err := recrypt.DecodeCapsule(tools.StringToByte(rekeySerialize.Capsule))
	//if err != nil {
	//	fmt.Println("DecodeCapsule err", err)
	//	return
	//}
	capsuleBig := big.NewInt(1)

	rekey := ReKey{
		Fdenc: StringToByte(rekeySerialize.Fdenc),
		Rk:    rk,
		XA:    xa.(*ecdsa.PublicKey),
		Capsule: &recrypt.Capsule{
			E: ce.(*ecdsa.PublicKey),
			V: cv.(*ecdsa.PublicKey),
			S: capsuleBig.SetBytes(StringToByte(rekeySerialize.CapsuleBint)),
		},
	}

	newCapsule, err := recrypt.ReEncryption(rekey.Rk, rekey.Capsule)
	if err != nil {
		return shim.Error("ReEncryption failure: " + err.Error())
	}

	newXa, err2 := x509.MarshalPKIXPublicKey(rekey.XA)
	newCe, err3 := x509.MarshalPKIXPublicKey(newCapsule.E)
	NewCv, err4 := x509.MarshalPKIXPublicKey(newCapsule.V)
	if err2 != nil || err3 != nil || err4 != nil {
		fmt.Println(err1, err2, err3)
		return shim.Error(err2.Error())
	}

	encryptEntity := EncryptEntity{FileID: args[1], XA: ByteToString(newXa), CapsuleE: ByteToString(newCe), CapsuleV: ByteToString(NewCv), CapsuleBint: ByteToString(newCapsule.S.Bytes()), CapsuleBintSign: fmt.Sprintf("%d", newCapsule.S.Sign()), Fdenc: ByteToString(rekey.Fdenc)}
	addrFile.FileEncrypt = append(addrFile.FileEncrypt, encryptEntity)

	newAddrFileBytes, err := json.Marshal(addrFile)
	if err != nil {
		return shim.Error("newAddrFileBytes Marshal failure:" + err.Error())
	}

	// put state
	err = stub.PutState(args[0], newAddrFileBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// queryAddressFile according address id to query all file info
func queryAddressFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	addrFileBytes, _ := stub.GetState(args[0])

	addrFile := AddressFile{}
	err := json.Unmarshal(addrFileBytes, &addrFile)
	if err != nil {
		return shim.Error("Unmarshal err")
	}
	for i := 0; i < len(addrFile.FileEncrypt); i++ {

		fileBytes, _ := stub.GetState(addrFile.FileEncrypt[i].FileID)
		file := File{}
		_ = json.Unmarshal(fileBytes, &file)
		addrFile.Files = append(addrFile.Files, file)
	}

	newAddrFileBytes, _ := json.Marshal(addrFile)

	return shim.Success(newAddrFileBytes)
}

// ByteToString 把字节数组转换为十六进制字符串
func ByteToString(b []byte) (s string) {
	return hex.EncodeToString(b)
}

// StringToByte 十六进制字符串转换为字节数组
func StringToByte(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
