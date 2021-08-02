package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/encrypt/esdsa"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "sfs system init",
	Long:  `Secure file sharing system based on blockchain init, At the same time generate the public and private keys of the elliptic curve algorithm`,
	Run: func(cmd *cobra.Command, args []string) {
		ecc := esdsa.NewECC()
		// create a private key
		priKey := ecc.GKey.GetPrivKey()
		log.Println("my private key is:", esdsa.ByteToString(priKey))
		// create a public key
		pubKey := ecc.GKey.GetPubKey()
		log.Println("my public key is:", esdsa.ByteToString(pubKey))
		// crate a sfs address
		address := ecc.GetAddress()
		log.Println("my sfs address is:", address)

		// upload address and public key to fabric
		err := insertPk(address, esdsa.ByteToString(pubKey))
		if err != nil {
			log.Println(err)
			log.Println("address and public key upload failure, please again init!!!!")
		}
		// write address and private key to file
		file.WriteWithFile("config/my.address", address)
		file.WriteWithFile("config/my.privatekey", esdsa.ByteToString(priKey))

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
