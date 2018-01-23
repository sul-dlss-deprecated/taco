package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRetrieve(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	//var controller = new(controllers.ResourceController)
	c, _ := gin.CreateTestContext(rec)
	//controller.Retrieve(c)

	c.Request, _ = http.NewRequest("GET", "/v1/resource/99", nil)
	fmt.Println("recorder", rec.Code)
	res := rec.Result()
	fmt.Println("results", res.StatusCode)
	assert.Equal(t, rec.Code, 200)
}
