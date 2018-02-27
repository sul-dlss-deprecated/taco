package handlers

import "github.com/sul-dlss-labs/taco/streaming"

type MockStream struct {
	Messages []streaming.Message
}

func NewMockStream() streaming.Stream {
	return &MockStream{Messages: []streaming.Message{}}
}

func (s *MockStream) SendMessage(message streaming.Message) error {
	s.Messages = append(s.Messages, message)
	return nil
}
