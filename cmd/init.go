package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/encrypt/esdsa"
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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
