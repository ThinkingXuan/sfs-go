package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show all files at this address.",
	Long:  `show all files at this address.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("The files at this address are as follows!!!")
		address := file.ReadWithFile("config/my.address")
		fileBytes := queryAddressFile(address)
		var filesCC fabservice.AddressFile
		err := json.Unmarshal(fileBytes, &filesCC)
		if err != nil {
			log.Println("Unmarshal err:", err)
		}

		log.Println("-----------------id--------------", "----------name-------", "----------hash-------")

		for i := 0; i < len(filesCC.Files); i++ {
			log.Println(filesCC.FileEncrypt[i].FileID, filesCC.Files[i].FileName, filesCC.Files[i].FileHash)
		}
	},
}

var attrsss []string

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolP("list", "l", false, "list file")
}

func queryAddressFile(address string) []byte {
	service := sdkInit.GetInstance().InitFabric()
	filesBytes, err := service.QueryAddressFile(address)
	if err != nil {
		log.Println(err)
	}
	return filesBytes
}
