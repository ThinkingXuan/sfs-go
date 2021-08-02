package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
)

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
	_, err := service.InsertAddressFile(address, fileID)
	if err != nil {
		log.Println("insert err:", err)
		return false
	}

	return true
}
