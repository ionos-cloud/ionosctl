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

// APIVersion struct for APIVersion
type APIVersion struct {
	Name       *string `json:"name,omitempty"`
	SwaggerUrl *string `json:"swaggerUrl,omitempty"`
}

// NewAPIVersion instantiates a new APIVersion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAPIVersion() *APIVersion {
	this := APIVersion{}

	return &this
}

// NewAPIVersionWithDefaults instantiates a new APIVersion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAPIVersionWithDefaults() *APIVersion {
	this := APIVersion{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *APIVersion) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *APIVersion) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *APIVersion) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *APIVersion) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetSwaggerUrl returns the SwaggerUrl field value
// If the value is explicit nil, the zero value for string will be returned
func (o *APIVersion) GetSwaggerUrl() *string {
	if o == nil {
		return nil
	}

	return o.SwaggerUrl

}

// GetSwaggerUrlOk returns a tuple with the SwaggerUrl field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *APIVersion) GetSwaggerUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.SwaggerUrl, true
}

// SetSwaggerUrl sets field value
func (o *APIVersion) SetSwaggerUrl(v string) {

	o.SwaggerUrl = &v

}

// HasSwaggerUrl returns a boolean if a field has been set.
func (o *APIVersion) HasSwaggerUrl() bool {
	if o != nil && o.SwaggerUrl != nil {
		return true
	}

	return false
}

func (o APIVersion) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.SwaggerUrl != nil {
		toSerialize["swaggerUrl"] = o.SwaggerUrl
	}
	return json.Marshal(toSerialize)
}

type NullableAPIVersion struct {
	value *APIVersion
	isSet bool
}

func (v NullableAPIVersion) Get() *APIVersion {
	return v.value
}

func (v *NullableAPIVersion) Set(val *APIVersion) {
	v.value = val
	v.isSet = true
}

func (v NullableAPIVersion) IsSet() bool {
	return v.isSet
}

func (v *NullableAPIVersion) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAPIVersion(val *APIVersion) *NullableAPIVersion {
	return &NullableAPIVersion{value: val, isSet: true}
}

func (v NullableAPIVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAPIVersion) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
