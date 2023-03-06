package fabservice

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) AddDoc(info Document) (string, error) {

	eventID := "eventAddDoc"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将med对象序列化成为字节数组
	b, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("指定的info对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addDoc", Args: [][]byte{b, []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) GetInfo(id string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getInfo", Args: [][]byte{[]byte(id)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

//peer chaincode invoke -n mycc -c '{"Args":["insert_pk","1234567890","0987654321"]}' -C myc

func (t *ServiceSetup) InsertPublicKey(address string, publicKey string) ([]byte, error) {

	newAddress := address + "0"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_pk", Args: [][]byte{[]byte(newAddress), []byte(publicKey)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) GetPublicKey(address string) ([]byte, error) {
	newAddress := address + "0"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "query_pk", Args: [][]byte{[]byte(newAddress)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

//peer chaincode invoke -n mycc -c '{"Args":["insert_file","fileID1","string1","string2","string3","string4","string5"]}' -C myc

// InsertFile insert file
func (t *ServiceSetup) InsertFile(file File) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_file", Args: [][]byte{[]byte(file.FileID), []byte(file.FileName), []byte(file.FileType), []byte(file.FileSize), []byte(file.FileDate), []byte(file.FileHash)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}
	return response.Payload, nil
}

//peer chaincode invoke -n mycc -c '{"Args":["query_file","fileID1"]}' -C myc

// QueryFile query file by file id
func (t *ServiceSetup) QueryFile(fileID string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "query_file", Args: [][]byte{[]byte(fileID)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

//peer chaincode invoke -n mycc -c '{"Args":["insert_addr_file","12345678901","fileID1","dsjfoijio"]}' -C myc

// InsertAddressFile to address send file
func (t *ServiceSetup) InsertAddressFile(address string, fileID string, fileEncryptCipher string) ([]byte, error) {
	newAddress := address + "1"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_addr_file", Args: [][]byte{[]byte(newAddress), []byte(fileID), []byte(fileEncryptCipher)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}
	return response.Payload, nil
}

// InsertShareAddressFile to address send file
func (t *ServiceSetup) InsertShareAddressFile(address string, fileID string, reKey []byte) ([]byte, error) {
	newAddress := address + "1"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_share_address_file", Args: [][]byte{[]byte(newAddress), []byte(fileID), reKey}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}
	return response.Payload, nil
}

//peer chaincode invoke -n mycc -c '{"Args":["query_addr_file","12345678901"]}' -C myc

// QueryAddressFile query address's file
func (t *ServiceSetup) QueryAddressFile(address string) ([]byte, error) {
	newAddress := address + "1"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "query_addr_file", Args: [][]byte{[]byte(newAddress)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

// InsertAbeAttrsAndId cp-abe insert attr and id
func (t *ServiceSetup) InsertAbeAttrsAndId(address string, attrsAndID string) ([]byte, error) {
	newAddress := address + "2"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_attrs_id", Args: [][]byte{[]byte(newAddress), []byte(attrsAndID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

// QueryAbeAttrsAndId
func (t *ServiceSetup) QueryAbeAttrsAndId(address string) ([]byte, error) {
	newAddress := address + "2"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "query_attrs_id", Args: [][]byte{[]byte(newAddress)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

// InsertAbeAuthPK cp-abe insert auth pk
func (t *ServiceSetup) InsertAbeAuthPK(address string, authPK []byte) ([]byte, error) {
	newAddress := address + "3"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_abe_auth_pk", Args: [][]byte{[]byte(newAddress), []byte(authPK)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) QueryAbeAuthPK(address string) ([]byte, error) {
	newAddress := address + "3"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "query_abe_auth_pk", Args: [][]byte{[]byte(newAddress)}}
	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}
	return response.Payload, nil
}

func (t *ServiceSetup) InsertAbeShareAddressFile(address string, fileID string, fdenc string) ([]byte, error) {
	newAddress := address + "1"
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "insert_abe_share_address_file", Args: [][]byte{[]byte(newAddress), []byte(fileID), []byte(fdenc)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}
	return response.Payload, nil
}
