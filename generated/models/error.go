// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Error error
// swagger:model Error
type Error struct {

	// a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// source
	Source *ErrorSource `json:"source,omitempty"`

	// a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem.
	Title string `json:"title,omitempty"`
}

// Validate validates this error
func (m *Error) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSource(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Error) validateSource(formats strfmt.Registry) error {

	if swag.IsZero(m.Source) { // not required
		return nil
	}

	if m.Source != nil {

		if err := m.Source.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("source")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Error) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Error) UnmarshalBinary(b []byte) error {
	var res Error
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
