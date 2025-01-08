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

// checks if the RecordReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RecordReadList{}

// RecordReadList struct for RecordReadList
type RecordReadList struct {
	// The resource's unique identifier.
	Id    string       `json:"id"`
	Type  string       `json:"type"`
	Href  string       `json:"href"`
	Items []RecordRead `json:"items"`
	// Pagination offset.
	Offset float32 `json:"offset"`
	// Pagination limit.
	Limit float32 `json:"limit"`
	Links Links   `json:"_links"`
}

// NewRecordReadList instantiates a new RecordReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRecordReadList(id string, type_ string, href string, items []RecordRead, offset float32, limit float32, links Links) *RecordReadList {
	this := RecordReadList{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Items = items
	this.Offset = offset
	this.Limit = limit
	this.Links = links

	return &this
}

// NewRecordReadListWithDefaults instantiates a new RecordReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRecordReadListWithDefaults() *RecordReadList {
	this := RecordReadList{}
	return &this
}

// GetId returns the Id field value
func (o *RecordReadList) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *RecordReadList) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *RecordReadList) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *RecordReadList) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *RecordReadList) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *RecordReadList) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value
func (o *RecordReadList) GetItems() []RecordRead {
	if o == nil {
		var ret []RecordRead
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetItemsOk() ([]RecordRead, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *RecordReadList) SetItems(v []RecordRead) {
	o.Items = v
}

// GetOffset returns the Offset field value
func (o *RecordReadList) GetOffset() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Offset
}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetOffsetOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Offset, true
}

// SetOffset sets field value
func (o *RecordReadList) SetOffset(v float32) {
	o.Offset = v
}

// GetLimit returns the Limit field value
func (o *RecordReadList) GetLimit() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Limit
}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetLimitOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Limit, true
}

// SetLimit sets field value
func (o *RecordReadList) SetLimit(v float32) {
	o.Limit = v
}

// GetLinks returns the Links field value
func (o *RecordReadList) GetLinks() Links {
	if o == nil {
		var ret Links
		return ret
	}

	return o.Links
}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
func (o *RecordReadList) GetLinksOk() (*Links, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Links, true
}

// SetLinks sets field value
func (o *RecordReadList) SetLinks(v Links) {
	o.Links = v
}

func (o RecordReadList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	toSerialize["items"] = o.Items
	toSerialize["offset"] = o.Offset
	toSerialize["limit"] = o.Limit
	toSerialize["_links"] = o.Links
	return toSerialize, nil
}

type NullableRecordReadList struct {
	value *RecordReadList
	isSet bool
}

func (v NullableRecordReadList) Get() *RecordReadList {
	return v.value
}

func (v *NullableRecordReadList) Set(val *RecordReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableRecordReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableRecordReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRecordReadList(val *RecordReadList) *NullableRecordReadList {
	return &NullableRecordReadList{value: val, isSet: true}
}

func (v NullableRecordReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRecordReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
