package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/handlers"
)

func main() {
	rt, err := taco.NewRuntime(config.NewConfig())
	if err != nil {
		log.Fatalln(err)
	}

	server := createServer(rt)

	// serve API
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

func createServer(rt *taco.Runtime) *gin.Engine {
	api := handlers.BuildAPI(rt)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	v1 := router.Group("/v1")
	{
		v1.POST("/resource", api.DepositResourceHandler)
		v1.GET("/resource/:id", api.RetrieveResourceHandler)
		v1.PATCH("/resource/:id", api.UpdateResourceHandler)
		v1.POST("/file", api.DepositFileHandler)
		v1.GET("/status/:id", func(c *gin.Context) {
			c.AbortWithStatusJSON(501, map[string]string{"status": "error", "message": "Not implemented"})
		})
		v1.GET("/healthcheck", api.HealthCheckHandler)
	}
	return router
}
