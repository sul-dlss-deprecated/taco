package taco

import (
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime() (*Runtime, error) {
	repository, err := persistence.NewRepository(db.NewConnection())
	if err != nil {
		return nil, err
	}
	return NewRuntimeForRepository(repository)
}

// NewRuntimeForRepository creates a new application level runtime that encapsulates the shared services for this application
func NewRuntimeForRepository(repository persistence.Repository) (*Runtime, error) {
	return &Runtime{
		repository: repository,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	repository persistence.Repository
}

// Repository returns the metadata store
func (r *Runtime) Repository() persistence.Repository {
	return r.repository
}
