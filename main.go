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
	"github.com/sul-dlss-labs/taco/middleware"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
)

type Taco struct {
	config     *config.Config
	server     *restapi.Server
	awsSession *session.Session
}

var tacoServer Taco

func main() {

	// Initialize our global struct
	tacoServer.config = config.NewConfig()
	tacoServer.awsSession = aws_session.Connect(tacoServer.config.AwsDisableSSL)
	database := &db.DynamodbDatabase{
		Connection: db.Connect(tacoServer.awsSession, tacoServer.config.DynamodbEndpoint),
		Table:      tacoServer.config.ResourceTableName,
	}
	stream := &streaming.KinesisStream{
		Connection: streaming.Connect(tacoServer.awsSession, tacoServer.config.KinesisEndpoint),
		StreamName: tacoServer.config.DepositStreamName,
	}
	storage := &storage.S3BucketStorage{
		Uploader:     connectToStorage(),
		S3BucketName: tacoServer.config.S3BucketName,
	}
	tacoServer.server = createServer(database, stream, storage)

	// serve API
	if err := tacoServer.server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

// NewS3Bucket creates a new storage adapter that uses S3 bucket storage to
// actually store the files
func connectToStorage() *s3manager.Uploader {
	forcePath := true // This is required for localstack
	s3Svc := s3.New(tacoServer.awsSession, &aws.Config{
		Endpoint:         aws.String(tacoServer.config.S3Endpoint),
		S3ForcePathStyle: &forcePath,
	})
	return s3manager.NewUploaderWithClient(s3Svc)
}

func createServer(database db.Database, stream streaming.Stream, storage storage.Storage) *restapi.Server {
	api := handlers.BuildAPI(database, stream, storage)
	server := restapi.NewServer(api)
	server.SetHandler(BuildHandler(api))
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = tacoServer.config.Port
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
