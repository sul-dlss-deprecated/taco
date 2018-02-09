package handlers

import (
	"errors"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/storage"
)

type MockStorage struct {
	CreatedFiles []*datautils.File
}

func NewMockStorage() storage.Storage {
	return &MockStorage{CreatedFiles: []*datautils.File{}}
}

func (s *MockStorage) UploadFile(id string, file *datautils.File) (*string, error) {
	s.CreatedFiles = append(s.CreatedFiles, file)
	path := "s3FileLocation"
	return &path, nil
}

func (s *MockStorage) CreateSignedURL(s3URI string) (*string, error) {
	path := "https://example.com/file-123"
	return &path, nil
}

type MockErrorStorage struct{}

func NewMockErrorStorage() storage.Storage {
	return &MockErrorStorage{}
}

func (f *MockErrorStorage) UploadFile(id string, file *datautils.File) (*string, error) {
	return nil, errors.New("broken")
}

func (f *MockErrorStorage) CreateSignedURL(s3URI string) (*string, error) {
	return nil, errors.New("Broken")
}
