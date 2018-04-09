// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDepositFileParams creates a new DepositFileParams object
// with the default values initialized.
func NewDepositFileParams() DepositFileParams {
	var ()
	return DepositFileParams{}
}

// DepositFileParams contains all the bound params for the deposit file operation
// typically these are obtained from a http.Request
//
// swagger:parameters depositFile
type DepositFileParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Fileset identifier. This points at the container for this file.
	  Required: true
	  In: path
	*/
	FilesetID string
	/*Binary to be added to an Object in TACO.
	  Required: true
	  In: formData
	*/
	Upload runtime.File
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *DepositFileParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return err
		} else if err := r.ParseForm(); err != nil {
			return err
		}
	}

	rFilesetID, rhkFilesetID, _ := route.Params.GetOK("FilesetID")
	if err := o.bindFilesetID(rFilesetID, rhkFilesetID, route.Formats); err != nil {
		res = append(res, err)
	}

	upload, uploadHeader, err := r.FormFile("upload")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "upload", err))
	} else if err := o.bindUpload(upload, uploadHeader); err != nil {
		res = append(res, err)
	} else {
		o.Upload = runtime.File{Data: upload, Header: uploadHeader}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DepositFileParams) bindFilesetID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.FilesetID = raw

	return nil
}

func (o *DepositFileParams) bindUpload(file multipart.File, header *multipart.FileHeader) error {

	return nil
}
