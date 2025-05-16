package uploader

import (
	"errors"
	"github.com/gomsr/atom-s3/configs3"
	"github.com/gomsr/atom-s3/util"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type local struct {
	*configs3.Local
}

func (c *local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := util.MD5V([]byte(strings.TrimSuffix(file.Filename, ext)))
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(c.StorePath, os.ModePerm)
	if mkdirErr != nil {
		return "", "", errors.New("function os.MkdirAll() failed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := c.StorePath + "/" + filename
	filepath := c.Path + "/" + filename

	f, openError := file.Open() // 读取文件
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		return "", "", errors.New("function os.Create() failed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		return "", "", errors.New("function io.Copy() failed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

func (c *local) DeleteFile(key string) error {
	p := c.StorePath + "/" + key
	if strings.Contains(p, c.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}

func NewLocalClient(c *configs3.Local) (*local, error) {
	return &local{
		Local: c,
	}, nil
}
