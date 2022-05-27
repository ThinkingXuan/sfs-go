package fabric

import (
	"fmt"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
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
