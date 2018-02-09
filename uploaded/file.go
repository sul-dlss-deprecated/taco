package uploaded

import (
	"mime/multipart"
)

// File represents a file and it's name & content-type
type File struct {
	Filename    string
	ContentType string
	Data        multipart.File
}

// NewFile creates a new uploade file instance
func NewFile(filename string, contentType string, file multipart.File) *File {
	return &File{Filename: filename, ContentType: contentType, Data: file}
}
