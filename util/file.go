package util

import (
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"
	"net/http"
)

func DetermineByFile(fileHeader *multipart.FileHeader) string {
	file, err := fileHeader.Open()
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	// 读取前 512 字节用于检测 MIME 类型
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "unknown"
	}

	// 使用 http.DetectContentType 自动检测
	return http.DetectContentType(buf)
}

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
