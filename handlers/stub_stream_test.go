package handlers

import (
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/streaming"
)

type MockStream struct {
	Messages []*streaming.Message
}

func NewMockStream() streaming.Stream {
	return &MockStream{Messages: []*streaming.Message{}}
}

func (s *MockStream) Send(action string, resource *datautils.Resource) error {
	message := streaming.NewMessage(action, resource)
	s.Messages = append(s.Messages, message)
	return nil
}
