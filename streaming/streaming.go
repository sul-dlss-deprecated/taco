package streaming

import "github.com/sul-dlss-labs/taco/datautils"

// Stream the interface for streaming pipeline
type Stream interface {
	Send(string, *datautils.Resource) error
}
