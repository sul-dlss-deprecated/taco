package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/uploaded"
)

// Storage the interface for file storage
type Storage interface {
	UploadFile(id string, file *uploaded.File) (*string, error)
}

// S3BucketStorage provides file storage using an S3 bucket
type S3BucketStorage struct {
	config   *config.Config
	uploader *s3manager.Uploader
}

// NewS3Bucket creates a new storage adapter that uses S3 bucket storage to
// actually store the files
func NewS3Bucket(config *config.Config, sess *session.Session) Storage {
	forcePath := true // This is required for localstack
	s3Svc := s3.New(sess, &aws.Config{
		Endpoint:         aws.String(config.S3Endpoint),
		S3ForcePathStyle: &forcePath,
	})
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	return &S3BucketStorage{config: config, uploader: uploader}
}

// UploadFile stores a file into S3
func (d *S3BucketStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", file.Filename)
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket:             &d.config.S3BucketName,
		Key:                &id,
		Body:               file.Data,
		ContentDisposition: &contentDisposition,
		ContentType:        &file.ContentType,
	}

	result, err := d.uploader.Upload(upParams)

	if err != nil {
		return nil, err
	}

	return &result.Location, nil
}
