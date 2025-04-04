/*
 * IONOS DBaaS MariaDB REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional MariaDB database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mariadb

import (
	"encoding/json"
)

// checks if the ClusterList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterList{}

// ClusterList List of clusters.
type ClusterList struct {
	// The unique ID of the resource.
	Id    *string           `json:"id,omitempty"`
	Items []ClusterResponse `json:"items,omitempty"`
	// The offset specified in the request (if none was specified, the default offset is 0).
	Offset *int32 `json:"offset,omitempty"`
	// The limit specified in the request (if none was specified, the default limit is 100).
	Limit *int32 `json:"limit,omitempty"`
	// The total number of elements matching the request (without pagination).
	Total *int32           `json:"total,omitempty"`
	Links *PaginationLinks `json:"_links,omitempty"`
}

// NewClusterList instantiates a new ClusterList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterList() *ClusterList {
	this := ClusterList{}

	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit
	var total int32 = 0
	this.Total = &total

	return &this
}

// NewClusterListWithDefaults instantiates a new ClusterList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterListWithDefaults() *ClusterList {
	this := ClusterList{}
	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit
	var total int32 = 0
	this.Total = &total
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ClusterList) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ClusterList) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ClusterList) SetId(v string) {
	o.Id = &v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *ClusterList) GetItems() []ClusterResponse {
	if o == nil || IsNil(o.Items) {
		var ret []ClusterResponse
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetItemsOk() ([]ClusterResponse, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *ClusterList) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []ClusterResponse and assigns it to the Items field.
func (o *ClusterList) SetItems(v []ClusterResponse) {
	o.Items = v
}

// GetOffset returns the Offset field value if set, zero value otherwise.
func (o *ClusterList) GetOffset() int32 {
	if o == nil || IsNil(o.Offset) {
		var ret int32
		return ret
	}
	return *o.Offset
}

// GetOffsetOk returns a tuple with the Offset field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetOffsetOk() (*int32, bool) {
	if o == nil || IsNil(o.Offset) {
		return nil, false
	}
	return o.Offset, true
}

// HasOffset returns a boolean if a field has been set.
func (o *ClusterList) HasOffset() bool {
	if o != nil && !IsNil(o.Offset) {
		return true
	}

	return false
}

// SetOffset gets a reference to the given int32 and assigns it to the Offset field.
func (o *ClusterList) SetOffset(v int32) {
	o.Offset = &v
}

// GetLimit returns the Limit field value if set, zero value otherwise.
func (o *ClusterList) GetLimit() int32 {
	if o == nil || IsNil(o.Limit) {
		var ret int32
		return ret
	}
	return *o.Limit
}

// GetLimitOk returns a tuple with the Limit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetLimitOk() (*int32, bool) {
	if o == nil || IsNil(o.Limit) {
		return nil, false
	}
	return o.Limit, true
}

// HasLimit returns a boolean if a field has been set.
func (o *ClusterList) HasLimit() bool {
	if o != nil && !IsNil(o.Limit) {
		return true
	}

	return false
}

// SetLimit gets a reference to the given int32 and assigns it to the Limit field.
func (o *ClusterList) SetLimit(v int32) {
	o.Limit = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *ClusterList) GetTotal() int32 {
	if o == nil || IsNil(o.Total) {
		var ret int32
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetTotalOk() (*int32, bool) {
	if o == nil || IsNil(o.Total) {
		return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *ClusterList) HasTotal() bool {
	if o != nil && !IsNil(o.Total) {
		return true
	}

	return false
}

// SetTotal gets a reference to the given int32 and assigns it to the Total field.
func (o *ClusterList) SetTotal(v int32) {
	o.Total = &v
}

// GetLinks returns the Links field value if set, zero value otherwise.
func (o *ClusterList) GetLinks() PaginationLinks {
	if o == nil || IsNil(o.Links) {
		var ret PaginationLinks
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterList) GetLinksOk() (*PaginationLinks, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}
	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ClusterList) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given PaginationLinks and assigns it to the Links field.
func (o *ClusterList) SetLinks(v PaginationLinks) {
	o.Links = &v
}

func (o ClusterList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	if !IsNil(o.Offset) {
		toSerialize["offset"] = o.Offset
	}
	if !IsNil(o.Limit) {
		toSerialize["limit"] = o.Limit
	}
	if !IsNil(o.Total) {
		toSerialize["total"] = o.Total
	}
	if !IsNil(o.Links) {
		toSerialize["_links"] = o.Links
	}
	return toSerialize, nil
}

type NullableClusterList struct {
	value *ClusterList
	isSet bool
}

func (v NullableClusterList) Get() *ClusterList {
	return v.value
}

func (v *NullableClusterList) Set(val *ClusterList) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterList) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterList(val *ClusterList) *NullableClusterList {
	return &NullableClusterList{value: val, isSet: true}
}

func (v NullableClusterList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
