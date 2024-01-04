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
	"time"
)

// MetadataBasic The resource metadata.
type MetadataBasic struct {
	// The date the resource was created.
	CreatedDate *IonosTime `json:"createdDate"`
	// The resource etag.
	Etag *string `json:"etag"`
	// The date the resource was last modified.
	LastModifiedDate *IonosTime `json:"lastModifiedDate"`
	// The resource state.
	State *string `json:"state"`
}

// NewMetadataBasic instantiates a new MetadataBasic object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataBasic(createdDate time.Time, etag string, lastModifiedDate time.Time, state string) *MetadataBasic {
	this := MetadataBasic{}

	this.CreatedDate = &IonosTime{createdDate}
	this.Etag = &etag
	this.LastModifiedDate = &IonosTime{lastModifiedDate}
	this.State = &state

	return &this
}

// NewMetadataBasicWithDefaults instantiates a new MetadataBasic object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataBasicWithDefaults() *MetadataBasic {
	this := MetadataBasic{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *MetadataBasic) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time

}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataBasic) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *MetadataBasic) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}

}

// HasCreatedDate returns a boolean if a field has been set.
func (o *MetadataBasic) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}

// GetEtag returns the Etag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataBasic) GetEtag() *string {
	if o == nil {
		return nil
	}

	return o.Etag

}

// GetEtagOk returns a tuple with the Etag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataBasic) GetEtagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Etag, true
}

// SetEtag sets field value
func (o *MetadataBasic) SetEtag(v string) {

	o.Etag = &v

}

// HasEtag returns a boolean if a field has been set.
func (o *MetadataBasic) HasEtag() bool {
	if o != nil && o.Etag != nil {
		return true
	}

	return false
}

// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *MetadataBasic) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time

}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataBasic) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *MetadataBasic) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}

}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *MetadataBasic) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}

// GetState returns the State field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataBasic) GetState() *string {
	if o == nil {
		return nil
	}

	return o.State

}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataBasic) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.State, true
}

// SetState sets field value
func (o *MetadataBasic) SetState(v string) {

	o.State = &v

}

// HasState returns a boolean if a field has been set.
func (o *MetadataBasic) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

func (o MetadataBasic) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}

	if o.Etag != nil {
		toSerialize["etag"] = o.Etag
	}

	if o.LastModifiedDate != nil {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}

	if o.State != nil {
		toSerialize["state"] = o.State
	}

	return json.Marshal(toSerialize)
}

type NullableMetadataBasic struct {
	value *MetadataBasic
	isSet bool
}

func (v NullableMetadataBasic) Get() *MetadataBasic {
	return v.value
}

func (v *NullableMetadataBasic) Set(val *MetadataBasic) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataBasic) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataBasic) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataBasic(val *MetadataBasic) *NullableMetadataBasic {
	return &NullableMetadataBasic{value: val, isSet: true}
}

func (v NullableMetadataBasic) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataBasic) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
