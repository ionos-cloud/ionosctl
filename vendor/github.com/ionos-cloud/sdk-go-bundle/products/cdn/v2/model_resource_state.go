/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the ResourceState type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ResourceState{}

// ResourceState The current status of the resource.
type ResourceState struct {
	// Represents one of the possible states of the resource.
	State string `json:"state"`
	// A human readable message describing the current state. In case of an error, the message will contain a detailed error message.
	Message *string `json:"message,omitempty"`
}

// NewResourceState instantiates a new ResourceState object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResourceState(state string) *ResourceState {
	this := ResourceState{}

	this.State = state

	return &this
}

// NewResourceStateWithDefaults instantiates a new ResourceState object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResourceStateWithDefaults() *ResourceState {
	this := ResourceState{}
	return &this
}

// GetState returns the State field value
func (o *ResourceState) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *ResourceState) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *ResourceState) SetState(v string) {
	o.State = v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ResourceState) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResourceState) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ResourceState) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ResourceState) SetMessage(v string) {
	o.Message = &v
}

func (o ResourceState) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ResourceState) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["state"] = o.State
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	return toSerialize, nil
}

type NullableResourceState struct {
	value *ResourceState
	isSet bool
}

func (v NullableResourceState) Get() *ResourceState {
	return v.value
}

func (v *NullableResourceState) Set(val *ResourceState) {
	v.value = val
	v.isSet = true
}

func (v NullableResourceState) IsSet() bool {
	return v.isSet
}

func (v *NullableResourceState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResourceState(val *ResourceState) *NullableResourceState {
	return &NullableResourceState{value: val, isSet: true}
}

func (v NullableResourceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResourceState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
