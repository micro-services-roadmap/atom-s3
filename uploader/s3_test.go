package uploader

import (
	"fmt"
	"github.com/gomsr/atom-s3/configs3"
	"testing"
)

func TestR2S3Client(t *testing.T) {
	s3, err := NewS3Client(&configs3.AwsS3{
		S3: configs3.S3{
			Endpoint:  "https://251167b1b2a63fb935d4a20e2c1966af.r2.cloudflarestorage.com",
			Bucket:    "na-test",
			Region:    "auto",
			SecretID:  "xxx",
			SecretKey: "xxxx",
			BaseURL:   "https://s3.hubby.top",
		},
	})
	if err != nil {
		panic(err)
	}

	preSigned, err := s3.PreSigned("test.png")
	if err != nil {
		return
	}
	fmt.Println(preSigned)
}

func Test_awsS3Client_UploadFile(t *testing.T) {

}

func TestOssS3Client(t *testing.T) {
	s3, err := NewS3Client(&configs3.AwsS3{
		S3: configs3.S3{
			Endpoint:  "https://oss-cn-hangzhou.aliyuncs.com",
			Bucket:    "project-ec",
			Region:    "hangzhou",
			SecretID:  "xx",
			SecretKey: "xx",
			BaseURL:   "https://s3.hubby.top",
		},
	})
	if err != nil {
		panic(err)
	}

	preSigned, err := s3.PreSigned("test.png")
	if err != nil {
		return
	}
	fmt.Println(preSigned)
}

func TestCosS3Client(t *testing.T) {
	s3, err := NewS3Client(&configs3.AwsS3{
		S3: configs3.S3{
			Endpoint:  "https://cos.ap-shanghai.myqcloud.com",
			Bucket:    "project-ec-1300043990",
			Region:    "ap-shanghai",
			SecretID:  "xx",
			SecretKey: "xx",
			BaseURL:   "https://s3.hubby.top",
		},
	})
	if err != nil {
		panic(err)
	}

	preSigned, err := s3.PreSigned("test.png")
	if err != nil {
		panic(err)
	}
	fmt.Println(preSigned)
}
