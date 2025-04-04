/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.16.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the DnssecKeyReadListPropertiesNsecParameters type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DnssecKeyReadListPropertiesNsecParameters{}

// DnssecKeyReadListPropertiesNsecParameters struct for DnssecKeyReadListPropertiesNsecParameters
type DnssecKeyReadListPropertiesNsecParameters struct {
	NsecMode *NsecMode `json:"nsecMode,omitempty"`
}

// NewDnssecKeyReadListPropertiesNsecParameters instantiates a new DnssecKeyReadListPropertiesNsecParameters object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnssecKeyReadListPropertiesNsecParameters() *DnssecKeyReadListPropertiesNsecParameters {
	this := DnssecKeyReadListPropertiesNsecParameters{}

	return &this
}

// NewDnssecKeyReadListPropertiesNsecParametersWithDefaults instantiates a new DnssecKeyReadListPropertiesNsecParameters object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnssecKeyReadListPropertiesNsecParametersWithDefaults() *DnssecKeyReadListPropertiesNsecParameters {
	this := DnssecKeyReadListPropertiesNsecParameters{}
	return &this
}

// GetNsecMode returns the NsecMode field value if set, zero value otherwise.
func (o *DnssecKeyReadListPropertiesNsecParameters) GetNsecMode() NsecMode {
	if o == nil || IsNil(o.NsecMode) {
		var ret NsecMode
		return ret
	}
	return *o.NsecMode
}

// GetNsecModeOk returns a tuple with the NsecMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadListPropertiesNsecParameters) GetNsecModeOk() (*NsecMode, bool) {
	if o == nil || IsNil(o.NsecMode) {
		return nil, false
	}
	return o.NsecMode, true
}

// HasNsecMode returns a boolean if a field has been set.
func (o *DnssecKeyReadListPropertiesNsecParameters) HasNsecMode() bool {
	if o != nil && !IsNil(o.NsecMode) {
		return true
	}

	return false
}

// SetNsecMode gets a reference to the given NsecMode and assigns it to the NsecMode field.
func (o *DnssecKeyReadListPropertiesNsecParameters) SetNsecMode(v NsecMode) {
	o.NsecMode = &v
}

func (o DnssecKeyReadListPropertiesNsecParameters) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.NsecMode) {
		toSerialize["nsecMode"] = o.NsecMode
	}
	return toSerialize, nil
}

type NullableDnssecKeyReadListPropertiesNsecParameters struct {
	value *DnssecKeyReadListPropertiesNsecParameters
	isSet bool
}

func (v NullableDnssecKeyReadListPropertiesNsecParameters) Get() *DnssecKeyReadListPropertiesNsecParameters {
	return v.value
}

func (v *NullableDnssecKeyReadListPropertiesNsecParameters) Set(val *DnssecKeyReadListPropertiesNsecParameters) {
	v.value = val
	v.isSet = true
}

func (v NullableDnssecKeyReadListPropertiesNsecParameters) IsSet() bool {
	return v.isSet
}

func (v *NullableDnssecKeyReadListPropertiesNsecParameters) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnssecKeyReadListPropertiesNsecParameters(val *DnssecKeyReadListPropertiesNsecParameters) *NullableDnssecKeyReadListPropertiesNsecParameters {
	return &NullableDnssecKeyReadListPropertiesNsecParameters{value: val, isSet: true}
}

func (v NullableDnssecKeyReadListPropertiesNsecParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnssecKeyReadListPropertiesNsecParameters) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
