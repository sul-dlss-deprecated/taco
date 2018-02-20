package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/uploaded"
)

func mockRepo(record *persistence.Resource) persistence.Repository {
	return &fakeRepository{record: record, CreatedResources: []persistence.Resource{}}
}

type fakeRepository struct {
	record           *persistence.Resource
	CreatedResources []persistence.Resource
}

func (f *fakeRepository) GetByID(id string) (*persistence.Resource, error) {

	if f.record != nil {
		return f.record, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepository) CreateItem(resource *persistence.Resource) error {
	f.CreatedResources = append(f.CreatedResources, *resource)
	return nil
}

func (f *fakeRepository) UpdateItem(resource *persistence.Resource) error {
	return nil
}

func mockStream() streaming.Stream {
	return &fakeStream{}
}

type fakeStream struct {
}

func (d fakeStream) SendMessage(message string) error { return nil }

func mockStorage() storage.Storage {
	return &fakeStorage{CreatedFiles: []*uploaded.File{}}
}

type fakeStorage struct {
	CreatedFiles []*uploaded.File
}

func (f *fakeStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	f.CreatedFiles = append(f.CreatedFiles, file)
	path := "s3FileLocation"
	return &path, nil
}

func setupFakeRuntime() *TestEnv {
	return &TestEnv{}
}

type TestEnv struct {
	storage storage.Storage
	repo    persistence.Repository
}

func (d *TestEnv) WithRepository(repo persistence.Repository) *TestEnv {
	d.repo = repo
	return d
}

func (d *TestEnv) WithStorage(storage storage.Storage) *TestEnv {
	d.storage = storage
	return d
}

func (d *TestEnv) Handler() http.Handler {
	if d.repo == nil {
		d.repo = mockRepo(nil)
	}

	if d.storage == nil {
		d.storage = &fakeStorage{}
	}

	rt, _ := taco.NewRuntime(nil)
	rt = rt.WithRepository(d.repo).WithStorage(d.storage).
		WithStreaming(mockStream()).
		WithIdentifierService(identifier.NewUUIDService())
	return BuildAPI(rt).Serve(nil)
}

func mockErrorRepo() persistence.Repository {
	return &fakeErroringRepository{}
}

type fakeErroringRepository struct{}

func (f *fakeErroringRepository) GetByID(id string) (*persistence.Resource, error) {
	return nil, errors.New("broken")
}

func (f *fakeErroringRepository) CreateItem(resource *persistence.Resource) error {
	return errors.New("broken")
}

func (f *fakeErroringRepository) UpdateItem(resource *persistence.Resource) error {
	return errors.New("broken")
}

func mockErrorStorage() storage.Storage {
	return &fakeErroringStorage{}
}

type fakeErroringStorage struct{}

func (f *fakeErroringStorage) UploadFile(id string, file *uploaded.File) (*string, error) {
	return nil, errors.New("broken")
}
