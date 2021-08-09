/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.2
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// LabelResource struct for LabelResource
type LabelResource struct {
	// Label on a resource is identified using label key.
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *string `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	Metadata *NoStateMetaData `json:"metadata,omitempty"`
	Properties *LabelResourceProperties `json:"properties"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *LabelResource) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LabelResource) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *LabelResource) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *LabelResource) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *LabelResource) GetType() *string {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LabelResource) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *LabelResource) SetType(v string) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *LabelResource) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *LabelResource) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LabelResource) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *LabelResource) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *LabelResource) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for NoStateMetaData will be returned
func (o *LabelResource) GetMetadata() *NoStateMetaData {
	if o == nil {
		return nil
	}


	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LabelResource) GetMetadataOk() (*NoStateMetaData, bool) {
	if o == nil {
		return nil, false
	}


	return o.Metadata, true
}

// SetMetadata sets field value
func (o *LabelResource) SetMetadata(v NoStateMetaData) {


	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *LabelResource) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}



// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for LabelResourceProperties will be returned
func (o *LabelResource) GetProperties() *LabelResourceProperties {
	if o == nil {
		return nil
	}


	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LabelResource) GetPropertiesOk() (*LabelResourceProperties, bool) {
	if o == nil {
		return nil, false
	}


	return o.Properties, true
}

// SetProperties sets field value
func (o *LabelResource) SetProperties(v LabelResourceProperties) {


	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *LabelResource) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}


func (o LabelResource) MarshalJSON() ([]byte, error) {
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

type NullableLabelResource struct {
	value *LabelResource
	isSet bool
}

func (v NullableLabelResource) Get() *LabelResource {
	return v.value
}

func (v *NullableLabelResource) Set(val *LabelResource) {
	v.value = val
	v.isSet = true
}

func (v NullableLabelResource) IsSet() bool {
	return v.isSet
}

func (v *NullableLabelResource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLabelResource(val *LabelResource) *NullableLabelResource {
	return &NullableLabelResource{value: val, isSet: true}
}

func (v NullableLabelResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLabelResource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


