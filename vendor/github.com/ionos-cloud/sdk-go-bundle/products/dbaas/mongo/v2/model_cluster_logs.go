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

// checks if the ClusterLogs type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterLogs{}

// ClusterLogs The logs of the MongoDB cluster.
type ClusterLogs struct {
	Instances []ClusterLogsInstances `json:"instances,omitempty"`
}

// NewClusterLogs instantiates a new ClusterLogs object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterLogs() *ClusterLogs {
	this := ClusterLogs{}

	return &this
}

// NewClusterLogsWithDefaults instantiates a new ClusterLogs object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterLogsWithDefaults() *ClusterLogs {
	this := ClusterLogs{}
	return &this
}

// GetInstances returns the Instances field value if set, zero value otherwise.
func (o *ClusterLogs) GetInstances() []ClusterLogsInstances {
	if o == nil || IsNil(o.Instances) {
		var ret []ClusterLogsInstances
		return ret
	}
	return o.Instances
}

// GetInstancesOk returns a tuple with the Instances field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLogs) GetInstancesOk() ([]ClusterLogsInstances, bool) {
	if o == nil || IsNil(o.Instances) {
		return nil, false
	}
	return o.Instances, true
}

// HasInstances returns a boolean if a field has been set.
func (o *ClusterLogs) HasInstances() bool {
	if o != nil && !IsNil(o.Instances) {
		return true
	}

	return false
}

// SetInstances gets a reference to the given []ClusterLogsInstances and assigns it to the Instances field.
func (o *ClusterLogs) SetInstances(v []ClusterLogsInstances) {
	o.Instances = v
}

func (o ClusterLogs) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClusterLogs) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Instances) {
		toSerialize["instances"] = o.Instances
	}
	return toSerialize, nil
}

type NullableClusterLogs struct {
	value *ClusterLogs
	isSet bool
}

func (v NullableClusterLogs) Get() *ClusterLogs {
	return v.value
}

func (v *NullableClusterLogs) Set(val *ClusterLogs) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLogs) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLogs) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLogs(val *ClusterLogs) *NullableClusterLogs {
	return &NullableClusterLogs{value: val, isSet: true}
}

func (v NullableClusterLogs) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLogs) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
