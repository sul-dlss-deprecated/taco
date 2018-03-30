package streaming

import (
	"github.com/sul-dlss-labs/taco/datautils"
)

// Message represents the message we send to kinesis
type Message struct {
	Action       string
	ResourceType string
	ID           string
}

// NewMessage builds a message given an action and resource
func NewMessage(action string, resource *datautils.Resource) *Message {
	return &Message{
		ID:           resource.ID(),
		ResourceType: resource.Type(),
		Action:       "deposit",
	}
}
