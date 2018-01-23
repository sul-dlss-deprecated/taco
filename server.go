package server

import (
	"./controllers"
	"github.com/gin-gonic/gin"
)

func main() {
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
