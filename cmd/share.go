package cmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"sfs-go/internal/encrypt/pre/recrypt"
	"sfs-go/internal/encrypt/util"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
)

type ReKey struct {
	Fdenc   []byte          `json:"fdenc"`
	Rk      big.Int         `json:"rk"`
	XA      ecdsa.PublicKey `json:"xa"`
	Capsule recrypt.Capsule `json:"capsule"`
}

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share file to a address.",
	Long:  `share file to a address.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := judgeAddress(shareAddress)
		if err != nil {
			log.Println(err)
			return
		}

		err = judgeFile(fileID)
		if err != nil {
			log.Println(err)
			return
		}

		if !shareFile(shareAddress, shareFileID) {
			log.Println("share file failure!!!")
			return
		}

		log.Println("share file success!!!")
	},
}

var (
	shareAddress string
	shareFileID  string
)

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringVarP(&shareAddress, "address", "a", "", "file share address")
	shareCmd.Flags().StringVarP(&shareFileID, "fid", "i", "", "file id")

	_ = shareCmd.MarkFlagRequired("address")
	_ = shareCmd.MarkFlagRequired("fid")
}

func judgeAddress(address string) error {

	// myself address
	myAddress := file.ReadWithFile("config/my.address")
	if address == myAddress {
		return errors.New("can't to share file to a myself address")
	}
	// address  correct
	if len(address) != 34 {
		return errors.New("address error, the length of correct address is 34")
	}

	// address exist
	if !addressExist(address) {
		return errors.New("address not exist")
	}

	return nil
}

func judgeFile(fileID string) error {

	// address  correct
	if len(fileID) != len("be3e742d-7c84-af78-2311-187cf73487bc") {
		return errors.New("file id error")
	}

	// address exist
	if !fileExist(fileID) {
		return errors.New("file not exist")
	}
	return nil
}

func addressExist(address string) bool {
	service := sdkInit.GetInstance().InitFabric()
	pkBytes, err := service.GetPublicKey(address)
	if err != nil {
		log.Println("query err:", err)
		return false
	}
	if len(pkBytes) <= 0 {
		return false
	}
	return true
}

func fileExist(fileID string) bool {
	service := sdkInit.GetInstance().InitFabric()
	fileBytes, err := service.QueryFile(fileID)
	if err != nil {
		log.Println("query err:", err)
		return false
	}
	if len(fileBytes) <= 0 {
		return false
	}
	return true
}

func shareFile(address, fileID string) bool {
	service := sdkInit.GetInstance().InitFabric()

	//var fileHash string
	var fileEncryptEntity fabservice.EncryptEntity

	// get file info
	if fileID != "" {
		_, fileEncryptEntity = QueryFile(fileID)
		//fileHash = fileCC.FileHash
	} else {
		log.Println("file id not is empty")
		return false
	}

	// get Aes encrypt key
	fd, err := GetAESKey(fileEncryptEntity.FileEncryptCipher)
	if err != nil {
		log.Println("get aes encrypt key failure,", err)
		return false
	}

	// pre process

	// alice
	myPriKey, err := util.GetPriKey("config/")
	myPubKey, err := util.GetPubKey("config/")

	// bob
	bobPubKeyByte, err := service.GetPublicKey(address)
	bobPubKey := util.KeyFromByte(bobPubKeyByte)

	// 本地加密后，上传密文和胶囊
	fdenc, capsule, err := recrypt.Encrypt(string(fd), myPubKey)
	if err != nil {
		fmt.Println(err)
	}

	// 本地分享时，先产生重加密密钥，上传到rk,Pubx
	rk, pubX, err := recrypt.ReKeyGen(myPriKey, bobPubKey)
	if err != nil {
		fmt.Println(err)
	}

	rekey := ReKey{
		Fdenc:   fdenc,
		Rk:      *rk,
		XA:      *pubX,
		Capsule: *capsule,
	}
	rekeyByte, err := json.Marshal(rekey)

	_, err = service.InsertShareAddressFile(address, fileID, rekeyByte)
	if err != nil {
		log.Println("insert err:", err)
		return false
	}
	return true
}
