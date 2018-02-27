package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-openapi/runtime"
	"github.com/sul-dlss-labs/taco/config"
)

// Storage the interface for file storage
type Storage interface {
	UploadFile(id string, file runtime.File) (*string, error)
}

// S3BucketStorage provides file storage using an S3 bucket
type S3BucketStorage struct {
	config   *config.Config
	uploader *s3manager.Uploader
}

// NewS3Bucket creates a new S3 bucket storage
func NewS3Bucket(config *config.Config) Storage {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:   aws.String(config.S3Endpoint),
		DisableSSL: aws.Bool(config.S3DisableSSL),
	})
	if err != nil {
		panic(err)
	}
	// This is required for localstack:
	sess.Config.WithS3ForcePathStyle(true)
	uploader := s3manager.NewUploader(session.Must(sess, err))
	return &S3BucketStorage{config: config, uploader: uploader}
}

// UploadFile stores a file into S3
func (d *S3BucketStorage) UploadFile(id string, file runtime.File) (*string, error) {
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &d.config.S3BucketName,
		Key:    &id,
		Body:   file.Data,
	}

	result, err := d.uploader.Upload(upParams)

	if err != nil {
		return nil, err
	}

	return &result.Location, nil
}
