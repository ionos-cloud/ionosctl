/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ApplicationLoadBalancerPut struct for ApplicationLoadBalancerPut
type ApplicationLoadBalancerPut struct {
	// The resource's unique identifier
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	Properties *ApplicationLoadBalancerProperties `json:"properties"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerPut) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerPut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *ApplicationLoadBalancerPut) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerPut) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for Type will be returned
func (o *ApplicationLoadBalancerPut) GetType() *Type {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerPut) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *ApplicationLoadBalancerPut) SetType(v Type) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerPut) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerPut) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerPut) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *ApplicationLoadBalancerPut) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerPut) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for ApplicationLoadBalancerProperties will be returned
func (o *ApplicationLoadBalancerPut) GetProperties() *ApplicationLoadBalancerProperties {
	if o == nil {
		return nil
	}


	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerPut) GetPropertiesOk() (*ApplicationLoadBalancerProperties, bool) {
	if o == nil {
		return nil, false
	}


	return o.Properties, true
}

// SetProperties sets field value
func (o *ApplicationLoadBalancerPut) SetProperties(v ApplicationLoadBalancerProperties) {


	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerPut) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}


func (o ApplicationLoadBalancerPut) MarshalJSON() ([]byte, error) {
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

type NullableApplicationLoadBalancerPut struct {
	value *ApplicationLoadBalancerPut
	isSet bool
}

func (v NullableApplicationLoadBalancerPut) Get() *ApplicationLoadBalancerPut {
	return v.value
}

func (v *NullableApplicationLoadBalancerPut) Set(val *ApplicationLoadBalancerPut) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationLoadBalancerPut) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationLoadBalancerPut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationLoadBalancerPut(val *ApplicationLoadBalancerPut) *NullableApplicationLoadBalancerPut {
	return &NullableApplicationLoadBalancerPut{value: val, isSet: true}
}

func (v NullableApplicationLoadBalancerPut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationLoadBalancerPut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


