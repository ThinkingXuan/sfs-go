package cmd

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"sfs-go/internal/encrypt"
	ecc2 "sfs-go/internal/encrypt/ecc"
	"sfs-go/internal/fabric/fabservice"
	"sfs-go/internal/fabric/sdkInit"
	"sfs-go/internal/file"
	"sfs-go/internal/tools"
	"time"
)

var (
	filePath  string
	startTime = time.Now()
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
		// read file byte
		log.Println("start file read!")
		fileBytes, err := file.ReadFileBytes(filePath)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("读文件：", time.Since(startTime))
		// file mixed encryption
		fileCipherText, eccCipherText, err := filehHybridEncryption(fileBytes)
		if err != nil {
			log.Println(err)
			return
		}

		// upload file to ipfs, get file hash
		hash, err := file.UploadFileToIPFS(fileCipherText)
		if err != nil {
			log.Printf("file[%s] upload failure", filePath)
			return
		}
		fmt.Println("上传IPFS：", time.Since(startTime))

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

		// ecc cipher text bytes to string
		fileEncryptCipher := tools.ByteToString(eccCipherText)

		// insert fabric
		err = insertFileToFabric(fileCC, fileEncryptCipher)
		if err != nil {
			log.Println("insert fabric failure,", err)
			return
		}

		fmt.Println("上传Fabric：", time.Since(startTime))

		log.Println("file upload success!!!!")

	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&filePath, "path", "p", "", "file upload path")
	//uploadCmd.Flags().Bool("pri", false, "if set private, file no share")
	_ = uploadCmd.MarkFlagRequired("path")
}

func insertFileToFabric(fileCC fabservice.File, fileEncryptCipher string) error {
	service := sdkInit.GetInstance().InitFabric()
	_, err := service.InsertFile(fileCC)
	if err != nil {
		log.Println(err)
		return err
	}
	address := file.ReadWithFile("config/my.address")
	_, err = service.InsertAddressFile(address, fileCC.FileID, fileEncryptCipher)
	if err != nil {
		log.Println("failure")
		return err
	}

	return nil

	//newFIleBytes, _ := service.QueryFile(fileCC.FileID)
	//fmt.Println(newFIleBytes)
	//var newFile fabservice.File
	//_ = json.Unmarshal(newFIleBytes, &newFile)
	//fmt.Println(newFile.FileEncryptCipher)
	//
	//fmt.Println("is", newFile.FileEncryptCipher == fileCC.FileEncryptCipher)
}

// filehHybridEncryption
func filehHybridEncryption(srcFileBytes []byte) ([]byte, []byte, error) {
	// generate Aes encryption key
	fileID, _ := uuid.GenerateUUID()
	aesEncryptKeyBytes := generateEncryptKey(fileID)
	//create a new AES struct
	aes := encrypt.NewAes()
	//aes := encrypt.NewAesCipher128(aesEncryptKeyBytes[:16], aesEncryptKeyBytes[:16])
	// Aes encrypt
	fileCipherText, err := aes.AESEncrypt(srcFileBytes, aesEncryptKeyBytes)
	if err != nil {
		log.Println("aes encrypt err:", err)
		return nil, nil, err
	}

	fmt.Println("AES加密：", time.Since(startTime))

	// create a new ecc struct
	ecc := ecc2.NewECC("config/")
	eccCipherText, err := ecc.EccEncrypt(aesEncryptKeyBytes)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	fmt.Println("ECC加密：", time.Since(startTime))
	return fileCipherText, eccCipherText, nil
}

func generateEncryptKey(fileID string) []byte {
	return []byte(fileID)[:32]
}
