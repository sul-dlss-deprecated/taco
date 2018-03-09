package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
)

// NewHealthCheck will return the service health
func NewHealthCheck(rt *taco.Runtime) func(*gin.Context) {
	return func(c *gin.Context) {
		entry := &healthCheck{}
		entry.Handle(c)
	}
}

type healthCheck struct{}

// Handle the health check request
func (d *healthCheck) Handle(c *gin.Context) {
	c.JSON(200, map[string]string{"status": "OK"})
}
