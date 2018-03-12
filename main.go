package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/handlers"
	"github.com/sul-dlss-labs/taco/middleware"
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
	tacoServer.awsSession = newAwsSession()
	database := &db.DynamodbDatabase{
		Connection: connectToDatabase(),
		Table:      tacoServer.config.ResourceTableName,
	}
	stream := &streaming.KinesisStream{
		Connection: connectToStream(),
		StreamName: tacoServer.config.DepositStreamName,
	}

	tacoServer.server = createServer(database, stream)

	//	storage := storage.NewS3Bucket(config, awsSession)
	//	stream := streaming.NewKinesisStream(config, awsSession)

	// serve API
	if err := tacoServer.server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func connectToDatabase() *dynamodb.DynamoDB {
	dynamoConfig := &aws.Config{Endpoint: aws.String(tacoServer.config.DynamodbEndpoint)}
	return dynamodb.New(tacoServer.awsSession, dynamoConfig)
}

func connectToStream() *kinesis.Kinesis {
	return kinesis.New(tacoServer.awsSession, &aws.Config{Endpoint: aws.String(tacoServer.config.KinesisEndpoint)})
}

func createServer(database db.Database, stream streaming.Stream) *restapi.Server {
	api := handlers.BuildAPI(database, stream)
	server := restapi.NewServer(api)
	server.SetHandler(BuildHandler(api))
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = tacoServer.config.Port
	return server
}

func newAwsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(tacoServer.config.AwsDisableSSL),
	}))
}

// BuildHandler sets up the middleware that wraps the API
func BuildHandler(api *operations.TacoAPI) http.Handler {
	return alice.New(
		middleware.NewHoneyBadgerMW(),
		middleware.NewRecoveryMW(),
		middleware.NewRequestLoggerMW(),
	).Then(api.Serve(nil))
}
