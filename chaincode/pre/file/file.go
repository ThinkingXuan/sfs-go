package file

import (
	"io/ioutil"
	"strings"
)


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
