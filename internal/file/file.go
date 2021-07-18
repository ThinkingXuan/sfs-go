package file

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
)

// DownFile 连接ipfs进行下载操作
func DownFile(fileHash string) ([]byte, error) {
	iPAddr := ""
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
	iPAddr := ""
	sh := shell.NewShell(iPAddr + ":5001")

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//上传文件生成hash
	hash, err := sh.Add(bytes.NewBuffer(file))
	if err != nil {
		fmt.Println("上传ipfs时错误：", err)
		return "", err
	}
	return hash, nil
}
