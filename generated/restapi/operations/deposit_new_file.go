// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DepositNewFileHandlerFunc turns a function with the right signature into a deposit new file handler
type DepositNewFileHandlerFunc func(DepositNewFileParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DepositNewFileHandlerFunc) Handle(params DepositNewFileParams) middleware.Responder {
	return fn(params)
}

// DepositNewFileHandler interface for that can handle valid deposit new file params
type DepositNewFileHandler interface {
	Handle(DepositNewFileParams) middleware.Responder
}

// NewDepositNewFile creates a new http.Handler for the deposit new file operation
func NewDepositNewFile(ctx *middleware.Context, handler DepositNewFileHandler) *DepositNewFile {
	return &DepositNewFile{Context: ctx, Handler: handler}
}

/*DepositNewFile swagger:route POST /file depositNewFile

Deposit a new File (binary) into SDR.

Deposits a new File (binary) into SDR. Will return the SDR identifier for the File resource (aka the metadata object generated and persisted for File management).

*/
type DepositNewFile struct {
	Context *middleware.Context
	Handler DepositNewFileHandler
}

func (o *DepositNewFile) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDepositNewFileParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
