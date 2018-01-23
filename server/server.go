package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco/controllers"
)

// Init -- setup the server
func Init() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		resourceGroup := v1.Group("resource")
		{
			resource := new(controllers.ResourceController)
			resourceGroup.GET("/:id", resource.Retrieve)
		}
	}
	router.Run() // listen and serve on 0.0.0.0:8080
}
