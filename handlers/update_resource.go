package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(rt *taco.Runtime) func(*gin.Context) {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps/Resource.json")
	validator := validators.NewUpdateResourceValidator(rt.Repository(), schemaPath)

	return func(c *gin.Context) {
		entry := &updateResourceEntry{rt: rt, validator: validator}
		entry.Handle(c)
	}
}

type updateResourceEntry struct {
	rt        *taco.Runtime
	validator *validators.UpdateResourceValidator
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(c *gin.Context) {
	buff := new(bytes.Buffer)
	buff.ReadFrom(c.Request.Body)
	s := buff.String()

	if err := d.validator.ValidateResource(s); err != nil {
		c.AbortWithError(422, err)
		return
	}

	id := c.Param("id")
	resource, err := d.rt.Repository().GetByID(id)

	if err != nil {
		if err.Error() == "not found" {
			c.AbortWithError(404, err)
			return
		}
		panic(err)
	}
	var data gin.H

	if err = json.Unmarshal(buff.Bytes(), &data); err != nil {
		panic(err)
	}

	if err = d.updateResource(id, data); err != nil {
		panic(err)
	}

	if err = d.addToStream(id); err != nil {
		panic(err)
	}

	c.JSON(200, resource)
}

func (d *updateResourceEntry) updateResource(resourceID string, data gin.H) error {
	resource := d.persistableResourceFromParams(resourceID, data)
	return d.rt.Repository().UpdateItem(resource)
}

func (d *updateResourceEntry) persistableResourceFromParams(resourceID string, data gin.H) persistence.Resource {
	resource := persistence.NewResource(data)
	resource["id"] = resourceID
	// TODO: expand this mapping
	// resource.Access = *params.Payload.Access
	// resource.AtContext = *params.Payload.AtContext
	// resource.AtType = *params.Payload.AtType
	// resource.Label = *params.Payload.Label
	// resource.Administrative = *params.Payload.Administrative
	// resource.Identification = *params.Payload.Identification
	// resource.Publish = *params.Payload.Publish
	return resource
}

func (d *updateResourceEntry) addToStream(id string) error {
	message, err := json.Marshal(id)
	if err != nil {
		return err
	}
	if d.rt.Stream() == nil {
		log.Printf("Stream is nil")
	}
	if err := d.rt.Stream().SendMessage(string(message)); err != nil {
		return err
	}
	return nil
}
