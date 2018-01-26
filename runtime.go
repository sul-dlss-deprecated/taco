package taco

import (
	"github.com/spf13/viper"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
)

// NewRuntime creates a new application level runtime that
// encapsulates the shared services for this application
func NewRuntime(config *viper.Viper) (*Runtime, error) {
	repository, err := persistence.NewRepository(*config, db.NewConnection(config))
	if err != nil {
		return nil, err
	}

	stream, err := streaming.NewKinesisStream(*config)
	if err != nil {
		return nil, err
	}
	return NewRuntimeWithServices(config, repository, stream)
}

// NewRuntimeWithServices creates a new application level runtime that encapsulates the shared services for this application
func NewRuntimeWithServices(config *viper.Viper, repository persistence.Repository, stream streaming.Stream) (*Runtime, error) {
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
	config     *viper.Viper
}

// Repository returns the metadata store
func (r *Runtime) Repository() persistence.Repository {
	return r.repository
}

// Stream returns the kinesis stream
func (r *Runtime) Stream() streaming.Stream {
	return r.stream
}

// Config returns the viper config for this application
func (r *Runtime) Config() *viper.Viper {
	return r.config
}
