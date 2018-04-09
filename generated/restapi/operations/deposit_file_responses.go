// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sul-dlss-labs/taco/generated/models"
)

// DepositFileCreatedCode is the HTTP code returned for type DepositFileCreated
const DepositFileCreatedCode int = 201

/*DepositFileCreated TACO binary ingested, File management metadata created, & File processing started.

swagger:response depositFileCreated
*/
type DepositFileCreated struct {

	/*
	  In: Body
	*/
	Payload models.ResourceResponse `json:"body,omitempty"`
}

// NewDepositFileCreated creates DepositFileCreated with default headers values
func NewDepositFileCreated() *DepositFileCreated {
	return &DepositFileCreated{}
}

// WithPayload adds the payload to the deposit file created response
func (o *DepositFileCreated) WithPayload(payload models.ResourceResponse) *DepositFileCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the deposit file created response
func (o *DepositFileCreated) SetPayload(payload models.ResourceResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DepositFileCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DepositFileUnauthorizedCode is the HTTP code returned for type DepositFileUnauthorized
const DepositFileUnauthorizedCode int = 401

/*DepositFileUnauthorized You are not authorized to ingest a File into TACO.

swagger:response depositFileUnauthorized
*/
type DepositFileUnauthorized struct {
}

// NewDepositFileUnauthorized creates DepositFileUnauthorized with default headers values
func NewDepositFileUnauthorized() *DepositFileUnauthorized {
	return &DepositFileUnauthorized{}
}

// WriteResponse to the client
func (o *DepositFileUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// DepositFileNotFoundCode is the HTTP code returned for type DepositFileNotFound
const DepositFileNotFoundCode int = 404

/*DepositFileNotFound Resource not found. Check that the provide identifier is correct.

swagger:response depositFileNotFound
*/
type DepositFileNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewDepositFileNotFound creates DepositFileNotFound with default headers values
func NewDepositFileNotFound() *DepositFileNotFound {
	return &DepositFileNotFound{}
}

// WithPayload adds the payload to the deposit file not found response
func (o *DepositFileNotFound) WithPayload(payload *models.ErrorResponse) *DepositFileNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the deposit file not found response
func (o *DepositFileNotFound) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DepositFileNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DepositFileUnsupportedMediaTypeCode is the HTTP code returned for type DepositFileUnsupportedMediaType
const DepositFileUnsupportedMediaTypeCode int = 415

/*DepositFileUnsupportedMediaType Unsupported file type provided.

swagger:response depositFileUnsupportedMediaType
*/
type DepositFileUnsupportedMediaType struct {
}

// NewDepositFileUnsupportedMediaType creates DepositFileUnsupportedMediaType with default headers values
func NewDepositFileUnsupportedMediaType() *DepositFileUnsupportedMediaType {
	return &DepositFileUnsupportedMediaType{}
}

// WriteResponse to the client
func (o *DepositFileUnsupportedMediaType) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(415)
}

// DepositFileInternalServerErrorCode is the HTTP code returned for type DepositFileInternalServerError
const DepositFileInternalServerErrorCode int = 500

/*DepositFileInternalServerError This file could be ingested at this time by TACO.

swagger:response depositFileInternalServerError
*/
type DepositFileInternalServerError struct {
}

// NewDepositFileInternalServerError creates DepositFileInternalServerError with default headers values
func NewDepositFileInternalServerError() *DepositFileInternalServerError {
	return &DepositFileInternalServerError{}
}

// WriteResponse to the client
func (o *DepositFileInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
