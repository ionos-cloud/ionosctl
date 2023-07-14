/*
 * IONOS Cloud - DNS API
 *
 * DNS API Specification
 *
 * API version: 1.2.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RecordCreate struct for RecordCreate
type RecordCreate struct {
	Properties *Record `json:"properties"`
}

// NewRecordCreate instantiates a new RecordCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRecordCreate(properties Record) *RecordCreate {
	this := RecordCreate{}

	this.Properties = &properties

	return &this
}

// NewRecordCreateWithDefaults instantiates a new RecordCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRecordCreateWithDefaults() *RecordCreate {
	this := RecordCreate{}
	return &this
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Record will be returned
func (o *RecordCreate) GetProperties() *Record {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RecordCreate) GetPropertiesOk() (*Record, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *RecordCreate) SetProperties(v Record) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *RecordCreate) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o RecordCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableRecordCreate struct {
	value *RecordCreate
	isSet bool
}

func (v NullableRecordCreate) Get() *RecordCreate {
	return v.value
}

func (v *NullableRecordCreate) Set(val *RecordCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableRecordCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableRecordCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRecordCreate(val *RecordCreate) *NullableRecordCreate {
	return &NullableRecordCreate{value: val, isSet: true}
}

func (v NullableRecordCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRecordCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
