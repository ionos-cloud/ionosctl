/*
 * Auth API
 *
 * Use the Auth API to manage tokens for secure access to IONOS Cloud  APIs (Auth API, Cloud API, Reseller API, Activity Log API, and others).
 *
 * API version: 1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// DeleteResponse struct for DeleteResponse
type DeleteResponse struct {
	Success *bool `json:"success,omitempty"`
}

// GetSuccess returns the Success field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *DeleteResponse) GetSuccess() *bool {
	if o == nil {
		return nil
	}

	return o.Success

}

// GetSuccessOk returns a tuple with the Success field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DeleteResponse) GetSuccessOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}

	return o.Success, true
}

// SetSuccess sets field value
func (o *DeleteResponse) SetSuccess(v bool) {

	o.Success = &v

}

// HasSuccess returns a boolean if a field has been set.
func (o *DeleteResponse) HasSuccess() bool {
	if o != nil && o.Success != nil {
		return true
	}

	return false
}

func (o DeleteResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Success != nil {
		toSerialize["success"] = o.Success
	}
	return json.Marshal(toSerialize)
}

type NullableDeleteResponse struct {
	value *DeleteResponse
	isSet bool
}

func (v NullableDeleteResponse) Get() *DeleteResponse {
	return v.value
}

func (v *NullableDeleteResponse) Set(val *DeleteResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableDeleteResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableDeleteResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDeleteResponse(val *DeleteResponse) *NullableDeleteResponse {
	return &NullableDeleteResponse{value: val, isSet: true}
}

func (v NullableDeleteResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDeleteResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
