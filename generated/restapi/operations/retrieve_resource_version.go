// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// RetrieveResourceVersionHandlerFunc turns a function with the right signature into a retrieve resource version handler
type RetrieveResourceVersionHandlerFunc func(RetrieveResourceVersionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RetrieveResourceVersionHandlerFunc) Handle(params RetrieveResourceVersionParams) middleware.Responder {
	return fn(params)
}

// RetrieveResourceVersionHandler interface for that can handle valid retrieve resource version params
type RetrieveResourceVersionHandler interface {
	Handle(RetrieveResourceVersionParams) middleware.Responder
}

// NewRetrieveResourceVersion creates a new http.Handler for the retrieve resource version operation
func NewRetrieveResourceVersion(ctx *middleware.Context, handler RetrieveResourceVersionHandler) *RetrieveResourceVersion {
	return &RetrieveResourceVersion{Context: ctx, Handler: handler}
}

/*RetrieveResourceVersion swagger:route GET /resource/{ID}/{Version} retrieveResourceVersion

Retrieve TACO Resource metadata for a specific version.

Retrieves the metadata (as JSON-LD following our SDR3 MAP v.1) for an existing version of a TACO resource (Collection, Digital Repository Object, File metadata object [not binary] or subclass of those). The resource is identified by the TACO identifier.

*/
type RetrieveResourceVersion struct {
	Context *middleware.Context
	Handler RetrieveResourceVersionHandler
}

func (o *RetrieveResourceVersion) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRetrieveResourceVersionParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}