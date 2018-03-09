package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(rt *taco.Runtime) func(*gin.Context) {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps/DepositResource.json")
	validator := validators.NewDepositResourceValidator(rt.Repository(), schemaPath)

	return func(c *gin.Context) {
		entry := &depositResourceEntry{rt: rt, validator: validator}
		entry.Handle(c)
	}
}

type depositResourceEntry struct {
	rt        *taco.Runtime
	validator *validators.DepositResourceValidator
}

// Handle the delete entry request
func (d *depositResourceEntry) Handle(c *gin.Context) {
	buff := new(bytes.Buffer)
	buff.ReadFrom(c.Request.Body)

	if err := d.validator.ValidateResource(buff.String()); err != nil {
		c.AbortWithError(422, err)
		return
	}

	var data gin.H

	if err := json.Unmarshal(buff.Bytes(), &data); err != nil {
		panic(err)
	}
	resourceID, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}
	if err := d.persistResource(resourceID, data); err != nil {
		// TODO: handle this with an error response

		panic(err)
	}

	if err := d.addToStream(&resourceID); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}
	c.JSON(201, map[string]string{"id": resourceID})
}

func (d *depositResourceEntry) persistResource(resourceID string, data gin.H) error {
	resource := d.persistableResourceFromParams(resourceID, data)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositResourceEntry) persistableResourceFromParams(resourceID string, data gin.H) persistence.Resource {
	resource := persistence.NewResource(data)
	resource["id"] = resourceID

	//resource.Access = map["Access"] : params.Payload.Access.Access}
	// resource.AtContext = *params.Payload.AtContext
	// resource.AtType = *params.Payload.AtType
	// resource.Label = *params.Payload.Label
	// TODO: ResourceIdentification has no SourceID?
	//resource.Identification = models.ResourceIdentification{SourceID: params.Payload.Identification.SourceID}
	return resource
}

func (d *depositResourceEntry) addToStream(id *string) error {
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
