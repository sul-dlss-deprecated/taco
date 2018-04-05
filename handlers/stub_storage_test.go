package handlers

import (
	"errors"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/storage"
)

type MockStorage struct {
	CreatedFiles []*datautils.File
}

func (s *MockStorage) UploadFile(id string, file *datautils.File) (*string, error) {
	s.CreatedFiles = append(s.CreatedFiles, file)
	path := "s3FileLocation"
	return &path, nil
}

func NewMockStorage() storage.Storage {
	return &MockStorage{CreatedFiles: []*datautils.File{}}
}

func NewMockErrorStorage() storage.Storage {
	return &fakeErroringStorage{}
}

type fakeErroringStorage struct{}

func (f *fakeErroringStorage) UploadFile(id string, file *datautils.File) (*string, error) {
	return nil, errors.New("broken")
}
