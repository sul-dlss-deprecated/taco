package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco/models"
)

// ResourceController -- The controller that manages resources
type ResourceController struct{}

var repo = new(models.Resource)

// Retrieve -- returns a resource by its identifier
func (u ResourceController) Retrieve(c *gin.Context) {
	if c.Param("id") != "" {
		resource, err := repo.GetByID(c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"message": "Error to retrieve resource", "error": err})
			c.Abort()
			return
		}
		c.JSON(200, gin.H{"message": "Resource found!", "resource": resource})
		return
	}
	c.JSON(400, gin.H{"message": "bad request"})
	c.Abort()
	return
}
