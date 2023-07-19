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

// DatabaseItems struct for DatabaseItems
type DatabaseItems struct {
	Items *[]DatabaseResource `json:"items"`
}

// NewDatabaseItems instantiates a new DatabaseItems object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabaseItems(items []DatabaseResource) *DatabaseItems {
	this := DatabaseItems{}

	this.Items = &items

	return &this
}

// NewDatabaseItemsWithDefaults instantiates a new DatabaseItems object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabaseItemsWithDefaults() *DatabaseItems {
	this := DatabaseItems{}
	return &this
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []DatabaseResource will be returned
func (o *DatabaseItems) GetItems() *[]DatabaseResource {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DatabaseItems) GetItemsOk() (*[]DatabaseResource, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *DatabaseItems) SetItems(v []DatabaseResource) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *DatabaseItems) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o DatabaseItems) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	return json.Marshal(toSerialize)
}

type NullableDatabaseItems struct {
	value *DatabaseItems
	isSet bool
}

func (v NullableDatabaseItems) Get() *DatabaseItems {
	return v.value
}

func (v *NullableDatabaseItems) Set(val *DatabaseItems) {
	v.value = val
	v.isSet = true
}

func (v NullableDatabaseItems) IsSet() bool {
	return v.isSet
}

func (v *NullableDatabaseItems) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDatabaseItems(val *DatabaseItems) *NullableDatabaseItems {
	return &NullableDatabaseItems{value: val, isSet: true}
}

func (v NullableDatabaseItems) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDatabaseItems) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
