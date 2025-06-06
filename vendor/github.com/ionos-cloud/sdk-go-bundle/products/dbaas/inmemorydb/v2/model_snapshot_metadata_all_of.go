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

	"time"
)

// checks if the SnapshotMetadataAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SnapshotMetadataAllOf{}

// SnapshotMetadataAllOf struct for SnapshotMetadataAllOf
type SnapshotMetadataAllOf struct {
	// The ID of the In-Memory DB replica set the snapshot is taken from.
	ReplicasetId string `json:"replicasetId"`
	// The time the snapshot was dumped from the replica set.
	SnapshotTime *IonosTime `json:"snapshotTime,omitempty"`
	// The ID of the datacenter the snapshot was created in. Please mind, that the snapshot is not available in other datacenters.
	DatacenterId string `json:"datacenterId"`
}

// NewSnapshotMetadataAllOf instantiates a new SnapshotMetadataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSnapshotMetadataAllOf(replicasetId string, datacenterId string) *SnapshotMetadataAllOf {
	this := SnapshotMetadataAllOf{}

	this.ReplicasetId = replicasetId
	this.DatacenterId = datacenterId

	return &this
}

// NewSnapshotMetadataAllOfWithDefaults instantiates a new SnapshotMetadataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSnapshotMetadataAllOfWithDefaults() *SnapshotMetadataAllOf {
	this := SnapshotMetadataAllOf{}
	return &this
}

// GetReplicasetId returns the ReplicasetId field value
func (o *SnapshotMetadataAllOf) GetReplicasetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ReplicasetId
}

// GetReplicasetIdOk returns a tuple with the ReplicasetId field value
// and a boolean to check if the value has been set.
func (o *SnapshotMetadataAllOf) GetReplicasetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReplicasetId, true
}

// SetReplicasetId sets field value
func (o *SnapshotMetadataAllOf) SetReplicasetId(v string) {
	o.ReplicasetId = v
}

// GetSnapshotTime returns the SnapshotTime field value if set, zero value otherwise.
func (o *SnapshotMetadataAllOf) GetSnapshotTime() time.Time {
	if o == nil || IsNil(o.SnapshotTime) {
		var ret time.Time
		return ret
	}
	return o.SnapshotTime.Time
}

// GetSnapshotTimeOk returns a tuple with the SnapshotTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SnapshotMetadataAllOf) GetSnapshotTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.SnapshotTime) {
		return nil, false
	}
	return &o.SnapshotTime.Time, true
}

// HasSnapshotTime returns a boolean if a field has been set.
func (o *SnapshotMetadataAllOf) HasSnapshotTime() bool {
	if o != nil && !IsNil(o.SnapshotTime) {
		return true
	}

	return false
}

// SetSnapshotTime gets a reference to the given time.Time and assigns it to the SnapshotTime field.
func (o *SnapshotMetadataAllOf) SetSnapshotTime(v time.Time) {
	o.SnapshotTime = &IonosTime{v}
}

// GetDatacenterId returns the DatacenterId field value
func (o *SnapshotMetadataAllOf) GetDatacenterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DatacenterId
}

// GetDatacenterIdOk returns a tuple with the DatacenterId field value
// and a boolean to check if the value has been set.
func (o *SnapshotMetadataAllOf) GetDatacenterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DatacenterId, true
}

// SetDatacenterId sets field value
func (o *SnapshotMetadataAllOf) SetDatacenterId(v string) {
	o.DatacenterId = v
}

func (o SnapshotMetadataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SnapshotMetadataAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["replicasetId"] = o.ReplicasetId
	if !IsNil(o.SnapshotTime) {
		toSerialize["snapshotTime"] = o.SnapshotTime
	}
	toSerialize["datacenterId"] = o.DatacenterId
	return toSerialize, nil
}

type NullableSnapshotMetadataAllOf struct {
	value *SnapshotMetadataAllOf
	isSet bool
}

func (v NullableSnapshotMetadataAllOf) Get() *SnapshotMetadataAllOf {
	return v.value
}

func (v *NullableSnapshotMetadataAllOf) Set(val *SnapshotMetadataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableSnapshotMetadataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableSnapshotMetadataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSnapshotMetadataAllOf(val *SnapshotMetadataAllOf) *NullableSnapshotMetadataAllOf {
	return &NullableSnapshotMetadataAllOf{value: val, isSet: true}
}

func (v NullableSnapshotMetadataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSnapshotMetadataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
