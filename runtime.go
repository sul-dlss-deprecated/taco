package taco

import (
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
)

// NewRuntime creates a new application level runtime that
// encapsulates the shared services for this application
func NewRuntime(config *config.Config) (*Runtime, error) {
	repository, err := persistence.NewRepository(config, db.NewConnection(config))
	if err != nil {
		return nil, err
	}

	stream, err := streaming.NewKinesisStream(config)
	if err != nil {
		return nil, err
	}
	return NewRuntimeWithServices(config, repository, stream)
}

// NewRuntimeWithServices creates a new application level runtime that encapsulates the shared services for this application
func NewRuntimeWithServices(config *config.Config, repository persistence.Repository, stream streaming.Stream) (*Runtime, error) {
	return &Runtime{
		repository: repository,
		stream:     stream,
		config:     config,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	repository persistence.Repository
	stream     streaming.Stream
	config     *config.Config
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
