package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

const filesetID = "99999"

var filePath = fmt.Sprintf("/v1/resource/%s/file", filesetID)

const contentType = "multipart/form-data; boundary=------------------------a31e2ddd4b2c0d92"
const body = `--------------------------a31e2ddd4b2c0d92
Content-Disposition: form-data; name="upload"; filename="foo.txt"
Content-Type: text/plain

Hello

--------------------------a31e2ddd4b2c0d92--`

func TestCreateFileHappyPath(t *testing.T) {
	r := gofight.New()
	storage := NewMockStorage()
	fileSet := datautils.NewResource(nil).
		WithType(datautils.FilesetType).
		WithID(filesetID)
	repo := NewMockDatabase(fileSet, nil)

	r.POST(filePath).
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
			"Content-Type": contentType,
		}).
		SetBody(body).
		Run(handler(repo, storage, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(storage.(*MockStorage).CreatedFiles))
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				fileResource := repo.(*MockDatabase).CreatedResources[0]
				fileName := fileResource.Identification().GetS("filename")
				assert.Equal(t, fileName, "foo.txt")
				assert.Equal(t, fileResource.Label(), "foo.txt")
				assert.Equal(t, "text/plain", fileResource.MimeType())
				assert.Equal(t, "s3FileLocation", fileResource.FileLocation())
				assert.Equal(t, fileSet.ID(), fileResource.Structural().GetS("isContainedBy"))
			})
}

func TestCreateFileWrongContentType(t *testing.T) {
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
			"Content-Type": "application/xml",
		}).
		SetBody(``).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
			})
}

// TODO: https://github.com/go-swagger/go-swagger/issues/1400
// func TestCreateFileMissingFile(t *testing.T) {
// 	r := gofight.New()
// 	r.POST(filePath).
// 		SetHeader(gofight.H{
// 			"Content-Type": contentType,
// 		}).
// 		SetBody(``).
// 		Run(setupFakeRuntime().Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
// 			})
// }

func TestCreateFileNoApiKey(t *testing.T) {
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"Content-Type": contentType,
		}).
		SetBody(body).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateFileNoPermissions(t *testing.T) {
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"On-Behalf-Of": "blalbrit@stanford.edu", // The dummy authZ service is set to only allow lmcrae
			"Content-Type": contentType,
		}).
		SetBody(body).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateFileNoFileset(t *testing.T) {
	// Trying to attach to a nonexistent fileset returns a 404
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
			"Content-Type": contentType,
		}).
		SetBody(body).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
				assert.Contains(t, r.Body.String(), "Validation Error")
			})
}

func TestCreateFileFailure(t *testing.T) {
	// Failure in storing the binary causes a panic
	r := gofight.New()
	fileSet := datautils.NewResource(nil).
		WithType(datautils.FilesetType).
		WithID(filesetID)
	repo := NewMockDatabase(fileSet, nil)
	storage := NewMockErrorStorage()
	assert.Panics(t,
		func() {
			r.POST(filePath).
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
					"Content-Type": contentType,
				}).
				SetBody(body).
				Run(handler(repo, storage, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}

func TestCreateFileResourceFailure(t *testing.T) {
	// Failure in storing the metadata causes a panic
	r := gofight.New()
	fileSet := datautils.NewResource(nil).
		WithType(datautils.FilesetType).
		WithID(filesetID)
	repo := NewMockErrorDatabase(fileSet)
	assert.Panics(t,
		func() {

			r.POST(filePath).
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
					"Content-Type": contentType,
				}).
				SetBody(body).
				Run(handler(repo, nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
