package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Storage the interface for file storage
type Storage interface {
	UploadFile(id string, file *datautils.File) (*string, error)
}

// S3BucketStorage provides file storage using an S3 bucket
type S3BucketStorage struct {
	Uploader     *s3manager.Uploader
	S3BucketName string
}

// UploadFile stores a file into S3
func (d *S3BucketStorage) UploadFile(id string, file *datautils.File) (*string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", file.Metadata.Filename)
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket:             &d.S3BucketName,
		Key:                &id,
		Body:               file.Data,
		ContentDisposition: &contentDisposition,
		ContentType:        &file.Metadata.ContentType,
	}

	result, err := d.Uploader.Upload(upParams)

	if err != nil {
		return nil, err
	}

	return &result.Location, nil
}
