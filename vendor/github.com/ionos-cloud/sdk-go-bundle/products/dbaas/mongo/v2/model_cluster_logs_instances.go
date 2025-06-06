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

// checks if the ClusterLogsInstances type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterLogsInstances{}

// ClusterLogsInstances struct for ClusterLogsInstances
type ClusterLogsInstances struct {
	// The name of the MongoDB instance.
	Name     *string                        `json:"name,omitempty"`
	Messages []ClusterLogsInstancesMessages `json:"messages,omitempty"`
}

// NewClusterLogsInstances instantiates a new ClusterLogsInstances object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterLogsInstances() *ClusterLogsInstances {
	this := ClusterLogsInstances{}

	return &this
}

// NewClusterLogsInstancesWithDefaults instantiates a new ClusterLogsInstances object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterLogsInstancesWithDefaults() *ClusterLogsInstances {
	this := ClusterLogsInstances{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ClusterLogsInstances) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLogsInstances) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ClusterLogsInstances) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ClusterLogsInstances) SetName(v string) {
	o.Name = &v
}

// GetMessages returns the Messages field value if set, zero value otherwise.
func (o *ClusterLogsInstances) GetMessages() []ClusterLogsInstancesMessages {
	if o == nil || IsNil(o.Messages) {
		var ret []ClusterLogsInstancesMessages
		return ret
	}
	return o.Messages
}

// GetMessagesOk returns a tuple with the Messages field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterLogsInstances) GetMessagesOk() ([]ClusterLogsInstancesMessages, bool) {
	if o == nil || IsNil(o.Messages) {
		return nil, false
	}
	return o.Messages, true
}

// HasMessages returns a boolean if a field has been set.
func (o *ClusterLogsInstances) HasMessages() bool {
	if o != nil && !IsNil(o.Messages) {
		return true
	}

	return false
}

// SetMessages gets a reference to the given []ClusterLogsInstancesMessages and assigns it to the Messages field.
func (o *ClusterLogsInstances) SetMessages(v []ClusterLogsInstancesMessages) {
	o.Messages = v
}

func (o ClusterLogsInstances) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClusterLogsInstances) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Messages) {
		toSerialize["messages"] = o.Messages
	}
	return toSerialize, nil
}

type NullableClusterLogsInstances struct {
	value *ClusterLogsInstances
	isSet bool
}

func (v NullableClusterLogsInstances) Get() *ClusterLogsInstances {
	return v.value
}

func (v *NullableClusterLogsInstances) Set(val *ClusterLogsInstances) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLogsInstances) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLogsInstances) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLogsInstances(val *ClusterLogsInstances) *NullableClusterLogsInstances {
	return &NullableClusterLogsInstances{value: val, isSet: true}
}

func (v NullableClusterLogsInstances) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLogsInstances) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
