package operations

import "github.com/gin-gonic/gin"

type TacoAPI struct {
	RetrieveResourceHandler (func(*gin.Context))
	DepositResourceHandler  (func(*gin.Context))
	UpdateResourceHandler   (func(*gin.Context))
	DepositFileHandler      (func(*gin.Context))
	HealthCheckHandler      (func(*gin.Context))
}

func NewTacoAPI() *TacoAPI {
	return &TacoAPI{}
}
