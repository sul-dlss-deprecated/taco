package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/streaming"
)

// NewDeleteResource -- Accepts requests to remove a resource and pushes them to Kinesis.
func NewDeleteResource(repository db.Database, stream streaming.Stream) operations.DeleteResourceHandler {
	return &deleteResourceEntry{
		repository: repository,
		stream:     stream,
	}
}

type deleteResourceEntry struct {
	stream     streaming.Stream
	repository db.Database
	// s3
}

// Handle the delete entry request
// TODO: Delete from S3
func (d *deleteResourceEntry) Handle(params operations.DeleteResourceParams) middleware.Responder {
	if err := d.repository.DeleteAllVersions(params.ID); err != nil {
		panic(err)
	}

	if err := d.addToStream(params.ID); err != nil {
		panic(err)
	}
	return operations.NewDeleteResourceNoContent()
}

func (d *deleteResourceEntry) addToStream(id string) error {
	message := streaming.Message{ID: id, Action: "delete"}
	return d.stream.SendMessage(message)
}
