package handlers

// import (
// 	"net/http"
// 	"testing"
//
// 	"github.com/appleboy/gofight"
// 	"github.com/stretchr/testify/assert"
// )
//
// const path = "/v1/file"
// const contentType = "multipart/form-data; boundary=------------------------a31e2ddd4b2c0d92"
// const body = `--------------------------a31e2ddd4b2c0d92
// Content-Disposition: form-data; name="upload"; filename="foo.txt"
// Content-Type: text/plain
//
// Hello
//
// --------------------------a31e2ddd4b2c0d92--`
//
// func TestCreateFileHappyPath(t *testing.T) {
// 	r := gofight.New()
// 	storage := mockStorage()
// 	repo := mockRepo(nil)
// 	handler := setupFakeRuntime().
// 		WithRepository(repo).
// 		WithStorage(storage).
// 		Handler()
//
// 	r.POST(path).
// 		SetHeader(gofight.H{
// 			"Content-Type": contentType,
// 		}).
// 		SetBody(body).
// 		Run(handler,
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusCreated, r.Code)
// 				assert.Equal(t, 1, len(storage.(*fakeStorage).CreatedFiles))
// 				assert.Equal(t, 1, len(repo.(*fakeRepository).CreatedResources))
// 			})
// }
//
// func TestCreateFileWrongContentType(t *testing.T) {
// 	r := gofight.New()
// 	r.POST(path).
// 		SetHeader(gofight.H{
// 			"Content-Type": "application/xml",
// 		}).
// 		SetBody(``).
// 		Run(setupFakeRuntime().Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusBadRequest, r.Code)
// 			})
// }
//
// // TODO: https://github.com/go-swagger/go-swagger/issues/1400
// // func TestCreateFileMissingFile(t *testing.T) {
// // 	r := gofight.New()
// // 	r.POST(path).
// // 		SetHeader(gofight.H{
// // 			"Content-Type": contentType,
// // 		}).
// // 		SetBody(``).
// // 		Run(setupFakeRuntime().Handler(),
// // 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// // 				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
// // 			})
// // }
//
// func TestCreateFileFailure(t *testing.T) {
// 	r := gofight.New()
// 	storage := mockErrorStorage()
// 	r.POST(path).
// 		SetHeader(gofight.H{
// 			"Content-Type": contentType,
// 		}).
// 		SetBody(body).
// 		Run(setupFakeRuntime().WithStorage(storage).Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusInternalServerError, r.Code)
// 			})
// }
//
// func TestCreateFileResourceFailure(t *testing.T) {
// 	r := gofight.New()
// 	repo := mockErrorRepo()
// 	r.POST(path).
// 		SetHeader(gofight.H{
// 			"Content-Type": contentType,
// 		}).
// 		SetBody(body).
// 		Run(setupFakeRuntime().WithRepository(repo).Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusInternalServerError, r.Code)
// 			})
// }
