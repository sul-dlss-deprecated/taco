package handlers

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/uploaded"
	"github.com/sul-dlss-labs/taco/validators"
)

const atContext = "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld"
const fileType = "http://sdr.sul.stanford.edu/contexts/sdr3-file.jsonld"

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(rt *taco.Runtime) func(*gin.Context) {
	return func(c *gin.Context) {
		entry := &depositFileEntry{rt: rt}
		entry.Handle(c)
	}
}

type depositFileEntry struct {
	rt *taco.Runtime
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(c *gin.Context) {
	if c.ContentType() != "multipart/form-data" {
		c.AbortWithStatusJSON(422, "You must attach file as multipart/form-data")
	}
	validator := validators.NewDepositFileValidator(d.rt.Repository())
	fileHeader, _ := c.FormFile("upload")
	if err := validator.ValidateResource(fileHeader); err != nil {
		c.AbortWithError(422, err)
	}

	id, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}

	_, err = d.copyFileToStorage(id, fileHeader)
	if err != nil {
		panic(err)
	}

	if err := d.createFileResource(id, fileHeader.Filename); err != nil {
		panic(err)
	}

	location := fmt.Sprintf("/v1/resource/%s", id)
	c.Header("Location", location)
	c.JSON(201, map[string]string{"id": id})
}

func (d *depositFileEntry) copyFileToStorage(id string, file *multipart.FileHeader) (*string, error) {
	filename := file.Filename
	contentType := file.Header.Get("Content-Type")
	log.Printf("Saving file \"%s\" with content-type: %s", filename, contentType)

	data, err := file.Open()
	if err != nil {
		return nil, err
	}
	upload := uploaded.NewFile(filename, contentType, data)
	return d.rt.FileStorage().UploadFile(id, upload)
}

func (d *depositFileEntry) createFileResource(resourceID string, filename string) error {
	resource := d.buildPersistableResource(resourceID, filename)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositFileEntry) buildPersistableResource(resourceID string, filename string) persistence.Resource {
	resource := persistence.Resource{"id": resourceID}
	// TODO: Where should Access come from/default to?
	// access := "private"
	// preserve := false
	// resource.Access = models.ResourceAccess{Access: &access}
	// resource.AtContext = atContext
	// resource.AtType = fileType
	// resource.Label = filename
	// resource.Administrative = models.ResourceAdministrative{SdrPreserve: &preserve}
	// resource.Publish = false
	return resource
}
