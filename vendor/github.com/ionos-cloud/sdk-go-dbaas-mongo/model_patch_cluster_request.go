/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.   MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// PatchClusterRequest Request payload to change a cluster.
type PatchClusterRequest struct {
	Metadata   *Metadata               `json:"metadata,omitempty"`
	Properties *PatchClusterProperties `json:"properties,omitempty"`
}

// NewPatchClusterRequest instantiates a new PatchClusterRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchClusterRequest() *PatchClusterRequest {
	this := PatchClusterRequest{}

	return &this
}

// NewPatchClusterRequestWithDefaults instantiates a new PatchClusterRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchClusterRequestWithDefaults() *PatchClusterRequest {
	this := PatchClusterRequest{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for Metadata will be returned
func (o *PatchClusterRequest) GetMetadata() *Metadata {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchClusterRequest) GetMetadataOk() (*Metadata, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *PatchClusterRequest) SetMetadata(v Metadata) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *PatchClusterRequest) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for PatchClusterProperties will be returned
func (o *PatchClusterRequest) GetProperties() *PatchClusterProperties {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchClusterRequest) GetPropertiesOk() (*PatchClusterProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *PatchClusterRequest) SetProperties(v PatchClusterProperties) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *PatchClusterRequest) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o PatchClusterRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullablePatchClusterRequest struct {
	value *PatchClusterRequest
	isSet bool
}

func (v NullablePatchClusterRequest) Get() *PatchClusterRequest {
	return v.value
}

func (v *NullablePatchClusterRequest) Set(val *PatchClusterRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchClusterRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchClusterRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchClusterRequest(val *PatchClusterRequest) *NullablePatchClusterRequest {
	return &NullablePatchClusterRequest{value: val, isSet: true}
}

func (v NullablePatchClusterRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchClusterRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
