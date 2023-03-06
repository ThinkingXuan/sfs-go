package main

import (
	"encoding/json"
	"fmt"
	"sfs-go/cmd"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/tools"
	"testing"
)

func TestFabricConn(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()
	fmt.Println(service)
}

func TestFabricQueryPk(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()
	pkBytes, err := service.GetPublicKey("1234567890")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("pk:", string(pkBytes))
}

func TestFabricInsertPk(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.InsertPublicKey("1234567234890", "32479238894027398478293")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricInsertFile(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	f := fabservice.File{
		FileID:   "file_id_3",
		FileName: "FileName3",
		FileType: "FileType3",
		FileSize: "FileSize3",
		FileDate: "FileDate3",
		FileHash: "FileHash3",
	}

	pkBytes, err := service.InsertFile(f)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricQueryFile(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.QueryFile("file_id_1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricInsertAddrFile(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.InsertAddressFile("1234567890", "file_id_2", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricQueryAddrFile(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.QueryAddressFile("1234567890")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricInsertAttrsAndID(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	attrs := cmd.IdAndAttr{
		Id:    "123",
		Attrs: []string{"aBE", "234"},
	}
	attrBytes, _ := json.Marshal(attrs)

	pkBytes, err := service.InsertAbeAttrsAndId("werwecewfewfwef", tools.ByteToString(attrBytes))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricQueryAttrsAndID(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.QueryAbeAttrsAndId("werwecewfewfwef")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricInsertAuthPK(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	//attrs := cmd.IdAndAttr{
	//	Id:    "123",
	//	Attrs: []string{"aBE", "234"},
	//}
	//attrBytes, _ := json.Marshal(attrs)

	pkBytes, err := service.InsertAbeAuthPK("werwecewfewfwef", []byte("234324234"))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestFabricQueryAuthPK(t *testing.T) {
	service := sdkInit.GetInstance().InitFabric()

	pkBytes, err := service.QueryAbeAuthPK("1G35GKZ1tqkfGRisvRiw6H3xMCmZpTzA46")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(pkBytes))
}

func TestJSON(t *testing.T) {

	type student struct {
		Name string `json:"name"`
	}

	var stu = student{"7b224964223a226175746831222c224174747273223a5b2261757468313a6174312061757468313a617432225d7d"}
	stuBytes, err := json.Marshal(stu)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stuBytes)
	var stu2 student
	err = json.Unmarshal(stuBytes, &stu2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stu2)
}
