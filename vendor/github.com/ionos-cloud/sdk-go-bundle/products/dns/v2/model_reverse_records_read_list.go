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

// checks if the ReverseRecordsReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReverseRecordsReadList{}

// ReverseRecordsReadList struct for ReverseRecordsReadList
type ReverseRecordsReadList struct {
	// ID (UUID) created to identify this action.
	Id    string              `json:"id"`
	Type  string              `json:"type"`
	Href  string              `json:"href"`
	Items []ReverseRecordRead `json:"items"`
	// Pagination offset.
	Offset float32 `json:"offset"`
	// Pagination limit.
	Limit float32 `json:"limit"`
	Links Links   `json:"_links"`
}

// NewReverseRecordsReadList instantiates a new ReverseRecordsReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReverseRecordsReadList(id string, type_ string, href string, items []ReverseRecordRead, offset float32, limit float32, links Links) *ReverseRecordsReadList {
	this := ReverseRecordsReadList{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Items = items
	this.Offset = offset
	this.Limit = limit
	this.Links = links

	return &this
}

// NewReverseRecordsReadListWithDefaults instantiates a new ReverseRecordsReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReverseRecordsReadListWithDefaults() *ReverseRecordsReadList {
	this := ReverseRecordsReadList{}
	return &this
}

// GetId returns the Id field value
func (o *ReverseRecordsReadList) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *ReverseRecordsReadList) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *ReverseRecordsReadList) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *ReverseRecordsReadList) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *ReverseRecordsReadList) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *ReverseRecordsReadList) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value
func (o *ReverseRecordsReadList) GetItems() []ReverseRecordRead {
	if o == nil {
		var ret []ReverseRecordRead
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetItemsOk() ([]ReverseRecordRead, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *ReverseRecordsReadList) SetItems(v []ReverseRecordRead) {
	o.Items = v
}

// GetOffset returns the Offset field value
func (o *ReverseRecordsReadList) GetOffset() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Offset
}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetOffsetOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Offset, true
}

// SetOffset sets field value
func (o *ReverseRecordsReadList) SetOffset(v float32) {
	o.Offset = v
}

// GetLimit returns the Limit field value
func (o *ReverseRecordsReadList) GetLimit() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Limit
}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetLimitOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Limit, true
}

// SetLimit sets field value
func (o *ReverseRecordsReadList) SetLimit(v float32) {
	o.Limit = v
}

// GetLinks returns the Links field value
func (o *ReverseRecordsReadList) GetLinks() Links {
	if o == nil {
		var ret Links
		return ret
	}

	return o.Links
}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
func (o *ReverseRecordsReadList) GetLinksOk() (*Links, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Links, true
}

// SetLinks sets field value
func (o *ReverseRecordsReadList) SetLinks(v Links) {
	o.Links = v
}

func (o ReverseRecordsReadList) ToMap() (map[string]interface{}, error) {
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

type NullableReverseRecordsReadList struct {
	value *ReverseRecordsReadList
	isSet bool
}

func (v NullableReverseRecordsReadList) Get() *ReverseRecordsReadList {
	return v.value
}

func (v *NullableReverseRecordsReadList) Set(val *ReverseRecordsReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableReverseRecordsReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableReverseRecordsReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReverseRecordsReadList(val *ReverseRecordsReadList) *NullableReverseRecordsReadList {
	return &NullableReverseRecordsReadList{value: val, isSet: true}
}

func (v NullableReverseRecordsReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReverseRecordsReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
