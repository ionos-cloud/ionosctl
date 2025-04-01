/*
 * CLOUD API
 *
 *  IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesNodePoolForPost struct for KubernetesNodePoolForPost
type KubernetesNodePoolForPost struct {
	// The resource's unique identifier.
	Id *string `json:"id,omitempty"`
	// The object type.
	Type *string `json:"type,omitempty"`
	// The URL to the object representation (absolute path).
	Href       *string                              `json:"href,omitempty"`
	Metadata   *DatacenterElementMetadata           `json:"metadata,omitempty"`
	Properties *KubernetesNodePoolPropertiesForPost `json:"properties"`
}

// NewKubernetesNodePoolForPost instantiates a new KubernetesNodePoolForPost object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKubernetesNodePoolForPost(properties KubernetesNodePoolPropertiesForPost) *KubernetesNodePoolForPost {
	this := KubernetesNodePoolForPost{}

	this.Properties = &properties

	return &this
}

// NewKubernetesNodePoolForPostWithDefaults instantiates a new KubernetesNodePoolForPost object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKubernetesNodePoolForPostWithDefaults() *KubernetesNodePoolForPost {
	this := KubernetesNodePoolForPost{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, nil is returned
func (o *KubernetesNodePoolForPost) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPost) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *KubernetesNodePoolForPost) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPost) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, nil is returned
func (o *KubernetesNodePoolForPost) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPost) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *KubernetesNodePoolForPost) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPost) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, nil is returned
func (o *KubernetesNodePoolForPost) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPost) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *KubernetesNodePoolForPost) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPost) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, nil is returned
func (o *KubernetesNodePoolForPost) GetMetadata() *DatacenterElementMetadata {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPost) GetMetadataOk() (*DatacenterElementMetadata, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *KubernetesNodePoolForPost) SetMetadata(v DatacenterElementMetadata) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPost) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, nil is returned
func (o *KubernetesNodePoolForPost) GetProperties() *KubernetesNodePoolPropertiesForPost {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPost) GetPropertiesOk() (*KubernetesNodePoolPropertiesForPost, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *KubernetesNodePoolForPost) SetProperties(v KubernetesNodePoolPropertiesForPost) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPost) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o KubernetesNodePoolForPost) MarshalJSON() ([]byte, error) {
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

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableKubernetesNodePoolForPost struct {
	value *KubernetesNodePoolForPost
	isSet bool
}

func (v NullableKubernetesNodePoolForPost) Get() *KubernetesNodePoolForPost {
	return v.value
}

func (v *NullableKubernetesNodePoolForPost) Set(val *KubernetesNodePoolForPost) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodePoolForPost) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodePoolForPost) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodePoolForPost(val *KubernetesNodePoolForPost) *NullableKubernetesNodePoolForPost {
	return &NullableKubernetesNodePoolForPost{value: val, isSet: true}
}

func (v NullableKubernetesNodePoolForPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodePoolForPost) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
