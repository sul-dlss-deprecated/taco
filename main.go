package main

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/server"
	"github.com/sul-dlss-labs/taco/sessionbuilder"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
)

// Taco encapsulates the shared services for this application
type Taco struct {
	repository  persistence.Repository
	stream      streaming.Stream
	fileStorage storage.Storage
	config      *config.Config
}

func main() {

	rt := Taco{}

	awsSession := sessionbuilder.NewAwsSession(config)
	dbConn := db.NewConnection(config, awsSession)
	repository := persistence.NewDynamoRepository(config, dbConn)
	storage := storage.NewS3Bucket(config, awsSession)
	stream := streaming.NewKinesisStream(config, awsSession)

	if err != nil {
		log.Fatalln(err)
	}

	server := createServer()
	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func createServer() *restapi.Server {
	api := server.BuildAPI()
	server := restapi.NewServer(api)
	server.SetHandler(server.BuildHandler(api))
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = rt.Config().Port
	return server
}

// NewRuntime creates a new application level runtime that
// encapsulates the shared services for this application
func NewRuntime(config *config.Config) (*Runtime, error) {

	return NewRuntimeWithServices(config, repository, storage, stream)
}

// NewRuntimeWithServices creates a new application level runtime that encapsulates the shared services for this application
func NewRuntimeWithServices(config *config.Config, repository persistence.Repository, fileStorage storage.Storage, stream streaming.Stream) (*Runtime, error) {
	return &Runtime{
		repository:  repository,
		stream:      stream,
		config:      config,
		fileStorage: fileStorage,
	}, nil
}

// Repository returns the metadata store
func (r *Runtime) Repository() persistence.Repository {
	return r.repository
}

// Stream returns the kinesis stream
func (r *Runtime) Stream() streaming.Stream {
	return r.stream
}

// Config returns the config for this application
func (r *Runtime) Config() *config.Config {
	return r.config
}

// FileStorage returns the file store
func (r *Runtime) FileStorage() storage.Storage {
	return r.fileStorage
}
