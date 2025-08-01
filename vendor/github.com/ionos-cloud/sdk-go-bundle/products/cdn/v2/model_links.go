/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the Links type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Links{}

// Links URLs to navigate the different pages. As of now we always only return a single page.
type Links struct {
	// URL (with offset and limit parameters) of the previous page; only present if offset is greater than 0.
	Prev *string `json:"prev,omitempty"`
	// URL (with offset and limit parameters) of the current page.
	Self *string `json:"self,omitempty"`
	// URL (with offset and limit parameters) of the next page; only present if offset + limit is less than the total number of elements.
	Next *string `json:"next,omitempty"`
}

// NewLinks instantiates a new Links object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinks() *Links {
	this := Links{}

	return &this
}

// NewLinksWithDefaults instantiates a new Links object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinksWithDefaults() *Links {
	this := Links{}
	return &this
}

// GetPrev returns the Prev field value if set, zero value otherwise.
func (o *Links) GetPrev() string {
	if o == nil || IsNil(o.Prev) {
		var ret string
		return ret
	}
	return *o.Prev
}

// GetPrevOk returns a tuple with the Prev field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Links) GetPrevOk() (*string, bool) {
	if o == nil || IsNil(o.Prev) {
		return nil, false
	}
	return o.Prev, true
}

// HasPrev returns a boolean if a field has been set.
func (o *Links) HasPrev() bool {
	if o != nil && !IsNil(o.Prev) {
		return true
	}

	return false
}

// SetPrev gets a reference to the given string and assigns it to the Prev field.
func (o *Links) SetPrev(v string) {
	o.Prev = &v
}

// GetSelf returns the Self field value if set, zero value otherwise.
func (o *Links) GetSelf() string {
	if o == nil || IsNil(o.Self) {
		var ret string
		return ret
	}
	return *o.Self
}

// GetSelfOk returns a tuple with the Self field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Links) GetSelfOk() (*string, bool) {
	if o == nil || IsNil(o.Self) {
		return nil, false
	}
	return o.Self, true
}

// HasSelf returns a boolean if a field has been set.
func (o *Links) HasSelf() bool {
	if o != nil && !IsNil(o.Self) {
		return true
	}

	return false
}

// SetSelf gets a reference to the given string and assigns it to the Self field.
func (o *Links) SetSelf(v string) {
	o.Self = &v
}

// GetNext returns the Next field value if set, zero value otherwise.
func (o *Links) GetNext() string {
	if o == nil || IsNil(o.Next) {
		var ret string
		return ret
	}
	return *o.Next
}

// GetNextOk returns a tuple with the Next field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Links) GetNextOk() (*string, bool) {
	if o == nil || IsNil(o.Next) {
		return nil, false
	}
	return o.Next, true
}

// HasNext returns a boolean if a field has been set.
func (o *Links) HasNext() bool {
	if o != nil && !IsNil(o.Next) {
		return true
	}

	return false
}

// SetNext gets a reference to the given string and assigns it to the Next field.
func (o *Links) SetNext(v string) {
	o.Next = &v
}

func (o Links) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Links) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Prev) {
		toSerialize["prev"] = o.Prev
	}
	if !IsNil(o.Self) {
		toSerialize["self"] = o.Self
	}
	if !IsNil(o.Next) {
		toSerialize["next"] = o.Next
	}
	return toSerialize, nil
}

type NullableLinks struct {
	value *Links
	isSet bool
}

func (v NullableLinks) Get() *Links {
	return v.value
}

func (v *NullableLinks) Set(val *Links) {
	v.value = val
	v.isSet = true
}

func (v NullableLinks) IsSet() bool {
	return v.isSet
}

func (v *NullableLinks) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinks(val *Links) *NullableLinks {
	return &NullableLinks{value: val, isSet: true}
}

func (v NullableLinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinks) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
