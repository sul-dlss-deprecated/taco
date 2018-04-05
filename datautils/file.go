package datautils

import (
	"mime/multipart"
)

// FileMetadata represents the metadata that comes with an uploaded file and is
// required for pushing to the storage adapter.
type FileMetadata struct {
	Filename    string
	ContentType string
}

// File represents a file and it's name & content-type
type File struct {
	Metadata FileMetadata
	Data     multipart.File
}

// NewFile creates a new uploaded file instance that has
// filename, mime-type and the binary data. This is the structure required by
// the file storage adapter.
func NewFile(metadata FileMetadata, data multipart.File) *File {
	return &File{Metadata: metadata, Data: data}
}
