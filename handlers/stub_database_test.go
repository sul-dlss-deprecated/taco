package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/uploaded"
)

func handler(database db.Database, stream streaming.Stream, storage storage.Storage) http.Handler {
	if database == nil {
		database = NewMockDatabase(nil)
	}
	if stream == nil {
		stream = NewMockStream("")
	}
	if storage == nil {
		storage = NewMockStorage()
	}
	return BuildAPI(database, stream, storage).Serve(nil)
}

type MockDatabase struct {
	record           *db.Resource
	CreatedResources []db.Resource
}

func NewMockDatabase(record *db.Resource) db.Database {
	return &MockDatabase{CreatedResources: []db.Resource{}, record: record}
}

type MockStream struct {
	message string
}

func NewMockStream(message string) streaming.Stream {
	return &MockStream{message: message}
}

type MockStorage struct {
	CreatedFiles []*uploaded.File
}

func NewMockStorage() storage.Storage {
	return &MockStorage{CreatedFiles: []*uploaded.File{}}
}

func (d *MockDatabase) Insert(params db.Resource) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) Read(id string) (*db.Resource, error) {
	if d.record != nil {
		return d.record, nil
	}
	return nil, errors.New("not found")
}

func (d *MockDatabase) Update(params db.Resource) error {
	return nil
}

func (s *MockStream) SendMessage(message string) error {
	return nil
}

func (s *MockStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	s.CreatedFiles = append(s.CreatedFiles, file)
	path := "s3FileLocation"
	return &path, nil
}

type MockErrorDatabase struct {
}

func NewMockErrorDatabase() db.Database {
	return &MockErrorDatabase{}
}

func (d *MockErrorDatabase) Insert(params db.Resource) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) Update(params db.Resource) error {
	return nil
}

func (d *MockErrorDatabase) Read(id string) (*db.Resource, error) {
	return nil, errors.New("Broken")
}

type MockErrorStorage struct{}

func NewMockErrorStorage() storage.Storage {
	return &MockErrorStorage{}
}

func (d *MockErrorStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	return nil, errors.New("Broken")
}
