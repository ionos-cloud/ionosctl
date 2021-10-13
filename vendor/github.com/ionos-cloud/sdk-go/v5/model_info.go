/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Info struct for Info
type Info struct {
	// API entry point
	Href *string `json:"href,omitempty"`
	// Name of the API
	Name *string `json:"name,omitempty"`
	// Version of the API
	Version *string `json:"version,omitempty"`
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Info) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Info) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *Info) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *Info) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}


// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Info) GetName() *string {
	if o == nil {
		return nil
	}


	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Info) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Name, true
}

// SetName sets field value
func (o *Info) SetName(v string) {


	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *Info) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}


// GetVersion returns the Version field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Info) GetVersion() *string {
	if o == nil {
		return nil
	}


	return o.Version

}

// GetVersionOk returns a tuple with the Version field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Info) GetVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Version, true
}

// SetVersion sets field value
func (o *Info) SetVersion(v string) {


	o.Version = &v

}

// HasVersion returns a boolean if a field has been set.
func (o *Info) HasVersion() bool {
	if o != nil && o.Version != nil {
		return true
	}

	return false
}

func (o Info) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Version != nil {
		toSerialize["version"] = o.Version
	}
	return json.Marshal(toSerialize)
}

type NullableInfo struct {
	value *Info
	isSet bool
}

func (v NullableInfo) Get() *Info {
	return v.value
}

func (v *NullableInfo) Set(val *Info) {
	v.value = val
	v.isSet = true
}

func (v NullableInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInfo(val *Info) *NullableInfo {
	return &NullableInfo{value: val, isSet: true}
}

func (v NullableInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


