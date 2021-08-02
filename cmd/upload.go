package cmd

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
	"time"
)

var (
	filePath string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Print("filepath: ", filePath)

		// get file info detail
		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			log.Println("file no exist")
			return
		}

		// upload file to ipfs
		hash, err := file.UploadFileToIPFS(filePath)
		if err != nil {
			log.Printf("file[%s] upload failure", filePath)
		}

		// upload file to fabric
		fileID, _ := uuid.GenerateUUID()
		fileCC := fabservice.File{
			FileID:   fileID,
			FileName: fileInfo.Name(),
			FileType: path.Ext(fileInfo.Name()),
			FileSize: fmt.Sprintf("%d", fileInfo.Size()),
			FileDate: time.Now().Format("2006-01-02 15:04:05"),
			FileHash: hash,
		}
		insertFileToFabric(fileCC)

		log.Println("file upload success!!!!")

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&filePath, "path", "p", "", "file upload path")
	_ = uploadCmd.MarkFlagRequired("path")
}

func insertFileToFabric(fileCC fabservice.File) {
	service := sdkInit.GetInstance().InitFabric()
	_, err := service.InsertFile(fileCC)
	address := file.ReadWithFile("config/my.address")
	_, err = service.InsertAddressFile(address, fileCC.FileID)
	if err != nil {
		log.Println("failure")
	}
}
