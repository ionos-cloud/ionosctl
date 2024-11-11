/*
 * Kafka as a Service API
 *
 * An managed Apache Kafka cluster is designed to be highly fault-tolerant and scalable, allowing large volumes of data to be ingested, stored, and processed in real-time. By distributing data across multiple brokers, Kafka achieves high throughput and low latency, making it suitable for applications requiring real-time data processing and analytics.
 *
 * API version: 1.7.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ClusterCreate struct for ClusterCreate
type ClusterCreate struct {
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *Cluster                `json:"properties"`
}

// NewClusterCreate instantiates a new ClusterCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterCreate(properties Cluster) *ClusterCreate {
	this := ClusterCreate{}

	this.Properties = &properties

	return &this
}

// NewClusterCreateWithDefaults instantiates a new ClusterCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterCreateWithDefaults() *ClusterCreate {
	this := ClusterCreate{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *ClusterCreate) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterCreate) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *ClusterCreate) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *ClusterCreate) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Cluster will be returned
func (o *ClusterCreate) GetProperties() *Cluster {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterCreate) GetPropertiesOk() (*Cluster, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *ClusterCreate) SetProperties(v Cluster) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *ClusterCreate) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o ClusterCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableClusterCreate struct {
	value *ClusterCreate
	isSet bool
}

func (v NullableClusterCreate) Get() *ClusterCreate {
	return v.value
}

func (v *NullableClusterCreate) Set(val *ClusterCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterCreate(val *ClusterCreate) *NullableClusterCreate {
	return &NullableClusterCreate{value: val, isSet: true}
}

func (v NullableClusterCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
