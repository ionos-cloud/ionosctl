/*
 * In-Memory DB API
 *
 * API description for the IONOS In-Memory DB
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package inmemorydb

import (
	"encoding/json"
)

// checks if the ReplicaSetReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReplicaSetReadListAllOf{}

// ReplicaSetReadListAllOf struct for ReplicaSetReadListAllOf
type ReplicaSetReadListAllOf struct {
	// ID of the list of ReplicaSet resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of ReplicaSet resources.
	Href string `json:"href"`
	// The list of ReplicaSet resources.
	Items []ReplicaSetRead `json:"items,omitempty"`
}

// NewReplicaSetReadListAllOf instantiates a new ReplicaSetReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicaSetReadListAllOf(id string, type_ string, href string) *ReplicaSetReadListAllOf {
	this := ReplicaSetReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewReplicaSetReadListAllOfWithDefaults instantiates a new ReplicaSetReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicaSetReadListAllOfWithDefaults() *ReplicaSetReadListAllOf {
	this := ReplicaSetReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *ReplicaSetReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *ReplicaSetReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *ReplicaSetReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *ReplicaSetReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *ReplicaSetReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *ReplicaSetReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *ReplicaSetReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *ReplicaSetReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *ReplicaSetReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *ReplicaSetReadListAllOf) GetItems() []ReplicaSetRead {
	if o == nil || IsNil(o.Items) {
		var ret []ReplicaSetRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicaSetReadListAllOf) GetItemsOk() ([]ReplicaSetRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *ReplicaSetReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []ReplicaSetRead and assigns it to the Items field.
func (o *ReplicaSetReadListAllOf) SetItems(v []ReplicaSetRead) {
	o.Items = v
}

func (o ReplicaSetReadListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ReplicaSetReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableReplicaSetReadListAllOf struct {
	value *ReplicaSetReadListAllOf
	isSet bool
}

func (v NullableReplicaSetReadListAllOf) Get() *ReplicaSetReadListAllOf {
	return v.value
}

func (v *NullableReplicaSetReadListAllOf) Set(val *ReplicaSetReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicaSetReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicaSetReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicaSetReadListAllOf(val *ReplicaSetReadListAllOf) *NullableReplicaSetReadListAllOf {
	return &NullableReplicaSetReadListAllOf{value: val, isSet: true}
}

func (v NullableReplicaSetReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicaSetReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
