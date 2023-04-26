package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/encrypt/ecc"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "sfs system init",
	Long:  `Secure file sharing system based on blockchain init, At the same time generate the public and private keys of the elliptic curve algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		ecc := ecc.NewECC("config/")

		var err error
		if keySize == "" {
			err = ecc.GenerateECC(256, nil)
		} else if keySize == "521" {
			err = ecc.GenerateECC(521, nil)
		} else if keySize == "384" {
			err = ecc.GenerateECC(384, nil)
		} else if keySize == "256" {
			err = ecc.GenerateECC(256, nil)
		} else if keySize == "224" {
			err = ecc.GenerateECC(256, nil)
		} else {
			fmt.Println("key size error: must is null、256、521、384、256")
			return
		}

		if err != nil {
			log.Println("generate key error,", err)
			return
		}

		// crate a sfs address
		address := ecc.GetAddress()
		log.Println("my sfs address is:", address)

		// upload address and public key to fabric
		pubKey, err := ecc.GetECCPublicKey()
		if err != nil {
			log.Println("get public key failure:", err)
			return
		}
		err = insertPk(address, string(pubKey))
		if err != nil {
			log.Println(err)
			log.Println("address and public key upload failure, please again init!!!!")
			return
		}
		// write address and private key to file
		file.WriteWithFile("config/my.address", address)

	},
}

var recover = &cobra.Command{
	Use:   "recover",
	Short: "recover prikey",
	Long:  "recover prikey by note word",
	Run: func(cmd *cobra.Command, args []string) {
		ecc := ecc.NewECC("config/")

		// 读取注记词
		wordListStr := file.ReadWithFile("config/WordList")
		wordList := strings.Split(wordListStr, " ")

		err := ecc.GenerateECC(0, wordList)
		if err != nil {
			log.Println("generate key error,", err)
			return
		}

		// crate a sfs address
		address := ecc.GetAddress()
		log.Println("my sfs address is:", address)

		// write address and private key to file
		file.WriteWithFile("config/my.address", address)

		log.Println("key recover succcess!")
	},
}

// 密钥长度
var keySize string

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&keySize, "size", "s", "", "ecc key size, default 256")

	initCmd.AddCommand(recover)
}

func insertPk(address, pk string) error {
	// upload address and public key to fabric
	service := sdkInit.GetInstance().InitFabric()

	_, err := service.InsertPublicKey(address, pk)
	return err
}
