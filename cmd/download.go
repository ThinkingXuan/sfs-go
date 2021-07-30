package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"path/filepath"
	file2 "sfs-go/internal/file"
)

var (
	hash string
	dir  string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "file download",
	Long:  `file download`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := file2.DownFile(hash)
		if err != nil {
			log.Print("error:", err)
			return
		}

		dst := filepath.Join(dir, "test.zip")
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
	downloadCmd.Flags().StringVarP(&dir, "dir", "d", "", "download file to dir")
	_ = downloadCmd.MarkFlagRequired("hash")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
