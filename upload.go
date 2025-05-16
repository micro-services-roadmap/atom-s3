package atoms3

import (
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/gomsr/atom-s3/configs3"
	"github.com/gomsr/atom-s3/uploader"
	"mime/multipart"
)

type Uploader interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

// S3 对象存储接口
type S3 interface {
	Uploader
	GetFilename(key string, timestamp ...int64) (string, string)

	UploadFileV2(name string, bs []byte, ct string) (string, string, error)
	//UploadFile(file *multipart.FileHeader) (string, string, error)
	//DeleteFile(key string) error

	PreSigned(key string) (*v4.PresignedHTTPRequest, error)

	CdnHost() string
}

// NewS3 OSS的实例化方法
func NewS3(c *configs3.AwsS3) (S3, error) {
	return uploader.NewS3Client(c)
}

func NewLocal(c *configs3.Local) (Uploader, error) {
	return uploader.NewLocalClient(c)
}
