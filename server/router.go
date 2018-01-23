package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco/controllers"
)

// NewRouter -- The engine
func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		resourceGroup := v1.Group("resource")
		{
			resource := new(controllers.ResourceController)
			resourceGroup.GET("/:id", resource.Retrieve)
		}
	}
	return router
}
