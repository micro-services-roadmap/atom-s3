package uploader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomsr/atom-s3/configs3"
	"github.com/gomsr/atom-s3/util"
	"log"
	"mime/multipart"
	"path"
	"time"
)

type awsS3Client struct {
	*configs3.AwsS3
	*s3.Client
}

func (c *awsS3Client) CdnHost() string {
	return c.BaseURL
}

func (c *awsS3Client) GetFilename(key string, timestamp ...int64) (string, string) {
	var fileKey string
	if len(timestamp) > 0 {
		fileKey = fmt.Sprintf("%d_%s", time.Now().Unix(), key)
	} else {
		fileKey = fmt.Sprintf("%s", key)
	}

	var filename string
	if len(c.PathPrefix) == 0 {
		filename = fileKey
	} else {
		filename = c.PathPrefix + "/" + fileKey
	}

	return fileKey, filename
}

// Deprecated: use PreSigned
// @see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
func (c *awsS3Client) UploadFile(file *multipart.FileHeader) (string, string, error) {
	fileKey, filename := c.GetFilename(file.Filename, time.Now().Unix())
	reader, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(filename),
		Body:        reader,
		ContentType: aws.String(util.DetermineByFile(file)),
	})
	if err != nil {
		fmt.Println("function uploader.Upload() failed", err)
		log.Println("function uploader.Upload() failed", err)
		return "", "", err
	}

	return path.Join(c.BaseURL, filename), fileKey, nil
}

// Deprecated: use PreSigned
// @see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
func (c *awsS3Client) UploadFileV2(name string, bs []byte, ct string) (string, string, error) {
	fileKey, filename := c.GetFilename(name, time.Now().Unix())

	_, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(bs),
		ContentType: aws.String(ct),
	})
	if err != nil {
		fmt.Println("function uploader.Upload() failed", err)
		return "", "", err
	}

	return path.Join(c.BaseURL, filename), fileKey, nil

}

// DeleteFile
// @see https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-delete-bucket-item
func (c *awsS3Client) DeleteFile(key string) error {
	_, filename := c.GetFilename(key)
	bucket := c.Bucket

	// 删除对象
	_, err := c.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		log.Println("function svc.DeleteObject() failed", err)
		fmt.Println("function svc.DeleteObject() failed", err)
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}

	// 等待对象被删除
	objectWaiter := s3.NewObjectNotExistsWaiter(c.Client)
	err = objectWaiter.Wait(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}, 5*time.Minute) // 等待时间为5分钟
	if err != nil {
		log.Println("function objectWaiter.Wait() failed", err)
		fmt.Println("function objectWaiter.Wait() failed", err)
		return errors.New("function objectWaiter.Wait() failed, err:" + err.Error())
	}

	return nil
}

func (c *awsS3Client) PreSigned(key string) (*v4.PresignedHTTPRequest, error) {
	_, filename := c.GetFilename(key, time.Now().Unix())

	in := &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(filename),
	}

	return s3.NewPresignClient(c.Client).PresignPutObject(context.TODO(), in)
}

// NewS3Client Create S3 session
func NewS3Client(c *configs3.AwsS3) (*awsS3Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: c.Endpoint}, nil
		})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.SecretID, c.SecretKey, "")),
		config.WithRegion(c.Region),
	)
	// S3ForcePathStyle: aws.Bool(v.S3ForcePathStyle),
	// DisableSSL:       aws.Bool(v.DisableSSL),
	if err != nil {
		return nil, err
	}

	return &awsS3Client{
		AwsS3:  c,
		Client: s3.NewFromConfig(cfg),
	}, nil
}
