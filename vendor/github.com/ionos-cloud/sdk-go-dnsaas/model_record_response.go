/*
 * IONOS Cloud - DNS as a Service API
 *
 * DNS API Specification
 *
 * API version: 0.1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RecordResponse struct for RecordResponse
type RecordResponse struct {
	Id         *string           `json:"id,omitempty"`
	Metadata   *RecordMetadata   `json:"metadata,omitempty"`
	Properties *RecordProperties `json:"properties,omitempty"`
}

// NewRecordResponse instantiates a new RecordResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRecordResponse() *RecordResponse {
	this := RecordResponse{}

	return &this
}

// NewRecordResponseWithDefaults instantiates a new RecordResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRecordResponseWithDefaults() *RecordResponse {
	this := RecordResponse{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RecordResponse) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RecordResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *RecordResponse) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *RecordResponse) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for RecordMetadata will be returned
func (o *RecordResponse) GetMetadata() *RecordMetadata {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RecordResponse) GetMetadataOk() (*RecordMetadata, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *RecordResponse) SetMetadata(v RecordMetadata) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *RecordResponse) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for RecordProperties will be returned
func (o *RecordResponse) GetProperties() *RecordProperties {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RecordResponse) GetPropertiesOk() (*RecordProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *RecordResponse) SetProperties(v RecordProperties) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *RecordResponse) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o RecordResponse) MarshalJSON() ([]byte, error) {
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

type NullableRecordResponse struct {
	value *RecordResponse
	isSet bool
}

func (v NullableRecordResponse) Get() *RecordResponse {
	return v.value
}

func (v *NullableRecordResponse) Set(val *RecordResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableRecordResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableRecordResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRecordResponse(val *RecordResponse) *NullableRecordResponse {
	return &NullableRecordResponse{value: val, isSet: true}
}

func (v NullableRecordResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRecordResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}