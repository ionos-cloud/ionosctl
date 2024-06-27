/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.15.4
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RecordEnsure struct for RecordEnsure
type RecordEnsure struct {
	Properties *Record `json:"properties"`
}

// NewRecordEnsure instantiates a new RecordEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRecordEnsure(properties Record) *RecordEnsure {
	this := RecordEnsure{}

	this.Properties = &properties

	return &this
}

// NewRecordEnsureWithDefaults instantiates a new RecordEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRecordEnsureWithDefaults() *RecordEnsure {
	this := RecordEnsure{}
	return &this
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Record will be returned
func (o *RecordEnsure) GetProperties() *Record {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RecordEnsure) GetPropertiesOk() (*Record, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *RecordEnsure) SetProperties(v Record) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *RecordEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o RecordEnsure) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableRecordEnsure struct {
	value *RecordEnsure
	isSet bool
}

func (v NullableRecordEnsure) Get() *RecordEnsure {
	return v.value
}

func (v *NullableRecordEnsure) Set(val *RecordEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableRecordEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableRecordEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRecordEnsure(val *RecordEnsure) *NullableRecordEnsure {
	return &NullableRecordEnsure{value: val, isSet: true}
}

func (v NullableRecordEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRecordEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
