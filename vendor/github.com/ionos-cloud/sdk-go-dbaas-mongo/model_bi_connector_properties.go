/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.  MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// BiConnectorProperties struct for BiConnectorProperties
type BiConnectorProperties struct {
	// The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands to aid in data analysis.
	Enabled *bool `json:"enabled,omitempty"`
	// The host where this new BI Connector is installed.
	Host *string `json:"host,omitempty"`
	// Port number used when connecting to this new BI Connector.
	Port *string `json:"port,omitempty"`
}

// NewBiConnectorProperties instantiates a new BiConnectorProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBiConnectorProperties() *BiConnectorProperties {
	this := BiConnectorProperties{}

	var enabled bool = false
	this.Enabled = &enabled

	return &this
}

// NewBiConnectorPropertiesWithDefaults instantiates a new BiConnectorProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBiConnectorPropertiesWithDefaults() *BiConnectorProperties {
	this := BiConnectorProperties{}
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// GetEnabled returns the Enabled field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *BiConnectorProperties) GetEnabled() *bool {
	if o == nil {
		return nil
	}

	return o.Enabled

}

// GetEnabledOk returns a tuple with the Enabled field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BiConnectorProperties) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}

	return o.Enabled, true
}

// SetEnabled sets field value
func (o *BiConnectorProperties) SetEnabled(v bool) {

	o.Enabled = &v

}

// HasEnabled returns a boolean if a field has been set.
func (o *BiConnectorProperties) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// GetHost returns the Host field value
// If the value is explicit nil, the zero value for string will be returned
func (o *BiConnectorProperties) GetHost() *string {
	if o == nil {
		return nil
	}

	return o.Host

}

// GetHostOk returns a tuple with the Host field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BiConnectorProperties) GetHostOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Host, true
}

// SetHost sets field value
func (o *BiConnectorProperties) SetHost(v string) {

	o.Host = &v

}

// HasHost returns a boolean if a field has been set.
func (o *BiConnectorProperties) HasHost() bool {
	if o != nil && o.Host != nil {
		return true
	}

	return false
}

// GetPort returns the Port field value
// If the value is explicit nil, the zero value for string will be returned
func (o *BiConnectorProperties) GetPort() *string {
	if o == nil {
		return nil
	}

	return o.Port

}

// GetPortOk returns a tuple with the Port field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BiConnectorProperties) GetPortOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Port, true
}

// SetPort sets field value
func (o *BiConnectorProperties) SetPort(v string) {

	o.Port = &v

}

// HasPort returns a boolean if a field has been set.
func (o *BiConnectorProperties) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

func (o BiConnectorProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.Host != nil {
		toSerialize["host"] = o.Host
	}
	if o.Port != nil {
		toSerialize["port"] = o.Port
	}
	return json.Marshal(toSerialize)
}

type NullableBiConnectorProperties struct {
	value *BiConnectorProperties
	isSet bool
}

func (v NullableBiConnectorProperties) Get() *BiConnectorProperties {
	return v.value
}

func (v *NullableBiConnectorProperties) Set(val *BiConnectorProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableBiConnectorProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableBiConnectorProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBiConnectorProperties(val *BiConnectorProperties) *NullableBiConnectorProperties {
	return &NullableBiConnectorProperties{value: val, isSet: true}
}

func (v NullableBiConnectorProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBiConnectorProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}