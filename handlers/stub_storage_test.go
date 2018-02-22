package handlers

import (
	"errors"

	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/uploaded"
)

type MockStorage struct {
	CreatedFiles []*uploaded.File
}

func (s *MockStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	s.CreatedFiles = append(s.CreatedFiles, file)
	path := "s3FileLocation"
	return &path, nil
}

func NewMockStorage() storage.Storage {
	return &MockStorage{CreatedFiles: []*uploaded.File{}}
}

func NewMockErrorStorage() storage.Storage {
	return &fakeErroringStorage{}
}

type fakeErroringStorage struct{}

func (f *fakeErroringStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	return nil, errors.New("broken")
}
