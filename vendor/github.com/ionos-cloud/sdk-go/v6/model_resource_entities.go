/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ResourceEntities struct for ResourceEntities
type ResourceEntities struct {
	Groups *ResourceGroups `json:"groups,omitempty"`
}



// GetGroups returns the Groups field value
// If the value is explicit nil, the zero value for ResourceGroups will be returned
func (o *ResourceEntities) GetGroups() *ResourceGroups {
	if o == nil {
		return nil
	}


	return o.Groups

}

// GetGroupsOk returns a tuple with the Groups field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceEntities) GetGroupsOk() (*ResourceGroups, bool) {
	if o == nil {
		return nil, false
	}


	return o.Groups, true
}

// SetGroups sets field value
func (o *ResourceEntities) SetGroups(v ResourceGroups) {


	o.Groups = &v

}

// HasGroups returns a boolean if a field has been set.
func (o *ResourceEntities) HasGroups() bool {
	if o != nil && o.Groups != nil {
		return true
	}

	return false
}


func (o ResourceEntities) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Groups != nil {
		toSerialize["groups"] = o.Groups
	}
	
	return json.Marshal(toSerialize)
}

type NullableResourceEntities struct {
	value *ResourceEntities
	isSet bool
}

func (v NullableResourceEntities) Get() *ResourceEntities {
	return v.value
}

func (v *NullableResourceEntities) Set(val *ResourceEntities) {
	v.value = val
	v.isSet = true
}

func (v NullableResourceEntities) IsSet() bool {
	return v.isSet
}

func (v *NullableResourceEntities) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResourceEntities(val *ResourceEntities) *NullableResourceEntities {
	return &NullableResourceEntities{value: val, isSet: true}
}

func (v NullableResourceEntities) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResourceEntities) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


