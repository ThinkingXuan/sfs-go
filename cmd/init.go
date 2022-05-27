package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/encrypt/ecc"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
	"time"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "sfs system init",
	Long:  `Secure file sharing system based on blockchain init, At the same time generate the public and private keys of the elliptic curve algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		var allTime time.Duration
		for i := 0; i < 10; i++ {

			startTime := time.Now()

			ecc := ecc.NewECC("config/")
			// generate public key and private
			err := ecc.GenerateECCKey(256)
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

			elapsed := time.Since(startTime)
			fmt.Println("初始化过程所花费的时间：", elapsed)
			allTime += elapsed
		}
		fmt.Println("总时间：", allTime)
		log.Println("address and public key upload success!!!!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func insertPk(address, pk string) error {
	// upload address and public key to fabric
	service := sdkInit.GetInstance().InitFabric()

	_, err := service.InsertPublicKey(address, pk)
	return err
}
