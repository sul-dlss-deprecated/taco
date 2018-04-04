// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewRetrieveResourceVersionParams creates a new RetrieveResourceVersionParams object
// with the default values initialized.
func NewRetrieveResourceVersionParams() RetrieveResourceVersionParams {
	var ()
	return RetrieveResourceVersionParams{}
}

// RetrieveResourceVersionParams contains all the bound params for the retrieve resource version operation
// typically these are obtained from a http.Request
//
// swagger:parameters retrieveResourceVersion
type RetrieveResourceVersionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*TACO Resource Identifier.
	  Required: true
	  In: path
	*/
	ID string
	/*TACO resource version number.
	  Required: true
	  In: path
	*/
	Version string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *RetrieveResourceVersionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("ID")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	rVersion, rhkVersion, _ := route.Params.GetOK("Version")
	if err := o.bindVersion(rVersion, rhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RetrieveResourceVersionParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ID = raw

	return nil
}

func (o *RetrieveResourceVersionParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.Version = raw

	return nil
}
