package models

import "fmt"

// Resource -- The metadata object
type Resource struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// GetByID -- given an identifier, find the resource
func (h Resource) GetByID(id string) (*Resource, error) {
	fmt.Printf("%s", id)
	return new(Resource), nil
}
