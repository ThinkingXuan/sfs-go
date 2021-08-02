package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"path/filepath"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	file2 "sfs-go/internal/file"
)

var (
	hash     string
	fileName string
	fileID   string
	dir      string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "file download",
	Long:  `file download`,
	Run: func(cmd *cobra.Command, args []string) {

		var fileHash string
		var fileCC fabservice.File

		if fileID != "" {
			fileCC = QueryFile(fileID)
			fileHash = fileCC.FileHash
		}

		file, err := file2.DownFile(fileHash)
		if err != nil {
			log.Print("error:", err)
			return
		}

		dst := filepath.Join(dir, fileCC.FileName)
		err = ioutil.WriteFile(dst, file, 0666)
		if err != nil {
			log.Print("error:", err)
			return
		}
		log.Println("file download success")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVar(&hash, "hash", "", "file hash")
	downloadCmd.Flags().StringVarP(&fileID, "id", "i", "", "file id")
	downloadCmd.Flags().StringVar(&fileName, "name", "", "file name")
	downloadCmd.Flags().StringVarP(&dir, "dir", "d", "", "download file to dir")
	//_ = downloadCmd.MarkFlagRequired("hash")
}

func QueryFile(fileID string) fabservice.File {
	service := sdkInit.GetInstance().InitFabric()
	fileBytes, err := service.QueryFile(fileID)
	if err != nil {
		log.Println("query err:", err)
	}
	var file fabservice.File
	err = json.Unmarshal(fileBytes, &file)
	if err != nil {
		log.Println("Unmarshal err:", err)
	}
	return file
}
