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

// checks if the SecondaryZoneReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SecondaryZoneReadListAllOf{}

// SecondaryZoneReadListAllOf List of secondary zones
type SecondaryZoneReadListAllOf struct {
	Items []SecondaryZoneRead `json:"items"`
}

// NewSecondaryZoneReadListAllOf instantiates a new SecondaryZoneReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecondaryZoneReadListAllOf(items []SecondaryZoneRead) *SecondaryZoneReadListAllOf {
	this := SecondaryZoneReadListAllOf{}

	this.Items = items

	return &this
}

// NewSecondaryZoneReadListAllOfWithDefaults instantiates a new SecondaryZoneReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecondaryZoneReadListAllOfWithDefaults() *SecondaryZoneReadListAllOf {
	this := SecondaryZoneReadListAllOf{}
	return &this
}

// GetItems returns the Items field value
func (o *SecondaryZoneReadListAllOf) GetItems() []SecondaryZoneRead {
	if o == nil {
		var ret []SecondaryZoneRead
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneReadListAllOf) GetItemsOk() ([]SecondaryZoneRead, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *SecondaryZoneReadListAllOf) SetItems(v []SecondaryZoneRead) {
	o.Items = v
}

func (o SecondaryZoneReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["items"] = o.Items
	return toSerialize, nil
}

type NullableSecondaryZoneReadListAllOf struct {
	value *SecondaryZoneReadListAllOf
	isSet bool
}

func (v NullableSecondaryZoneReadListAllOf) Get() *SecondaryZoneReadListAllOf {
	return v.value
}

func (v *NullableSecondaryZoneReadListAllOf) Set(val *SecondaryZoneReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableSecondaryZoneReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableSecondaryZoneReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecondaryZoneReadListAllOf(val *SecondaryZoneReadListAllOf) *NullableSecondaryZoneReadListAllOf {
	return &NullableSecondaryZoneReadListAllOf{value: val, isSet: true}
}

func (v NullableSecondaryZoneReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecondaryZoneReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
