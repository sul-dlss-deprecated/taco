package uploaded

import (
	"mime/multipart"
)

// FileMetadata represents the metadata that comes with an uploaded file
type FileMetadata struct {
	Filename    string
	ContentType string
}

// File represents a file and it's name & content-type
type File struct {
	Metadata FileMetadata
	Data     multipart.File
}

// NewFile creates a new uploaded file instance
func NewFile(metadata FileMetadata, data multipart.File) *File {
	return &File{Metadata: metadata, Data: data}
}
