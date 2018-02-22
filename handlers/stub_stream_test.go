package handlers

import "github.com/sul-dlss-labs/taco/streaming"

type MockStream struct {
	message string
}

func NewMockStream(message string) streaming.Stream {
	return &MockStream{message: message}
}

func (d *MockStream) SendMessage(message string) error {
	return nil
}
