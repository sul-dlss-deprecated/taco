// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DepositNewResourceHandlerFunc turns a function with the right signature into a deposit new resource handler
type DepositNewResourceHandlerFunc func(DepositNewResourceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DepositNewResourceHandlerFunc) Handle(params DepositNewResourceParams) middleware.Responder {
	return fn(params)
}

// DepositNewResourceHandler interface for that can handle valid deposit new resource params
type DepositNewResourceHandler interface {
	Handle(DepositNewResourceParams) middleware.Responder
}

// NewDepositNewResource creates a new http.Handler for the deposit new resource operation
func NewDepositNewResource(ctx *middleware.Context, handler DepositNewResourceHandler) *DepositNewResource {
	return &DepositNewResource{Context: ctx, Handler: handler}
}

/*DepositNewResource swagger:route POST /resource depositNewResource

Deposit a new resource into SDR.

Deposits a new resource (Collection, Digital Repository Object, Fileset, or subclass of those) into SDR. Will return the SDR identifier for the resource.

*/
type DepositNewResource struct {
	Context *middleware.Context
	Handler DepositNewResourceHandler
}

func (o *DepositNewResource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDepositNewResourceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
