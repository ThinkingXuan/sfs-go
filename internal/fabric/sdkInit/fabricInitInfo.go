package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgName        string
	OrgAdmin       string
	OrdererOrgName string
	OrgResMgmt     *resmgmt.Client //资源管理端实例

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
}
