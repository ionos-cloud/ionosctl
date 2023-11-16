/*
 * IONOS Logging REST API
 *
 * Logging as a Service (LaaS) is a service that provides a centralized logging system where users are able to push and aggregate their system or application logs. This service also provides a visualization platform where users are able to observe, search and filter the logs and also create dashboards and alerts for their data points. This service can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an API. The API allows you to create logging pipelines or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// PipelineResponseAllOf1 struct for PipelineResponseAllOf1
type PipelineResponseAllOf1 struct {
	Destinations *[]Destination `json:"destinations,omitempty"`
}

// NewPipelineResponseAllOf1 instantiates a new PipelineResponseAllOf1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineResponseAllOf1() *PipelineResponseAllOf1 {
	this := PipelineResponseAllOf1{}

	return &this
}

// NewPipelineResponseAllOf1WithDefaults instantiates a new PipelineResponseAllOf1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineResponseAllOf1WithDefaults() *PipelineResponseAllOf1 {
	this := PipelineResponseAllOf1{}
	return &this
}

// GetDestinations returns the Destinations field value
// If the value is explicit nil, the zero value for []Destination will be returned
func (o *PipelineResponseAllOf1) GetDestinations() *[]Destination {
	if o == nil {
		return nil
	}

	return o.Destinations

}

// GetDestinationsOk returns a tuple with the Destinations field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineResponseAllOf1) GetDestinationsOk() (*[]Destination, bool) {
	if o == nil {
		return nil, false
	}

	return o.Destinations, true
}

// SetDestinations sets field value
func (o *PipelineResponseAllOf1) SetDestinations(v []Destination) {

	o.Destinations = &v

}

// HasDestinations returns a boolean if a field has been set.
func (o *PipelineResponseAllOf1) HasDestinations() bool {
	if o != nil && o.Destinations != nil {
		return true
	}

	return false
}

func (o PipelineResponseAllOf1) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Destinations != nil {
		toSerialize["destinations"] = o.Destinations
	}

	return json.Marshal(toSerialize)
}

type NullablePipelineResponseAllOf1 struct {
	value *PipelineResponseAllOf1
	isSet bool
}

func (v NullablePipelineResponseAllOf1) Get() *PipelineResponseAllOf1 {
	return v.value
}

func (v *NullablePipelineResponseAllOf1) Set(val *PipelineResponseAllOf1) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineResponseAllOf1) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineResponseAllOf1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineResponseAllOf1(val *PipelineResponseAllOf1) *NullablePipelineResponseAllOf1 {
	return &NullablePipelineResponseAllOf1{value: val, isSet: true}
}

func (v NullablePipelineResponseAllOf1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineResponseAllOf1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
