package cmd

import (
	"crypto/ecdsa"
	"crypto/x509"
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
	"sfs-go/internal/tools"
)

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

		err = judgeFile(shareFileID)
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
	//if !addressExist(address) {
	//	return errors.New("address not exist")
	//}

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

	fmt.Println("fd: ", string(fd))

	// pre process

	// alice (sender)
	myPriKey, err := util.GetPriKey("config/")
	myPubKey, err := util.GetPubKey("config/")

	// bob （receiver）
	bobPubKeyByte, err := service.GetPublicKey(address)
	bobPubKey := util.KeyFromByte(bobPubKeyByte)
	fmt.Println(bobPubKey)

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

	// 原始rekey
	rekey := ReKey{
		Fdenc:   fdenc,
		Rk:      rk,
		XA:      pubX,
		Capsule: capsule,
	}

	xa, err1 := x509.MarshalPKIXPublicKey(rekey.XA)
	ce, err2 := x509.MarshalPKIXPublicKey(rekey.Capsule.E)
	cv, err3 := x509.MarshalPKIXPublicKey(rekey.Capsule.V)

	fmt.Println("xa", xa, " ", tools.ByteToString(xa))
	fmt.Println("ce", ce, " ", tools.ByteToString(ce))
	fmt.Println("cv", cv, " ", tools.ByteToString(cv))

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println(err1, err2, err3)
	}

	// 需要序列化传输的rekey
	rekeySerialize := RekeySerialize{
		Fdenc:           tools.ByteToString(rekey.Fdenc),
		Rk:              tools.ByteToString(rekey.Rk.Bytes()),
		RkSign:          fmt.Sprintf("%d", rekey.Rk.Sign()),
		XA:              tools.ByteToString(xa),
		CapsuleE:        tools.ByteToString(ce),
		CapsuleV:        tools.ByteToString(cv),
		CapsuleBint:     tools.ByteToString(rekey.Capsule.S.Bytes()),
		CapsuleBintSign: fmt.Sprintf("%d", rekey.Capsule.S.Sign()),
	}

	rekeyByte, err := json.Marshal(rekeySerialize)
	if err != nil {
		fmt.Println("rekeybyte Marshal: " + err.Error())
	}

	_, err = service.InsertShareAddressFile(address, fileID, rekeyByte)
	if err != nil {
		log.Println("insert err:", err)
		return false
	}
	return true
}
