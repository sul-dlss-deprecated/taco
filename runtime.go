package taco

import (
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/sessionbuilder"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
)

// NewRuntime creates a new application level runtime that
// encapsulates the shared services for this application
func NewRuntime(config *config.Config) (*Runtime, error) {
	return &Runtime{config: config}, nil
}

// NewDefaultRuntime creates a new application level runtime with all
// necessary services pre-configured
func NewDefaultRuntime(config *config.Config) (*Runtime, error) {
	rt, _ := NewRuntime(config)
	awsSession := sessionbuilder.NewAwsSession(config)
	dbConn := db.NewConnection(config, awsSession)

	return rt.
		WithRepository(persistence.NewDynamoRepository(config, dbConn)).
		WithStorage(storage.NewS3Bucket(config, awsSession)).
		WithStreaming(streaming.NewKinesisStream(config, awsSession)).
		WithIdentifierService(identifier.NewService(config)), nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	repository  persistence.Repository
	stream      streaming.Stream
	fileStorage storage.Storage
	ident       *identifier.TypeAwareService
	config      *config.Config
}

// WithRepository sets the passed in repository on the runtime.
func (r *Runtime) WithRepository(repo persistence.Repository) *Runtime {
	r.repository = repo
	return r
}

// WithStorage sets the passed in storage on the runtime.
func (r *Runtime) WithStorage(store storage.Storage) *Runtime {
	r.fileStorage = store
	return r
}

// WithStreaming sets the passed in streaming on the runtime.
func (r *Runtime) WithStreaming(stream streaming.Stream) *Runtime {
	r.stream = stream
	return r
}

// WithIdentifierService sets the passed in streaming on the runtime.
func (r *Runtime) WithIdentifierService(idService *identifier.TypeAwareService) *Runtime {
	r.ident = idService
	return r
}

// Repository returns the metadata store
func (r *Runtime) Repository() persistence.Repository {
	if r.repository == nil {
		panic("Repository not initialized")
	}
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

// Identifier returns the identifier service
func (r *Runtime) Identifier() *identifier.TypeAwareService {
	if r.ident == nil {
		panic("Identifier service not initialized")
	}
	return r.ident
}
