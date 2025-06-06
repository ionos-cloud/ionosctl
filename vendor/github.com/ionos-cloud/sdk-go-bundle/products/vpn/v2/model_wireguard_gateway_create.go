/*
 * IONOS Cloud VPN Gateway API
 *
 * The Managed VPN Gateway service provides secure and scalable connectivity, enabling encrypted communication between your IONOS cloud resources in a VDC and remote networks (on-premises, multi-cloud, private LANs in other VDCs etc).
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vpn

import (
	"encoding/json"
)

// checks if the WireguardGatewayCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardGatewayCreate{}

// WireguardGatewayCreate struct for WireguardGatewayCreate
type WireguardGatewayCreate struct {
	// Metadata
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Properties WireguardGateway       `json:"properties"`
}

// NewWireguardGatewayCreate instantiates a new WireguardGatewayCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardGatewayCreate(properties WireguardGateway) *WireguardGatewayCreate {
	this := WireguardGatewayCreate{}

	this.Properties = properties

	return &this
}

// NewWireguardGatewayCreateWithDefaults instantiates a new WireguardGatewayCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardGatewayCreateWithDefaults() *WireguardGatewayCreate {
	this := WireguardGatewayCreate{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *WireguardGatewayCreate) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardGatewayCreate) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *WireguardGatewayCreate) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *WireguardGatewayCreate) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *WireguardGatewayCreate) GetProperties() WireguardGateway {
	if o == nil {
		var ret WireguardGateway
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *WireguardGatewayCreate) GetPropertiesOk() (*WireguardGateway, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *WireguardGatewayCreate) SetProperties(v WireguardGateway) {
	o.Properties = v
}

func (o WireguardGatewayCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableWireguardGatewayCreate struct {
	value *WireguardGatewayCreate
	isSet bool
}

func (v NullableWireguardGatewayCreate) Get() *WireguardGatewayCreate {
	return v.value
}

func (v *NullableWireguardGatewayCreate) Set(val *WireguardGatewayCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardGatewayCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardGatewayCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardGatewayCreate(val *WireguardGatewayCreate) *NullableWireguardGatewayCreate {
	return &NullableWireguardGatewayCreate{value: val, isSet: true}
}

func (v NullableWireguardGatewayCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardGatewayCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
