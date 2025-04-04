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

// checks if the ReverseRecordEnsure type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReverseRecordEnsure{}

// ReverseRecordEnsure struct for ReverseRecordEnsure
type ReverseRecordEnsure struct {
	Properties ReverseRecord `json:"properties"`
}

// NewReverseRecordEnsure instantiates a new ReverseRecordEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReverseRecordEnsure(properties ReverseRecord) *ReverseRecordEnsure {
	this := ReverseRecordEnsure{}

	this.Properties = properties

	return &this
}

// NewReverseRecordEnsureWithDefaults instantiates a new ReverseRecordEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReverseRecordEnsureWithDefaults() *ReverseRecordEnsure {
	this := ReverseRecordEnsure{}
	return &this
}

// GetProperties returns the Properties field value
func (o *ReverseRecordEnsure) GetProperties() ReverseRecord {
	if o == nil {
		var ret ReverseRecord
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordEnsure) GetPropertiesOk() (*ReverseRecord, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *ReverseRecordEnsure) SetProperties(v ReverseRecord) {
	o.Properties = v
}

func (o ReverseRecordEnsure) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableReverseRecordEnsure struct {
	value *ReverseRecordEnsure
	isSet bool
}

func (v NullableReverseRecordEnsure) Get() *ReverseRecordEnsure {
	return v.value
}

func (v *NullableReverseRecordEnsure) Set(val *ReverseRecordEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableReverseRecordEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableReverseRecordEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReverseRecordEnsure(val *ReverseRecordEnsure) *NullableReverseRecordEnsure {
	return &NullableReverseRecordEnsure{value: val, isSet: true}
}

func (v NullableReverseRecordEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReverseRecordEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
