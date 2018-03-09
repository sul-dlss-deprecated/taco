package handlers

import (
	"log"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/loads"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/operations"
)

// BuildAPI create new service API
func BuildAPI(rt *taco.Runtime) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RetrieveResource = NewRetrieveResource(rt)
	api.DepositResource = NewDepositResource(rt)
	api.UpdateResource = NewUpdateResource(rt)
	api.DepositFile = NewDepositFile(rt)
	api.HealthCheck = NewHealthCheck(rt)
	api.GetProcessStatus = func(c *gin.Context) { c.JSON(503, gin.H{"error": "not implemented"}) }
	api.DeleteResource = func(c *gin.Context) { c.JSON(503, gin.H{"error": "not implemented"}) }

	return api
}

//
// // BuildHandler sets up the middleware that wraps the API
// func BuildHandler(api *operations.TacoAPI) http.Handler {
// 	return alice.New(
// 		middleware.NewHoneyBadgerMW(),
// 		middleware.NewRecoveryMW(),
// 		middleware.NewRequestLoggerMW(),
// 	).Then(api.Serve(nil))
// }

func swaggerSpec() *loads.Document {
	_, filename, _, _ := runtime.Caller(0)
	swaggerDoc := path.Join(path.Dir(filename), "../swagger.json")
	swaggerSpec, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
