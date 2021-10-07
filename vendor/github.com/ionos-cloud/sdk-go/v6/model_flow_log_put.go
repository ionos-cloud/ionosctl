/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// FlowLogPut struct for FlowLogPut
type FlowLogPut struct {
	// The resource's unique identifier
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	Properties *FlowLogProperties `json:"properties"`
}


// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *FlowLogPut) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FlowLogPut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *FlowLogPut) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *FlowLogPut) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for Type will be returned
func (o *FlowLogPut) GetType() *Type {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FlowLogPut) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *FlowLogPut) SetType(v Type) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *FlowLogPut) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *FlowLogPut) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FlowLogPut) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *FlowLogPut) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *FlowLogPut) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for FlowLogProperties will be returned
func (o *FlowLogPut) GetProperties() *FlowLogProperties {
	if o == nil {
		return nil
	}


	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *FlowLogPut) GetPropertiesOk() (*FlowLogProperties, bool) {
	if o == nil {
		return nil, false
	}


	return o.Properties, true
}

// SetProperties sets field value
func (o *FlowLogPut) SetProperties(v FlowLogProperties) {


	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *FlowLogPut) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o FlowLogPut) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}
	return json.Marshal(toSerialize)
}
type NullableFlowLogPut struct {
	value *FlowLogPut
	isSet bool
}

func (v NullableFlowLogPut) Get() *FlowLogPut {
	return v.value
}

func (v *NullableFlowLogPut) Set(val *FlowLogPut) {
	v.value = val
	v.isSet = true
}

func (v NullableFlowLogPut) IsSet() bool {
	return v.isSet
}

func (v *NullableFlowLogPut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFlowLogPut(val *FlowLogPut) *NullableFlowLogPut {
	return &NullableFlowLogPut{value: val, isSet: true}
}

func (v NullableFlowLogPut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFlowLogPut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


