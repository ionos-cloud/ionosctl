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
	"time"
)

// ClusterMetadata struct for ClusterMetadata
type ClusterMetadata struct {
	// The ISO 8601 creation timestamp.
	CreatedDate *IonosTime `json:"createdDate,omitempty"`
	// Unique name of the identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`
	// Unique id of the identity that created the resource.
	CreatedByUserId *string `json:"createdByUserId,omitempty"`
	// The ISO 8601 modified timestamp.
	LastModifiedDate *IonosTime `json:"lastModifiedDate,omitempty"`
	// Unique name of the identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// Unique id of the identity that last modified the resource.
	LastModifiedByUserId *string `json:"lastModifiedByUserId,omitempty"`
	// Unique name of the resource.
	ResourceURN *string `json:"resourceURN,omitempty"`
	// State of the resource. Resource states: `AVAILABLE`: There are no pending modification requests for this item. `BUSY`: There is at least one modification request pending and all following requests will be queued. `DEPLOYING`: The resource is being created. `FAILED`: The creation of the resource failed. `UPDATING`: The resource is being updated. `FAILED_UPDATING`: An update to the resource was not successful. `DESTROYING`: A delete command was issued, and the resource is being deleted.
	State *string `json:"state"`
	// A human readable message describing the current state. In case of an error, the message will contain a detailed error message.
	Message *string `json:"message,omitempty"`
	// IP addresses and ports of cluster brokers.
	BrokerAddresses *[]string `json:"brokerAddresses,omitempty"`
}

// NewClusterMetadata instantiates a new ClusterMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterMetadata(state string) *ClusterMetadata {
	this := ClusterMetadata{}

	this.State = &state

	return &this
}

// NewClusterMetadataWithDefaults instantiates a new ClusterMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterMetadataWithDefaults() *ClusterMetadata {
	this := ClusterMetadata{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *ClusterMetadata) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time

}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *ClusterMetadata) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}

}

// HasCreatedDate returns a boolean if a field has been set.
func (o *ClusterMetadata) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}

// GetCreatedBy returns the CreatedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetCreatedBy() *string {
	if o == nil {
		return nil
	}

	return o.CreatedBy

}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetCreatedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *ClusterMetadata) SetCreatedBy(v string) {

	o.CreatedBy = &v

}

// HasCreatedBy returns a boolean if a field has been set.
func (o *ClusterMetadata) HasCreatedBy() bool {
	if o != nil && o.CreatedBy != nil {
		return true
	}

	return false
}

// GetCreatedByUserId returns the CreatedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetCreatedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.CreatedByUserId

}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedByUserId, true
}

// SetCreatedByUserId sets field value
func (o *ClusterMetadata) SetCreatedByUserId(v string) {

	o.CreatedByUserId = &v

}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *ClusterMetadata) HasCreatedByUserId() bool {
	if o != nil && o.CreatedByUserId != nil {
		return true
	}

	return false
}

// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *ClusterMetadata) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time

}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *ClusterMetadata) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}

}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *ClusterMetadata) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}

// GetLastModifiedBy returns the LastModifiedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetLastModifiedBy() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedBy

}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetLastModifiedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedBy, true
}

// SetLastModifiedBy sets field value
func (o *ClusterMetadata) SetLastModifiedBy(v string) {

	o.LastModifiedBy = &v

}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *ClusterMetadata) HasLastModifiedBy() bool {
	if o != nil && o.LastModifiedBy != nil {
		return true
	}

	return false
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetLastModifiedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedByUserId

}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedByUserId, true
}

// SetLastModifiedByUserId sets field value
func (o *ClusterMetadata) SetLastModifiedByUserId(v string) {

	o.LastModifiedByUserId = &v

}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *ClusterMetadata) HasLastModifiedByUserId() bool {
	if o != nil && o.LastModifiedByUserId != nil {
		return true
	}

	return false
}

// GetResourceURN returns the ResourceURN field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetResourceURN() *string {
	if o == nil {
		return nil
	}

	return o.ResourceURN

}

// GetResourceURNOk returns a tuple with the ResourceURN field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetResourceURNOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ResourceURN, true
}

// SetResourceURN sets field value
func (o *ClusterMetadata) SetResourceURN(v string) {

	o.ResourceURN = &v

}

// HasResourceURN returns a boolean if a field has been set.
func (o *ClusterMetadata) HasResourceURN() bool {
	if o != nil && o.ResourceURN != nil {
		return true
	}

	return false
}

// GetState returns the State field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetState() *string {
	if o == nil {
		return nil
	}

	return o.State

}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.State, true
}

// SetState sets field value
func (o *ClusterMetadata) SetState(v string) {

	o.State = &v

}

// HasState returns a boolean if a field has been set.
func (o *ClusterMetadata) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// GetMessage returns the Message field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterMetadata) GetMessage() *string {
	if o == nil {
		return nil
	}

	return o.Message

}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Message, true
}

// SetMessage sets field value
func (o *ClusterMetadata) SetMessage(v string) {

	o.Message = &v

}

// HasMessage returns a boolean if a field has been set.
func (o *ClusterMetadata) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

// GetBrokerAddresses returns the BrokerAddresses field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *ClusterMetadata) GetBrokerAddresses() *[]string {
	if o == nil {
		return nil
	}

	return o.BrokerAddresses

}

// GetBrokerAddressesOk returns a tuple with the BrokerAddresses field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterMetadata) GetBrokerAddressesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.BrokerAddresses, true
}

// SetBrokerAddresses sets field value
func (o *ClusterMetadata) SetBrokerAddresses(v []string) {

	o.BrokerAddresses = &v

}

// HasBrokerAddresses returns a boolean if a field has been set.
func (o *ClusterMetadata) HasBrokerAddresses() bool {
	if o != nil && o.BrokerAddresses != nil {
		return true
	}

	return false
}

func (o ClusterMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}

	if o.CreatedBy != nil {
		toSerialize["createdBy"] = o.CreatedBy
	}

	if o.CreatedByUserId != nil {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}

	if o.LastModifiedDate != nil {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}

	if o.LastModifiedBy != nil {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}

	if o.LastModifiedByUserId != nil {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}

	if o.ResourceURN != nil {
		toSerialize["resourceURN"] = o.ResourceURN
	}

	if o.State != nil {
		toSerialize["state"] = o.State
	}

	if o.Message != nil {
		toSerialize["message"] = o.Message
	}

	if o.BrokerAddresses != nil {
		toSerialize["brokerAddresses"] = o.BrokerAddresses
	}

	return json.Marshal(toSerialize)
}

type NullableClusterMetadata struct {
	value *ClusterMetadata
	isSet bool
}

func (v NullableClusterMetadata) Get() *ClusterMetadata {
	return v.value
}

func (v *NullableClusterMetadata) Set(val *ClusterMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterMetadata(val *ClusterMetadata) *NullableClusterMetadata {
	return &NullableClusterMetadata{value: val, isSet: true}
}

func (v NullableClusterMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
