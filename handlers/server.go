package handlers

import (
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/operations"
)

// BuildAPI create new service API
func BuildAPI(rt *taco.Runtime) *operations.TacoAPI {
	api := operations.NewTacoAPI() //(swaggerSpec())
	api.RetrieveResourceHandler = NewRetrieveResource(rt)
	api.DepositResourceHandler = NewDepositResource(rt)
	api.UpdateResourceHandler = NewUpdateResource(rt)
	api.DepositFileHandler = NewDepositFile(rt)
	api.HealthCheckHandler = NewHealthCheck(rt)
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
//
// func swaggerSpec() *loads.Document {
// 	// load embedded swagger file
// 	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return swaggerSpec
// }
