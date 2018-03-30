package handlers

import "github.com/sul-dlss-labs/taco/streaming"

type MockStream struct {
	Messages []string
}

func NewMockStream() streaming.Stream {
	return &MockStream{Messages: []string{}}
}

func (s *MockStream) SendMessage(message string) error {
	s.Messages = append(s.Messages, message)
	return nil
}
