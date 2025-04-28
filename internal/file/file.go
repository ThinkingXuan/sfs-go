package file

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	progressbar "sfs-go/internal/pb"
	"strings"
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
func UploadFileToIPFS(fileBytes []byte) (string, error) {

	iPAddr := viper.GetString("ipfs_ip")
	sh := shell.NewShell(iPAddr + ":5001")

	var limit = int64(len(fileBytes))
	// start new bar
	bar := progressbar.NewProcessBarReader(limit)

	// create proxy reader
	barReader := bar.NewProxyReader(bytes.NewBuffer(fileBytes))

	//上传文件生成hash
	hash, err := sh.Add(barReader)
	if err != nil {
		log.Println("upload ipfs error:", err)
		return "", err
	}
	bar.Finish()
	return hash, nil
}

// WriteWithFile 使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteWithFile(filename, content string) {
	data := []byte(content)
	if ioutil.WriteFile(filename, data, 0644) == nil {
		return
	}
}

func ReadWithFile(filename string) string {
	if contents, err := ioutil.ReadFile(filename); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		return result
	}
	return ""
}

// ReadFileBytes read file from filepath ,return file bytes and error
func ReadFileBytes(filePath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}
