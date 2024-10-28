/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// WireguardPeerEnsure struct for WireguardPeerEnsure
type WireguardPeerEnsure struct {
	// The ID (UUID) of the WireguardPeer.
	Id *string `json:"id"`
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *WireguardPeer          `json:"properties"`
}

// NewWireguardPeerEnsure instantiates a new WireguardPeerEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardPeerEnsure(id string, properties WireguardPeer) *WireguardPeerEnsure {
	this := WireguardPeerEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewWireguardPeerEnsureWithDefaults instantiates a new WireguardPeerEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardPeerEnsureWithDefaults() *WireguardPeerEnsure {
	this := WireguardPeerEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *WireguardPeerEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WireguardPeerEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *WireguardPeerEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *WireguardPeerEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *WireguardPeerEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WireguardPeerEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *WireguardPeerEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *WireguardPeerEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for WireguardPeer will be returned
func (o *WireguardPeerEnsure) GetProperties() *WireguardPeer {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WireguardPeerEnsure) GetPropertiesOk() (*WireguardPeer, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *WireguardPeerEnsure) SetProperties(v WireguardPeer) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *WireguardPeerEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o WireguardPeerEnsure) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableWireguardPeerEnsure struct {
	value *WireguardPeerEnsure
	isSet bool
}

func (v NullableWireguardPeerEnsure) Get() *WireguardPeerEnsure {
	return v.value
}

func (v *NullableWireguardPeerEnsure) Set(val *WireguardPeerEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardPeerEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardPeerEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardPeerEnsure(val *WireguardPeerEnsure) *NullableWireguardPeerEnsure {
	return &NullableWireguardPeerEnsure{value: val, isSet: true}
}

func (v NullableWireguardPeerEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardPeerEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
