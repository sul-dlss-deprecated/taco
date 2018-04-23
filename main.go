package main

//go:generate go run generate/includemaps.go
import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/aws_session"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/handlers"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/middleware"
	"github.com/sul-dlss-labs/taco/storage"
)

func main() {
	// Initialize our global struct
	config := config.NewConfig()
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
	server := createServer(database, storage, identifierService, config.Port)
	defer server.Shutdown()

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
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

func createServer(database db.Database, storage storage.Storage, identifierService identifier.Service, port int) *restapi.Server {
	api := handlers.BuildAPI(database, storage, identifierService)
	server := restapi.NewServer(api)
	server.SetHandler(BuildHandler(api))

	// set the port this service will be run on
	server.Port = port
	return server
}

// BuildHandler sets up the middleware that wraps the API
func BuildHandler(api *operations.TacoAPI) http.Handler {
	return alice.New(
		middleware.NewHoneyBadgerMW(),
		middleware.NewRecoveryMW(),
		middleware.NewRequestLoggerMW(),
	).Then(api.Serve(nil))
}
