package taco

import (
	"github.com/spf13/viper"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(config *viper.Viper) (*Runtime, error) {
	repository, err := persistence.NewRepository(*config, db.NewConnection(config))
	if err != nil {
		return nil, err
	}
	return NewRuntimeForRepository(config, repository)
}

// NewRuntimeForRepository creates a new application level runtime that encapsulates the shared services for this application
func NewRuntimeForRepository(config *viper.Viper, repository persistence.Repository) (*Runtime, error) {
	return &Runtime{
		repository: repository,
		config:     config,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	repository persistence.Repository
	config     *viper.Viper
}

// Repository returns the metadata store
func (r *Runtime) Repository() persistence.Repository {
	return r.repository
}

// Config returns the viper config for this application
func (r *Runtime) Config() *viper.Viper {
	return r.config
}
