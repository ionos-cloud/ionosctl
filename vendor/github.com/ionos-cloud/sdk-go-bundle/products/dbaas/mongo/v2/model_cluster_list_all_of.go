/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.  MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mongo

import (
	"encoding/json"
)

// checks if the ClusterListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterListAllOf{}

// ClusterListAllOf struct for ClusterListAllOf
type ClusterListAllOf struct {
	Type *ResourceType `json:"type,omitempty"`
	// The unique ID of the resource.
	Id    *string           `json:"id,omitempty"`
	Items []ClusterResponse `json:"items,omitempty"`
}

// NewClusterListAllOf instantiates a new ClusterListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterListAllOf() *ClusterListAllOf {
	this := ClusterListAllOf{}

	return &this
}

// NewClusterListAllOfWithDefaults instantiates a new ClusterListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterListAllOfWithDefaults() *ClusterListAllOf {
	this := ClusterListAllOf{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ClusterListAllOf) GetType() ResourceType {
	if o == nil || IsNil(o.Type) {
		var ret ResourceType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterListAllOf) GetTypeOk() (*ResourceType, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ClusterListAllOf) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given ResourceType and assigns it to the Type field.
func (o *ClusterListAllOf) SetType(v ResourceType) {
	o.Type = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ClusterListAllOf) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterListAllOf) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ClusterListAllOf) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ClusterListAllOf) SetId(v string) {
	o.Id = &v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *ClusterListAllOf) GetItems() []ClusterResponse {
	if o == nil || IsNil(o.Items) {
		var ret []ClusterResponse
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterListAllOf) GetItemsOk() ([]ClusterResponse, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *ClusterListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []ClusterResponse and assigns it to the Items field.
func (o *ClusterListAllOf) SetItems(v []ClusterResponse) {
	o.Items = v
}

func (o ClusterListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClusterListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableClusterListAllOf struct {
	value *ClusterListAllOf
	isSet bool
}

func (v NullableClusterListAllOf) Get() *ClusterListAllOf {
	return v.value
}

func (v *NullableClusterListAllOf) Set(val *ClusterListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterListAllOf(val *ClusterListAllOf) *NullableClusterListAllOf {
	return &NullableClusterListAllOf{value: val, isSet: true}
}

func (v NullableClusterListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
