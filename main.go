package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/handlers"
	"github.com/sul-dlss-labs/taco/middleware"
)

type Taco struct {
	config     *config.Config
	server     *restapi.Server
	awsSession *session.Session
	api        *operations.TacoAPI
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
	tacoServer.api = handlers.BuildAPI(database)
	tacoServer.server = createServer()

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

func createServer() *restapi.Server {
	server := restapi.NewServer(tacoServer.api)
	server.SetHandler(addMiddleware())
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = tacoServer.config.Port
	return server
}

func addMiddleware() http.Handler {
	return alice.New(
		middleware.NewHoneyBadgerMW(),
		middleware.NewRecoveryMW(),
		middleware.NewRequestLoggerMW(),
	).Then(tacoServer.api.Serve(nil))
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
