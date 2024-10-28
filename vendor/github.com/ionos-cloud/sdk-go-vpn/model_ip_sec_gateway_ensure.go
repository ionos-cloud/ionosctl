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

// IPSecGatewayEnsure struct for IPSecGatewayEnsure
type IPSecGatewayEnsure struct {
	// The ID (UUID) of the IPSecGateway.
	Id *string `json:"id"`
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *IPSecGateway           `json:"properties"`
}

// NewIPSecGatewayEnsure instantiates a new IPSecGatewayEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecGatewayEnsure(id string, properties IPSecGateway) *IPSecGatewayEnsure {
	this := IPSecGatewayEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewIPSecGatewayEnsureWithDefaults instantiates a new IPSecGatewayEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecGatewayEnsureWithDefaults() *IPSecGatewayEnsure {
	this := IPSecGatewayEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *IPSecGatewayEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IPSecGatewayEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *IPSecGatewayEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *IPSecGatewayEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *IPSecGatewayEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IPSecGatewayEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *IPSecGatewayEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *IPSecGatewayEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for IPSecGateway will be returned
func (o *IPSecGatewayEnsure) GetProperties() *IPSecGateway {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IPSecGatewayEnsure) GetPropertiesOk() (*IPSecGateway, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *IPSecGatewayEnsure) SetProperties(v IPSecGateway) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *IPSecGatewayEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o IPSecGatewayEnsure) MarshalJSON() ([]byte, error) {
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

type NullableIPSecGatewayEnsure struct {
	value *IPSecGatewayEnsure
	isSet bool
}

func (v NullableIPSecGatewayEnsure) Get() *IPSecGatewayEnsure {
	return v.value
}

func (v *NullableIPSecGatewayEnsure) Set(val *IPSecGatewayEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecGatewayEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecGatewayEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecGatewayEnsure(val *IPSecGatewayEnsure) *NullableIPSecGatewayEnsure {
	return &NullableIPSecGatewayEnsure{value: val, isSet: true}
}

func (v NullableIPSecGatewayEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecGatewayEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
