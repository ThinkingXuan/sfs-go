package encrypt

import (
	"fmt"
	"testing"
)

func TestReadFileAttribute(t *testing.T) {
	filePath := "D:\\document\\论文\\我的论文-安全文件共享系统\\实验资料\\文件\\文件_1000M.zip"
	fmt.Println(GenerateEncryptionMethod(filePath, confidential))
}
