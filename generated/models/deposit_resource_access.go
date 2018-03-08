// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DepositResourceAccess Access Metadata for the Resource.
// swagger:model depositResourceAccess
type DepositResourceAccess struct {

	// Access level for the resource.
	// Required: true
	Access *string `json:"access"`

	// Download level for the resource metadata.
	// Required: true
	Download *string `json:"download"`

	// deposit resource access
	DepositResourceAccess map[string]string `json:"-"`
}

// UnmarshalJSON unmarshals this object with additional properties from JSON
func (m *DepositResourceAccess) UnmarshalJSON(data []byte) error {
	// stage 1, bind the properties
	var stage1 struct {

		// Access level for the resource.
		// Required: true
		Access *string `json:"access"`

		// Download level for the resource metadata.
		// Required: true
		Download *string `json:"download"`
	}
	if err := json.Unmarshal(data, &stage1); err != nil {
		return err
	}
	var rcv DepositResourceAccess

	rcv.Access = stage1.Access

	rcv.Download = stage1.Download

	*m = rcv

	// stage 2, remove properties and add to map
	stage2 := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &stage2); err != nil {
		return err
	}

	delete(stage2, "access")

	delete(stage2, "download")

	// stage 3, add additional properties values
	if len(stage2) > 0 {
		result := make(map[string]string)
		for k, v := range stage2 {
			var toadd string
			if err := json.Unmarshal(v, &toadd); err != nil {
				return err
			}
			result[k] = toadd
		}
		m.DepositResourceAccess = result
	}

	return nil
}

// MarshalJSON marshals this object with additional properties into a JSON object
func (m DepositResourceAccess) MarshalJSON() ([]byte, error) {
	var stage1 struct {

		// Access level for the resource.
		// Required: true
		Access *string `json:"access"`

		// Download level for the resource metadata.
		// Required: true
		Download *string `json:"download"`
	}

	stage1.Access = m.Access

	stage1.Download = m.Download

	// make JSON object for known properties
	props, err := json.Marshal(stage1)
	if err != nil {
		return nil, err
	}

	if len(m.DepositResourceAccess) == 0 {
		return props, nil
	}

	// make JSON object for the additional properties
	additional, err := json.Marshal(m.DepositResourceAccess)
	if err != nil {
		return nil, err
	}

	if len(props) < 3 {
		return additional, nil
	}

	// concatenate the 2 objects
	props[len(props)-1] = ','
	return append(props, additional[1:]...), nil
}

// Validate validates this deposit resource access
func (m *DepositResourceAccess) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccess(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDownload(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var depositResourceAccessTypeAccessPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["world","stanford","location-based","citation-only","dark"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		depositResourceAccessTypeAccessPropEnum = append(depositResourceAccessTypeAccessPropEnum, v)
	}
}

const (
	// DepositResourceAccessAccessWorld captures enum value "world"
	DepositResourceAccessAccessWorld string = "world"
	// DepositResourceAccessAccessStanford captures enum value "stanford"
	DepositResourceAccessAccessStanford string = "stanford"
	// DepositResourceAccessAccessLocationBased captures enum value "location-based"
	DepositResourceAccessAccessLocationBased string = "location-based"
	// DepositResourceAccessAccessCitationOnly captures enum value "citation-only"
	DepositResourceAccessAccessCitationOnly string = "citation-only"
	// DepositResourceAccessAccessDark captures enum value "dark"
	DepositResourceAccessAccessDark string = "dark"
)

// prop value enum
func (m *DepositResourceAccess) validateAccessEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, depositResourceAccessTypeAccessPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DepositResourceAccess) validateAccess(formats strfmt.Registry) error {

	if err := validate.Required("access", "body", m.Access); err != nil {
		return err
	}

	// value enum
	if err := m.validateAccessEnum("access", "body", *m.Access); err != nil {
		return err
	}

	return nil
}

var depositResourceAccessTypeDownloadPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["world","stanford","location-based","citation-only","dark"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		depositResourceAccessTypeDownloadPropEnum = append(depositResourceAccessTypeDownloadPropEnum, v)
	}
}

const (
	// DepositResourceAccessDownloadWorld captures enum value "world"
	DepositResourceAccessDownloadWorld string = "world"
	// DepositResourceAccessDownloadStanford captures enum value "stanford"
	DepositResourceAccessDownloadStanford string = "stanford"
	// DepositResourceAccessDownloadLocationBased captures enum value "location-based"
	DepositResourceAccessDownloadLocationBased string = "location-based"
	// DepositResourceAccessDownloadCitationOnly captures enum value "citation-only"
	DepositResourceAccessDownloadCitationOnly string = "citation-only"
	// DepositResourceAccessDownloadDark captures enum value "dark"
	DepositResourceAccessDownloadDark string = "dark"
)

// prop value enum
func (m *DepositResourceAccess) validateDownloadEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, depositResourceAccessTypeDownloadPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DepositResourceAccess) validateDownload(formats strfmt.Registry) error {

	if err := validate.Required("download", "body", m.Download); err != nil {
		return err
	}

	// value enum
	if err := m.validateDownloadEnum("download", "body", *m.Download); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DepositResourceAccess) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepositResourceAccess) UnmarshalBinary(b []byte) error {
	var res DepositResourceAccess
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
