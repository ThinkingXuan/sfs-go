package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"sfs-go/internal/fabric/fabservice"
)

const (
	configFile  = "sdkInit/config.yaml"
	initialized = false
	SimpleCC    = "sfscc"
)

/*
初始化Fabric SDK,返回客户端配置信息实例
*/

type AppGoSDK struct {
	sdk *fabsdk.FabricSDK
}

/**
 * 建立私有变量
 */
var instance *AppGoSDK

// GetInstance 获取单例对象的方法，引用传递返回
func GetInstance() *AppGoSDK {
	if instance == nil {
		instance = new(AppGoSDK)
	}
	return instance
}

// InitFabric 初始化Fabric，供外部调用
func (app *AppGoSDK) InitFabric() fabservice.ServiceSetup {

	//获取客户端实例
	channelClient, err := InstallAndInstantiateCC(app.loadSDK(), app.loadConfigure())
	if err != nil {
		log.Println(err.Error())
	}
	//开启service（链码ID+客户端实例）
	serviceSetup := fabservice.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client:      channelClient,
	}

	return serviceSetup
}

// CloseSDK 关闭SDK
func (app *AppGoSDK) CloseSDK() {
	app.sdk.Close()
}

//加载SDK
func (app *AppGoSDK) loadSDK() *fabsdk.FabricSDK {
	//设置SDK
	sdk, err := SetupSDK(configFile, initialized)
	app.sdk = sdk
	if err != nil {
		log.Println(err.Error())
	}
	return sdk
}

//加载配置信息
func (app *AppGoSDK) loadConfigure() *InitInfo {
	initInfo := &InitInfo{

		ChannelID: "mychannel",
		//ChannelConfig: os.Getenv("GOPATH") + "/src/test2/fixtures/artifacts/channel.tx",
		OrgAdmin:       "Admin",
		OrgName:        "Org1",
		OrdererOrgName: "orderer.example.com",

		ChaincodeID: SimpleCC,
		//	ChaincodeGoPath: os.Getenv("GOPATH"),
		//ChaincodePath:   "test2/chaincode/",
		UserName: "User1",
	}
	return initInfo
}
