// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../generated --name  --spec ../swagger.json --principal authorization.Agent --exclude-main

func configureFlags(api *operations.TacoAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TacoAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "On-Behalf-Of" header is set
	api.RemoteUserAuth = func(token string) (*authorization.Agent, error) {
		return nil, errors.NotImplemented("api key auth (RemoteUser) On-Behalf-Of from header param [On-Behalf-Of] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.DeleteResourceHandler = operations.DeleteResourceHandlerFunc(func(params operations.DeleteResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteResource has not yet been implemented")
	})
	api.DepositFileHandler = operations.DepositFileHandlerFunc(func(params operations.DepositFileParams, principal *authorization.Agent) middleware.Responder {
		return middleware.NotImplemented("operation .DepositFile has not yet been implemented")
	})
	api.DepositResourceHandler = operations.DepositResourceHandlerFunc(func(params operations.DepositResourceParams, principal *authorization.Agent) middleware.Responder {
		return middleware.NotImplemented("operation .DepositResource has not yet been implemented")
	})
	api.GetProcessStatusHandler = operations.GetProcessStatusHandlerFunc(func(params operations.GetProcessStatusParams, principal *authorization.Agent) middleware.Responder {
		return middleware.NotImplemented("operation .GetProcessStatus has not yet been implemented")
	})
	api.HealthCheckHandler = operations.HealthCheckHandlerFunc(func(params operations.HealthCheckParams) middleware.Responder {
		return middleware.NotImplemented("operation .HealthCheck has not yet been implemented")
	})
	api.RetrieveResourceHandler = operations.RetrieveResourceHandlerFunc(func(params operations.RetrieveResourceParams, principal *authorization.Agent) middleware.Responder {
		return middleware.NotImplemented("operation .RetrieveResource has not yet been implemented")
	})
	api.UpdateResourceHandler = operations.UpdateResourceHandlerFunc(func(params operations.UpdateResourceParams, principal *authorization.Agent) middleware.Responder {
		return middleware.NotImplemented("operation .UpdateResource has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
