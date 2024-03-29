/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ServerProperties struct for ServerProperties
type ServerProperties struct {
	DatacenterServer *DatacenterServer `json:"datacenterServer"`
	Name             *string           `json:"name,omitempty"`
}

// NewServerProperties instantiates a new ServerProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerProperties(datacenterServer DatacenterServer) *ServerProperties {
	this := ServerProperties{}

	this.DatacenterServer = &datacenterServer

	return &this
}

// NewServerPropertiesWithDefaults instantiates a new ServerProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerPropertiesWithDefaults() *ServerProperties {
	this := ServerProperties{}
	return &this
}

// GetDatacenterServer returns the DatacenterServer field value
// If the value is explicit nil, the zero value for DatacenterServer will be returned
func (o *ServerProperties) GetDatacenterServer() *DatacenterServer {
	if o == nil {
		return nil
	}

	return o.DatacenterServer

}

// GetDatacenterServerOk returns a tuple with the DatacenterServer field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ServerProperties) GetDatacenterServerOk() (*DatacenterServer, bool) {
	if o == nil {
		return nil, false
	}

	return o.DatacenterServer, true
}

// SetDatacenterServer sets field value
func (o *ServerProperties) SetDatacenterServer(v DatacenterServer) {

	o.DatacenterServer = &v

}

// HasDatacenterServer returns a boolean if a field has been set.
func (o *ServerProperties) HasDatacenterServer() bool {
	if o != nil && o.DatacenterServer != nil {
		return true
	}

	return false
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ServerProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ServerProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *ServerProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *ServerProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

func (o ServerProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DatacenterServer != nil {
		toSerialize["datacenterServer"] = o.DatacenterServer
	}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	return json.Marshal(toSerialize)
}

type NullableServerProperties struct {
	value *ServerProperties
	isSet bool
}

func (v NullableServerProperties) Get() *ServerProperties {
	return v.value
}

func (v *NullableServerProperties) Set(val *ServerProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableServerProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableServerProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerProperties(val *ServerProperties) *NullableServerProperties {
	return &NullableServerProperties{value: val, isSet: true}
}

func (v NullableServerProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
