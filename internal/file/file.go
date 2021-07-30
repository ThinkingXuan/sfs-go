package file

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	progressbar "sfs-go/internal/pb"
)

// DownFile 连接ipfs进行下载操作
func DownFile(fileHash string) ([]byte, error) {
	iPAddr := viper.GetString("ipfs_ip")
	sh := shell.NewShell(iPAddr + ":5001")
	fmt.Println("下载文件哈希值：" + fileHash)
	//从ipfs下载数据
	read, err := sh.Cat(fileHash)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(read)

	return body, nil
}

// UploadFileToIPFS 上传文件到IPFS
func UploadFileToIPFS(filename string) (string, error) {
	iPAddr := viper.GetString("ipfs_ip")
	sh := shell.NewShell(iPAddr + ":5001")

	log.Println("start file read!")
	file, err := ioutil.ReadFile(filename)
	log.Println("file read finish, start upload!")

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var limit = int64(len(file))
	// start new bar
	bar := progressbar.NewProcessBarReader(limit)

	// create proxy reader
	barReader := bar.NewProxyReader(bytes.NewBuffer(file))

	//上传文件生成hash
	hash, err := sh.Add(barReader)
	if err != nil {
		log.Println("upload ipfs error:", err)
		return "", err
	}
	bar.Finish()
	return hash, nil
}
