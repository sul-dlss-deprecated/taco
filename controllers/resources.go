package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco/models"
)

// ResourceController -- The controller that manages resources
type ResourceController struct{}

var repo = new(models.Resource)

// Retrieve -- returns a resource by its identifier
func (u ResourceController) Retrieve(c *gin.Context) {
	log.Println("retrieving")

	if c.Param("id") != "" {
		resource, err := repo.GetByID(c.Param("id"))
		if err != nil {
			if fmt.Sprint(err) == "not found" {
				c.JSON(404, gin.H{"message": "not found"})
			} else {
				c.JSON(500, gin.H{"message": "Error retrieving resource", "error": fmt.Sprint(err)})
			}
			c.Abort()
			return
		}
		log.Println("Got here")
		c.JSON(200, gin.H{"message": "Resource found!", "resource": resource})
		return
	}
	c.JSON(400, gin.H{"message": "bad request"})
	c.Abort()
	return
}
