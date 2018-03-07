package operations

import "github.com/gin-gonic/gin"

type TacoAPI struct {
	RetrieveResourceHandler (func(*gin.Context))
	DepositResourceHandler  (func(*gin.Context))
	UpdateResourceHandler   (func(*gin.Context))
	DepositFileHandler      (func(*gin.Context))
	HealthCheckHandler      (func(*gin.Context))
}

// NewTacoAPI creates a new api instance
func NewTacoAPI() *TacoAPI {
	return &TacoAPI{}
}

// Engine returns the gin.Engine for this API
func (api *TacoAPI) Engine() *gin.Engine {
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
