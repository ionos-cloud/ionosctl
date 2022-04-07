/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// UserPut struct for UserPut
type UserPut struct {
	// The resource's unique identifier.
	Id         *string            `json:"id,omitempty"`
	Properties *UserPropertiesPut `json:"properties"`
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *UserPut) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserPut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *UserPut) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *UserPut) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for UserPropertiesPut will be returned
func (o *UserPut) GetProperties() *UserPropertiesPut {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserPut) GetPropertiesOk() (*UserPropertiesPut, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *UserPut) SetProperties(v UserPropertiesPut) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *UserPut) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o UserPut) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}
	return json.Marshal(toSerialize)
}

type NullableUserPut struct {
	value *UserPut
	isSet bool
}

func (v NullableUserPut) Get() *UserPut {
	return v.value
}

func (v *NullableUserPut) Set(val *UserPut) {
	v.value = val
	v.isSet = true
}

func (v NullableUserPut) IsSet() bool {
	return v.isSet
}

func (v *NullableUserPut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserPut(val *UserPut) *NullableUserPut {
	return &NullableUserPut{value: val, isSet: true}
}

func (v NullableUserPut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserPut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
