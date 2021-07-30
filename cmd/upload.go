package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"sfs-go/internal/file"
)

var (
	path string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Print("filepath: ", path)
		hash, err := file.UploadFileToIPFS(path)
		if err != nil {
			log.Printf("file[%s] upload failure", path)
		}
		log.Println("file upload success")
		log.Println(hash)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&path, "path", "p", "", "file upload path")
	_ = uploadCmd.MarkFlagRequired("path")
}
