package atoms3

import (
	"github.com/micro-services-roadmap/atom-s3/configs3"
	"testing"
)

func TestNewS3(t *testing.T) {
	s3, err := NewS3(&configs3.AwsS3{})
	if err != nil {
		panic(err)
	}

	_, err = s3.PreSigned("")
}

func TestNewLocal(t *testing.T) {
	s3, err := NewLocal(&configs3.Local{})
	if err != nil {
		panic(err)
	}

	_ = s3.DeleteFile("")
}
