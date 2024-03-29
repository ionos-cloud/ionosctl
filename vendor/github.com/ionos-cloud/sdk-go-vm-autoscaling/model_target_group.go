/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// TargetGroup In order to link VM to ALB, target group must be provided
type TargetGroup struct {
	// id
	TargetGroupId *string `json:"targetGroupId"`
	// port
	Port *int32 `json:"port"`
	// weight
	Weight *int32 `json:"weight"`
}

// NewTargetGroup instantiates a new TargetGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTargetGroup(targetGroupId string, port int32, weight int32) *TargetGroup {
	this := TargetGroup{}

	this.TargetGroupId = &targetGroupId
	this.Port = &port
	this.Weight = &weight

	return &this
}

// NewTargetGroupWithDefaults instantiates a new TargetGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTargetGroupWithDefaults() *TargetGroup {
	this := TargetGroup{}
	return &this
}

// GetTargetGroupId returns the TargetGroupId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TargetGroup) GetTargetGroupId() *string {
	if o == nil {
		return nil
	}

	return o.TargetGroupId

}

// GetTargetGroupIdOk returns a tuple with the TargetGroupId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroup) GetTargetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.TargetGroupId, true
}

// SetTargetGroupId sets field value
func (o *TargetGroup) SetTargetGroupId(v string) {

	o.TargetGroupId = &v

}

// HasTargetGroupId returns a boolean if a field has been set.
func (o *TargetGroup) HasTargetGroupId() bool {
	if o != nil && o.TargetGroupId != nil {
		return true
	}

	return false
}

// GetPort returns the Port field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroup) GetPort() *int32 {
	if o == nil {
		return nil
	}

	return o.Port

}

// GetPortOk returns a tuple with the Port field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroup) GetPortOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Port, true
}

// SetPort sets field value
func (o *TargetGroup) SetPort(v int32) {

	o.Port = &v

}

// HasPort returns a boolean if a field has been set.
func (o *TargetGroup) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

// GetWeight returns the Weight field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroup) GetWeight() *int32 {
	if o == nil {
		return nil
	}

	return o.Weight

}

// GetWeightOk returns a tuple with the Weight field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroup) GetWeightOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Weight, true
}

// SetWeight sets field value
func (o *TargetGroup) SetWeight(v int32) {

	o.Weight = &v

}

// HasWeight returns a boolean if a field has been set.
func (o *TargetGroup) HasWeight() bool {
	if o != nil && o.Weight != nil {
		return true
	}

	return false
}

func (o TargetGroup) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TargetGroupId != nil {
		toSerialize["targetGroupId"] = o.TargetGroupId
	}

	if o.Port != nil {
		toSerialize["port"] = o.Port
	}

	if o.Weight != nil {
		toSerialize["weight"] = o.Weight
	}

	return json.Marshal(toSerialize)
}

type NullableTargetGroup struct {
	value *TargetGroup
	isSet bool
}

func (v NullableTargetGroup) Get() *TargetGroup {
	return v.value
}

func (v *NullableTargetGroup) Set(val *TargetGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableTargetGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableTargetGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTargetGroup(val *TargetGroup) *NullableTargetGroup {
	return &NullableTargetGroup{value: val, isSet: true}
}

func (v NullableTargetGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTargetGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
