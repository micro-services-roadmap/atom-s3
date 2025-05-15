package uploader

import (
	"bytes"
	"fmt"
	"github.com/micro-services-roadmap/atom-s3/configs3"
	"io"
	"mime/multipart"
	"os"
	"testing"
)

func TestNewLocalClient(t *testing.T) {
	c, err := NewLocalClient(&configs3.Local{
		Path:      "uploads/file",
		StorePath: "uploads/file",
	})
	if err != nil {
		panic(err)
	}
	testFile := GetFile("aaa.png")

	filename, path, err := c.UploadFile(testFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(filename, path)
}

func GetFile(filePath string) *multipart.FileHeader {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// 创建一个缓冲区和 multipart writer
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// 创建 form 文件字段
	part, err := writer.CreateFormFile("file", fileInfo.Name())
	if err != nil {
		panic(err)
	}

	// 将文件内容拷贝到 part
	if _, err := io.Copy(part, file); err != nil {
		panic(err)
	}

	// 关闭 writer，生成完整的 multipart 数据
	writer.Close()

	// 解析 multipart 数据，生成 *multipart.FileHeader
	requestReader := bytes.NewReader(buffer.Bytes())
	multipartReader := multipart.NewReader(requestReader, writer.Boundary())
	form, err := multipartReader.ReadForm(10 << 20) // 最大 10MB
	if err != nil {
		panic(err)
	}

	// 获取文件头
	return form.File["file"][0]
}
