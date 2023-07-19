/*
 * IONOS DBaaS PostgreSQL REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional PostgreSQL database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// UserProperties struct for UserProperties
type UserProperties struct {
	// The username of a given user.
	Username *string `json:"username"`
	// The password of a given user.
	Password *string `json:"password,omitempty"`
	// Describes whether this user is a system user or not. A system user cannot be updated or deleted.
	System *bool `json:"system,omitempty"`
}

// NewUserProperties instantiates a new UserProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserProperties(username string) *UserProperties {
	this := UserProperties{}

	this.Username = &username

	return &this
}

// NewUserPropertiesWithDefaults instantiates a new UserProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserPropertiesWithDefaults() *UserProperties {
	this := UserProperties{}
	return &this
}

// GetUsername returns the Username field value
// If the value is explicit nil, the zero value for string will be returned
func (o *UserProperties) GetUsername() *string {
	if o == nil {
		return nil
	}

	return o.Username

}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserProperties) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Username, true
}

// SetUsername sets field value
func (o *UserProperties) SetUsername(v string) {

	o.Username = &v

}

// HasUsername returns a boolean if a field has been set.
func (o *UserProperties) HasUsername() bool {
	if o != nil && o.Username != nil {
		return true
	}

	return false
}

// GetPassword returns the Password field value
// If the value is explicit nil, the zero value for string will be returned
func (o *UserProperties) GetPassword() *string {
	if o == nil {
		return nil
	}

	return o.Password

}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserProperties) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Password, true
}

// SetPassword sets field value
func (o *UserProperties) SetPassword(v string) {

	o.Password = &v

}

// HasPassword returns a boolean if a field has been set.
func (o *UserProperties) HasPassword() bool {
	if o != nil && o.Password != nil {
		return true
	}

	return false
}

// GetSystem returns the System field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *UserProperties) GetSystem() *bool {
	if o == nil {
		return nil
	}

	return o.System

}

// GetSystemOk returns a tuple with the System field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserProperties) GetSystemOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}

	return o.System, true
}

// SetSystem sets field value
func (o *UserProperties) SetSystem(v bool) {

	o.System = &v

}

// HasSystem returns a boolean if a field has been set.
func (o *UserProperties) HasSystem() bool {
	if o != nil && o.System != nil {
		return true
	}

	return false
}

func (o UserProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Username != nil {
		toSerialize["username"] = o.Username
	}

	if o.Password != nil {
		toSerialize["password"] = o.Password
	}

	if o.System != nil {
		toSerialize["system"] = o.System
	}

	return json.Marshal(toSerialize)
}

type NullableUserProperties struct {
	value *UserProperties
	isSet bool
}

func (v NullableUserProperties) Get() *UserProperties {
	return v.value
}

func (v *NullableUserProperties) Set(val *UserProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableUserProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableUserProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserProperties(val *UserProperties) *NullableUserProperties {
	return &NullableUserProperties{value: val, isSet: true}
}

func (v NullableUserProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
