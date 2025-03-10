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

// checks if the IPSecTunnelCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecTunnelCreate{}

// IPSecTunnelCreate struct for IPSecTunnelCreate
type IPSecTunnelCreate struct {
	// Metadata
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Properties IPSecTunnel            `json:"properties"`
}

// NewIPSecTunnelCreate instantiates a new IPSecTunnelCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecTunnelCreate(properties IPSecTunnel) *IPSecTunnelCreate {
	this := IPSecTunnelCreate{}

	this.Properties = properties

	return &this
}

// NewIPSecTunnelCreateWithDefaults instantiates a new IPSecTunnelCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecTunnelCreateWithDefaults() *IPSecTunnelCreate {
	this := IPSecTunnelCreate{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *IPSecTunnelCreate) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecTunnelCreate) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *IPSecTunnelCreate) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *IPSecTunnelCreate) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *IPSecTunnelCreate) GetProperties() IPSecTunnel {
	if o == nil {
		var ret IPSecTunnel
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnelCreate) GetPropertiesOk() (*IPSecTunnel, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *IPSecTunnelCreate) SetProperties(v IPSecTunnel) {
	o.Properties = v
}

func (o IPSecTunnelCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableIPSecTunnelCreate struct {
	value *IPSecTunnelCreate
	isSet bool
}

func (v NullableIPSecTunnelCreate) Get() *IPSecTunnelCreate {
	return v.value
}

func (v *NullableIPSecTunnelCreate) Set(val *IPSecTunnelCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecTunnelCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecTunnelCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecTunnelCreate(val *IPSecTunnelCreate) *NullableIPSecTunnelCreate {
	return &NullableIPSecTunnelCreate{value: val, isSet: true}
}

func (v NullableIPSecTunnelCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecTunnelCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
