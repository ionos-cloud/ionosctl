/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.16.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the DnssecKeyReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DnssecKeyReadList{}

// DnssecKeyReadList struct for DnssecKeyReadList
type DnssecKeyReadList struct {
	Id         *string                      `json:"id,omitempty"`
	Type       *string                      `json:"type,omitempty"`
	Href       *string                      `json:"href,omitempty"`
	Metadata   *DnssecKeyReadListMetadata   `json:"metadata,omitempty"`
	Properties *DnssecKeyReadListProperties `json:"properties,omitempty"`
}

// NewDnssecKeyReadList instantiates a new DnssecKeyReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnssecKeyReadList() *DnssecKeyReadList {
	this := DnssecKeyReadList{}

	return &this
}

// NewDnssecKeyReadListWithDefaults instantiates a new DnssecKeyReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnssecKeyReadListWithDefaults() *DnssecKeyReadList {
	this := DnssecKeyReadList{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *DnssecKeyReadList) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadList) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DnssecKeyReadList) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DnssecKeyReadList) SetId(v string) {
	o.Id = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *DnssecKeyReadList) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadList) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *DnssecKeyReadList) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *DnssecKeyReadList) SetType(v string) {
	o.Type = &v
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *DnssecKeyReadList) GetHref() string {
	if o == nil || IsNil(o.Href) {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadList) GetHrefOk() (*string, bool) {
	if o == nil || IsNil(o.Href) {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *DnssecKeyReadList) HasHref() bool {
	if o != nil && !IsNil(o.Href) {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *DnssecKeyReadList) SetHref(v string) {
	o.Href = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *DnssecKeyReadList) GetMetadata() DnssecKeyReadListMetadata {
	if o == nil || IsNil(o.Metadata) {
		var ret DnssecKeyReadListMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadList) GetMetadataOk() (*DnssecKeyReadListMetadata, bool) {
	if o == nil || IsNil(o.Metadata) {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *DnssecKeyReadList) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given DnssecKeyReadListMetadata and assigns it to the Metadata field.
func (o *DnssecKeyReadList) SetMetadata(v DnssecKeyReadListMetadata) {
	o.Metadata = &v
}

// GetProperties returns the Properties field value if set, zero value otherwise.
func (o *DnssecKeyReadList) GetProperties() DnssecKeyReadListProperties {
	if o == nil || IsNil(o.Properties) {
		var ret DnssecKeyReadListProperties
		return ret
	}
	return *o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadList) GetPropertiesOk() (*DnssecKeyReadListProperties, bool) {
	if o == nil || IsNil(o.Properties) {
		return nil, false
	}
	return o.Properties, true
}

// HasProperties returns a boolean if a field has been set.
func (o *DnssecKeyReadList) HasProperties() bool {
	if o != nil && !IsNil(o.Properties) {
		return true
	}

	return false
}

// SetProperties gets a reference to the given DnssecKeyReadListProperties and assigns it to the Properties field.
func (o *DnssecKeyReadList) SetProperties(v DnssecKeyReadListProperties) {
	o.Properties = &v
}

func (o DnssecKeyReadList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Href) {
		toSerialize["href"] = o.Href
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsNil(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableDnssecKeyReadList struct {
	value *DnssecKeyReadList
	isSet bool
}

func (v NullableDnssecKeyReadList) Get() *DnssecKeyReadList {
	return v.value
}

func (v *NullableDnssecKeyReadList) Set(val *DnssecKeyReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableDnssecKeyReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableDnssecKeyReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnssecKeyReadList(val *DnssecKeyReadList) *NullableDnssecKeyReadList {
	return &NullableDnssecKeyReadList{value: val, isSet: true}
}

func (v NullableDnssecKeyReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnssecKeyReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
