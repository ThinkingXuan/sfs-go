package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
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
	FileID []string `json:"file_id,omitempty"`
	Files  []File   `json:"files,omitempty"`
}

// File file info
type File struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
	FileData string `json:"file_data"`
	FileHash string `json:"file_hash"`
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
	return shim.Success(pkBytes)
}

// insertFile insert file
// k:v  key = fileID, value = File's json
func insertFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var file = File{args[1], args[2], args[3], args[4], args[5]}

	fileBytes, _ := json.Marshal(file)
	err := stub.PutState(args[0], fileBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

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
// k:v  key = address2, value = fileID
func insertAddressFile(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	addrFileBytes, _ := stub.GetState(args[0])

	//// add file to address files
	addrFile := AddressFile{}
	if len(addrFileBytes) > 0 {
		err := json.Unmarshal(addrFileBytes, &addrFile)
		if err != nil {
			return shim.Error("Unmarshal err")
		}
	}

	addrFile.FileID = append(addrFile.FileID, args[1])

	newAddrFileBytes, _ := json.Marshal(addrFile)
	// put state
	err := stub.PutState(args[0], newAddrFileBytes)
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
	for i := 0; i < len(addrFile.FileID); i++ {

		fileBytes,_ := stub.GetState(addrFile.FileID[i])
		file := File{}
		json.Unmarshal(fileBytes, &file)
		addrFile.Files = append(addrFile.Files, file)
	}

	newAddrFileBytes, _ := json.Marshal(addrFile)

	return shim.Success(newAddrFileBytes)
}
