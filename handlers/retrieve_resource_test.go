package handlers

// import (
// 	"net/http"
// 	"testing"
//
// 	"github.com/appleboy/gofight"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/sul-dlss-labs/taco/persistence"
// )
//
// func TestRetrieveHappyPath(t *testing.T) {
// 	r := gofight.New()
// 	repo := mockRepo(new(persistence.Resource))
// 	r.GET("/v1/resource/99").
// 		Run(setupFakeRuntime().WithRepository(repo).Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusOK, r.Code)
// 			})
// }
//
// func TestRetrieveNotFound(t *testing.T) {
// 	r := gofight.New()
// 	r.GET("/v1/resource/100").
// 		Run(setupFakeRuntime().Handler(),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusNotFound, r.Code)
// 			})
// }
//
// // TODO: Error handling
// // func TestRetrieveError(t *testing.T) {
// // 	r := gofight.New()
// // 	r.GET("/v1/resource/100").
// // 		Run(setupFakeRuntime(mockErrorRepo()),
// // 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// // 				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
// // 			})
// // }
