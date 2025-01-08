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

// checks if the SecondaryZoneEnsure type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SecondaryZoneEnsure{}

// SecondaryZoneEnsure struct for SecondaryZoneEnsure
type SecondaryZoneEnsure struct {
	Properties SecondaryZone `json:"properties"`
}

// NewSecondaryZoneEnsure instantiates a new SecondaryZoneEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecondaryZoneEnsure(properties SecondaryZone) *SecondaryZoneEnsure {
	this := SecondaryZoneEnsure{}

	this.Properties = properties

	return &this
}

// NewSecondaryZoneEnsureWithDefaults instantiates a new SecondaryZoneEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecondaryZoneEnsureWithDefaults() *SecondaryZoneEnsure {
	this := SecondaryZoneEnsure{}
	return &this
}

// GetProperties returns the Properties field value
func (o *SecondaryZoneEnsure) GetProperties() SecondaryZone {
	if o == nil {
		var ret SecondaryZone
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneEnsure) GetPropertiesOk() (*SecondaryZone, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *SecondaryZoneEnsure) SetProperties(v SecondaryZone) {
	o.Properties = v
}

func (o SecondaryZoneEnsure) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableSecondaryZoneEnsure struct {
	value *SecondaryZoneEnsure
	isSet bool
}

func (v NullableSecondaryZoneEnsure) Get() *SecondaryZoneEnsure {
	return v.value
}

func (v *NullableSecondaryZoneEnsure) Set(val *SecondaryZoneEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableSecondaryZoneEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableSecondaryZoneEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecondaryZoneEnsure(val *SecondaryZoneEnsure) *NullableSecondaryZoneEnsure {
	return &NullableSecondaryZoneEnsure{value: val, isSet: true}
}

func (v NullableSecondaryZoneEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecondaryZoneEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
