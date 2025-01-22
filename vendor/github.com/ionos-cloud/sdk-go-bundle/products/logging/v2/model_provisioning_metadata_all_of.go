/*
 * IONOS Logging REST API
 *
 * The logging service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data.  This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the ProvisioningMetadataAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ProvisioningMetadataAllOf{}

// ProvisioningMetadataAllOf struct for ProvisioningMetadataAllOf
type ProvisioningMetadataAllOf struct {
	// The current state reported back by the pipeline.
	State *string `json:"state,omitempty"`
}

// NewProvisioningMetadataAllOf instantiates a new ProvisioningMetadataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProvisioningMetadataAllOf() *ProvisioningMetadataAllOf {
	this := ProvisioningMetadataAllOf{}

	return &this
}

// NewProvisioningMetadataAllOfWithDefaults instantiates a new ProvisioningMetadataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProvisioningMetadataAllOfWithDefaults() *ProvisioningMetadataAllOf {
	this := ProvisioningMetadataAllOf{}
	return &this
}

// GetState returns the State field value if set, zero value otherwise.
func (o *ProvisioningMetadataAllOf) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProvisioningMetadataAllOf) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *ProvisioningMetadataAllOf) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *ProvisioningMetadataAllOf) SetState(v string) {
	o.State = &v
}

func (o ProvisioningMetadataAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	return toSerialize, nil
}

type NullableProvisioningMetadataAllOf struct {
	value *ProvisioningMetadataAllOf
	isSet bool
}

func (v NullableProvisioningMetadataAllOf) Get() *ProvisioningMetadataAllOf {
	return v.value
}

func (v *NullableProvisioningMetadataAllOf) Set(val *ProvisioningMetadataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableProvisioningMetadataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableProvisioningMetadataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProvisioningMetadataAllOf(val *ProvisioningMetadataAllOf) *NullableProvisioningMetadataAllOf {
	return &NullableProvisioningMetadataAllOf{value: val, isSet: true}
}

func (v NullableProvisioningMetadataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProvisioningMetadataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
