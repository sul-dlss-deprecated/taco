package storage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Storage the interface for file storage
type Storage interface {
	UploadFile(id string, file *datautils.File) (*string, error)
	CreateSignedURL(s3url string) (*string, error)
}

// NewS3Bucket creates a new storage adapter that uses S3 bucket storage to
// actually store the files
func NewS3BucketStorage(session *session.Session, bucketName string, endpoint string) Storage {
	forcePath := true // This is required for localstack
	s3Svc := s3.New(session, &aws.Config{
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: &forcePath,
	})
	return &S3BucketStorage{BucketName: bucketName, service: s3Svc}
}

// S3BucketStorage provides file storage using an S3 bucket
type S3BucketStorage struct {
	BucketName string
	service    *s3.S3
}

// CreateSignedURL returns a signed URL to the file specifed by key that
//   expires in 15 minutes
func (d *S3BucketStorage) CreateSignedURL(s3URL string) (*string, error) {
	u, err := url.Parse(s3URL)
	if err != nil {
		return nil, err
	}
	req, _ := d.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(u.Hostname()),
		Key:    aws.String(u.Path),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return nil, err
	}

	return &urlStr, nil
}

func (d *S3BucketStorage) uploader() *s3manager.Uploader {
	return s3manager.NewUploaderWithClient(d.service)
}

// UploadFile stores a file into S3
func (d *S3BucketStorage) UploadFile(key string, file *datautils.File) (*string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", file.Metadata.Filename)
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket:             &d.BucketName,
		Key:                &key,
		Body:               file.Data,
		ContentDisposition: &contentDisposition,
		ContentType:        &file.Metadata.ContentType,
	}

	_, err := d.uploader().Upload(upParams)

	if err != nil {
		return nil, err
	}

	return d.s3URI(key), nil
}

func (d *S3BucketStorage) s3URI(key string) *string {
	uri := fmt.Sprintf("s3:/%s/%s", d.BucketName, key)
	return &uri
}
