package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

const filePath = "/v1/file"
const contentType = "multipart/form-data; boundary=------------------------a31e2ddd4b2c0d92"
const body = `--------------------------a31e2ddd4b2c0d92
Content-Disposition: form-data; name="upload"; filename="foo.txt"
Content-Type: text/plain

Hello

--------------------------a31e2ddd4b2c0d92--`

func TestCreateFileHappyPath(t *testing.T) {
	r := gofight.New()
	storage := NewMockStorage()
	repo := NewMockDatabase(nil)

	r.POST(filePath).
		SetHeader(gofight.H{
			"Authorization": "lmcrae@stanford.edu",
			"Content-Type":  contentType,
		}).
		SetBody(body).
		Run(handler(repo, storage),
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

			})
}

func TestCreateFileWrongContentType(t *testing.T) {
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"Authorization": "lmcrae@stanford.edu",
			"Content-Type":  "application/xml",
		}).
		SetBody(``).
		Run(handler(nil, NewMockStorage()),
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
		Run(handler(nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateFileNoPermissions(t *testing.T) {
	r := gofight.New()
	r.POST(filePath).
		SetHeader(gofight.H{
			"Authorization": "blalbrit@stanford.edu", // The dummy authZ service is set to only allow lmcrae
			"Content-Type":  contentType,
		}).
		SetBody(body).
		Run(handler(nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateFileFailure(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(nil)
	storage := NewMockErrorStorage()
	assert.Panics(t,
		func() {
			r.POST(filePath).
				SetHeader(gofight.H{
					"Authorization": "lmcrae@stanford.edu",
					"Content-Type":  contentType,
				}).
				SetBody(body).
				Run(handler(repo, storage),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}

func TestCreateFileResourceFailure(t *testing.T) {
	r := gofight.New()
	repo := NewMockErrorDatabase()
	assert.Panics(t,
		func() {

			r.POST(filePath).
				SetHeader(gofight.H{
					"Authorization": "lmcrae@stanford.edu",
					"Content-Type":  contentType,
				}).
				SetBody(body).
				Run(handler(repo, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
