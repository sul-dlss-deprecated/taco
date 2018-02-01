// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sul-dlss-labs/taco/generated/models"
)

// GetProcessStatusOKCode is the HTTP code returned for type GetProcessStatusOK
const GetProcessStatusOKCode int = 200

/*GetProcessStatusOK Processing status for the TACO resource.

swagger:response getProcessStatusOK
*/
type GetProcessStatusOK struct {

	/*
	  In: Body
	*/
	Payload *models.ProcessResponse `json:"body,omitempty"`
}

// NewGetProcessStatusOK creates GetProcessStatusOK with default headers values
func NewGetProcessStatusOK() *GetProcessStatusOK {
	return &GetProcessStatusOK{}
}

// WithPayload adds the payload to the get process status o k response
func (o *GetProcessStatusOK) WithPayload(payload *models.ProcessResponse) *GetProcessStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get process status o k response
func (o *GetProcessStatusOK) SetPayload(payload *models.ProcessResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProcessStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProcessStatusUnauthorizedCode is the HTTP code returned for type GetProcessStatusUnauthorized
const GetProcessStatusUnauthorizedCode int = 401

/*GetProcessStatusUnauthorized You are not authorized to view this resource's processing status in TACO.

swagger:response getProcessStatusUnauthorized
*/
type GetProcessStatusUnauthorized struct {
}

// NewGetProcessStatusUnauthorized creates GetProcessStatusUnauthorized with default headers values
func NewGetProcessStatusUnauthorized() *GetProcessStatusUnauthorized {
	return &GetProcessStatusUnauthorized{}
}

// WriteResponse to the client
func (o *GetProcessStatusUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetProcessStatusNotFoundCode is the HTTP code returned for type GetProcessStatusNotFound
const GetProcessStatusNotFoundCode int = 404

/*GetProcessStatusNotFound Resource not found. Please check your provided TACO identifier.

swagger:response getProcessStatusNotFound
*/
type GetProcessStatusNotFound struct {
}

// NewGetProcessStatusNotFound creates GetProcessStatusNotFound with default headers values
func NewGetProcessStatusNotFound() *GetProcessStatusNotFound {
	return &GetProcessStatusNotFound{}
}

// WriteResponse to the client
func (o *GetProcessStatusNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetProcessStatusInternalServerErrorCode is the HTTP code returned for type GetProcessStatusInternalServerError
const GetProcessStatusInternalServerErrorCode int = 500

/*GetProcessStatusInternalServerError This resource's processing status could be retrieved at this time by TACO.

swagger:response getProcessStatusInternalServerError
*/
type GetProcessStatusInternalServerError struct {
}

// NewGetProcessStatusInternalServerError creates GetProcessStatusInternalServerError with default headers values
func NewGetProcessStatusInternalServerError() *GetProcessStatusInternalServerError {
	return &GetProcessStatusInternalServerError{}
}

// WriteResponse to the client
func (o *GetProcessStatusInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}