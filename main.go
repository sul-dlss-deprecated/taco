package main

//go:generate go run generate/generate.go

import (
	"log"
	"net/http"

	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/handlers"
	"github.com/sul-dlss-labs/taco/middleware"
	"github.com/sul-dlss-labs/taco/runtime"
)

func main() {
	// Initialize our global struct
	config := config.NewConfig()
	rt := runtime.NewRuntime(config)

	server := createServer(rt, config.Port)
	defer server.Shutdown()

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func createServer(rt *runtime.Runtime, port int) *restapi.Server {
	api := handlers.BuildAPI(rt)
	server := restapi.NewServer(api)
	server.SetHandler(BuildHandler(api))

	// set the port this service will be run on
	server.Port = port
	return server
}

// BuildHandler sets up the middleware that wraps the API
func BuildHandler(api *operations.TacoAPI) http.Handler {
	return alice.New(
		middleware.NewHoneyBadgerMW(),
		middleware.NewRecoveryMW(),
		middleware.NewRequestLoggerMW(),
	).Then(api.Serve(nil))
}
