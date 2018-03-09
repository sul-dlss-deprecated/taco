package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(rt *taco.Runtime) func(*gin.Context) {
	return func(c *gin.Context) {
		log.Println("*** GOT STUFF *** ")
		entry := &retrieveResourceEntry{repository: rt.Repository()}
		entry.Handle(c)
	}
}

// retrieveResourceEntry handles a request for finding & returning an entry
type retrieveResourceEntry struct {
	repository persistence.Repository
}

// Handle the delete entry request
func (d *retrieveResourceEntry) Handle(c *gin.Context) {
	log.Printf("*** HERE WE ARE (%s) ***", c.Param("id"))
	resource, err := d.repository.GetByID(c.Param("id"))
	if err == nil {
		response := buildResponse(resource)
		c.JSON(200, response)
	} else if err.Error() == "not found" {
		c.AbortWithError(404, err)
	} else {
		panic(err)
	}
}

// TODO: expand this mapping
func buildResponse(resource *persistence.Resource) *persistence.Resource {
	return resource
}
