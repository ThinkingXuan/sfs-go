package fabservice

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

type Document struct {
	ObjectType     string `json:"docType"`
	Id             string `json:"id"`
	SentAccount    string `json:"sent_account"`
	ReceiveAccount string `json:"receive_account"`
	Department     string `json:"department"`
	IPFSHash       string `json:"ipfs_hash"`
	UserID         string `json:"user_id"`
	FileName       string `json:"file_name"`
	Size           string `json:"file_size"`
	Type           string `json:"file_type"`
	SentDate       string `json:"upload_date"`
	Classification string `json:"classification"`
	Approver       string `json:"approver"`
	SubApprover    string `json:"sub_approver"`
	Class          string `json:"class"`
}

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}
