/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.15.4
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// MetadataWithStateFqdnZoneIdAllOf struct for MetadataWithStateFqdnZoneIdAllOf
type MetadataWithStateFqdnZoneIdAllOf struct {
	State *ProvisioningState `json:"state"`
	// A fully qualified domain name. FQDN consists of two parts - the hostname and the domain name.
	Fqdn *string `json:"fqdn"`
	// The ID (UUID) of the DNS zone of which record belongs to.
	ZoneId *string `json:"zoneId"`
}

// NewMetadataWithStateFqdnZoneIdAllOf instantiates a new MetadataWithStateFqdnZoneIdAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataWithStateFqdnZoneIdAllOf(state ProvisioningState, fqdn string, zoneId string) *MetadataWithStateFqdnZoneIdAllOf {
	this := MetadataWithStateFqdnZoneIdAllOf{}

	this.State = &state
	this.Fqdn = &fqdn
	this.ZoneId = &zoneId

	return &this
}

// NewMetadataWithStateFqdnZoneIdAllOfWithDefaults instantiates a new MetadataWithStateFqdnZoneIdAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithStateFqdnZoneIdAllOfWithDefaults() *MetadataWithStateFqdnZoneIdAllOf {
	this := MetadataWithStateFqdnZoneIdAllOf{}
	return &this
}

// GetState returns the State field value
// If the value is explicit nil, the zero value for ProvisioningState will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetState() *ProvisioningState {
	if o == nil {
		return nil
	}

	return o.State

}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetStateOk() (*ProvisioningState, bool) {
	if o == nil {
		return nil, false
	}

	return o.State, true
}

// SetState sets field value
func (o *MetadataWithStateFqdnZoneIdAllOf) SetState(v ProvisioningState) {

	o.State = &v

}

// HasState returns a boolean if a field has been set.
func (o *MetadataWithStateFqdnZoneIdAllOf) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// GetFqdn returns the Fqdn field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetFqdn() *string {
	if o == nil {
		return nil
	}

	return o.Fqdn

}

// GetFqdnOk returns a tuple with the Fqdn field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetFqdnOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Fqdn, true
}

// SetFqdn sets field value
func (o *MetadataWithStateFqdnZoneIdAllOf) SetFqdn(v string) {

	o.Fqdn = &v

}

// HasFqdn returns a boolean if a field has been set.
func (o *MetadataWithStateFqdnZoneIdAllOf) HasFqdn() bool {
	if o != nil && o.Fqdn != nil {
		return true
	}

	return false
}

// GetZoneId returns the ZoneId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetZoneId() *string {
	if o == nil {
		return nil
	}

	return o.ZoneId

}

// GetZoneIdOk returns a tuple with the ZoneId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithStateFqdnZoneIdAllOf) GetZoneIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ZoneId, true
}

// SetZoneId sets field value
func (o *MetadataWithStateFqdnZoneIdAllOf) SetZoneId(v string) {

	o.ZoneId = &v

}

// HasZoneId returns a boolean if a field has been set.
func (o *MetadataWithStateFqdnZoneIdAllOf) HasZoneId() bool {
	if o != nil && o.ZoneId != nil {
		return true
	}

	return false
}

func (o MetadataWithStateFqdnZoneIdAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.State != nil {
		toSerialize["state"] = o.State
	}

	if o.Fqdn != nil {
		toSerialize["fqdn"] = o.Fqdn
	}

	if o.ZoneId != nil {
		toSerialize["zoneId"] = o.ZoneId
	}

	return json.Marshal(toSerialize)
}

type NullableMetadataWithStateFqdnZoneIdAllOf struct {
	value *MetadataWithStateFqdnZoneIdAllOf
	isSet bool
}

func (v NullableMetadataWithStateFqdnZoneIdAllOf) Get() *MetadataWithStateFqdnZoneIdAllOf {
	return v.value
}

func (v *NullableMetadataWithStateFqdnZoneIdAllOf) Set(val *MetadataWithStateFqdnZoneIdAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataWithStateFqdnZoneIdAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataWithStateFqdnZoneIdAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataWithStateFqdnZoneIdAllOf(val *MetadataWithStateFqdnZoneIdAllOf) *NullableMetadataWithStateFqdnZoneIdAllOf {
	return &NullableMetadataWithStateFqdnZoneIdAllOf{value: val, isSet: true}
}

func (v NullableMetadataWithStateFqdnZoneIdAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataWithStateFqdnZoneIdAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
