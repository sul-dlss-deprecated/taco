package runtime

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sul-dlss-labs/taco/aws_session"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
)

// Runtime represents the instantiated runtime objects
type Runtime struct {
	Database          db.Database
	Storage           storage.Storage
	IdentifierService identifier.Service
}

// NewRuntime instantiates all the runtime services
func NewRuntime(config *config.Config) *Runtime {
	awsSession := aws_session.Connect(config.AwsDisableSSL)
	database := &db.DynamodbDatabase{
		Connection: db.Connect(awsSession, config.DynamodbEndpoint),
		Table:      config.ResourceTableName,
	}

	storage := &storage.S3BucketStorage{
		Uploader:     connectToStorage(awsSession, config.S3Endpoint),
		S3BucketName: config.S3BucketName,
	}

	identifierService := identifier.NewService(config)
	return &Runtime{
		Database:          database,
		Storage:           storage,
		IdentifierService: identifierService,
	}
}

// NewS3Bucket creates a new storage adapter that uses S3 bucket storage to
// actually store the files
func connectToStorage(awsSession *session.Session, endpoint string) *s3manager.Uploader {
	forcePath := true // This is required for localstack
	s3Svc := s3.New(awsSession, &aws.Config{
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: &forcePath,
	})
	return s3manager.NewUploaderWithClient(s3Svc)
}
