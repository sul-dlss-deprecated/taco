package operations

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/loads"
)

type TacoAPI struct {
	specDoc          *loads.Document
	RetrieveResource (func(*gin.Context))
	DepositResource  (func(*gin.Context))
	UpdateResource   (func(*gin.Context))
	DepositFile      (func(*gin.Context))
	HealthCheck      (func(*gin.Context))
	GetProcessStatus (func(*gin.Context))
	DeleteResource   (func(*gin.Context))
}

// NewTacoAPI creates a new api instance
func NewTacoAPI(specDoc *loads.Document) *TacoAPI {
	return &TacoAPI{specDoc: specDoc}
}

// CallFuncByName given a string, return that method
func (api *TacoAPI) CallFuncByName(funcName string) (func(*gin.Context), error) {
	myClassValue := reflect.ValueOf(api)
	m := reflect.Indirect(myClassValue).FieldByName(funcName)

	if !m.IsValid() {
		return nil, fmt.Errorf("method \"%s\" was not defined ", funcName)
	}
	return m.Interface().(func(*gin.Context)), nil
}

func (api *TacoAPI) operationFor(path string, operationId string) func(*gin.Context) {
	funcName := strings.Title(operationId)
	result, err := api.CallFuncByName(funcName)
	if err != nil {
		log.Panicf("Unable to create route for GET %s%s because %s", api.specDoc.BasePath(), path, err)
	}
	return result
}

var re = regexp.MustCompile(`\{(.*)\}`)

func (api *TacoAPI) pathFor(key string) string {
	// replace {ID} with :id
	return re.ReplaceAllStringFunc(key, func(in string) string {
		return fmt.Sprintf(":%s", strings.ToLower(strings.Trim(in, "{}")))
	})
}

// Engine returns the gin.Engine for this API
func (api *TacoAPI) Engine() *gin.Engine {
	router := gin.Default()
	base := router.Group(api.specDoc.BasePath())
	{
		for key, value := range api.specDoc.Spec().Paths.Paths {
			if value.Get != nil {
				result := api.operationFor(key, value.Get.OperationProps.ID)
				base.GET(api.pathFor(key), result)
			}
			if value.Post != nil {
				result := api.operationFor(key, value.Post.OperationProps.ID)
				base.POST(api.pathFor(key), result)
			}
			if value.Patch != nil {
				result := api.operationFor(key, value.Patch.OperationProps.ID)
				base.PATCH(api.pathFor(key), result)
			}
			if value.Put != nil {
				result := api.operationFor(key, value.Put.OperationProps.ID)
				base.PUT(api.pathFor(key), result)
			}
			if value.Delete != nil {
				result := api.operationFor(key, value.Delete.OperationProps.ID)
				base.DELETE(api.pathFor(key), result)
			}
		}
	}
	return router
}
