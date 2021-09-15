/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ApplicationLoadBalancerHttpRuleCondition struct for ApplicationLoadBalancerHttpRuleCondition
type ApplicationLoadBalancerHttpRuleCondition struct {
	// Type of the Http Rule condition.
	Type *string `json:"type"`
	// Matching rule for the Http Rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST and COOKIE types; must be null when type is SOURCE_IP.
	Condition *string `json:"condition"`
	// Specifies whether the condition is negated or not; default: false
	Negate *bool `json:"negate,omitempty"`
	// Must be null when type is PATH, METHOD, HOST or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, QUERY.
	Key *string `json:"key,omitempty"`
	// Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.
	Value *string `json:"value,omitempty"`
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetType() *string {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *ApplicationLoadBalancerHttpRuleCondition) SetType(v string) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerHttpRuleCondition) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetCondition returns the Condition field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetCondition() *string {
	if o == nil {
		return nil
	}


	return o.Condition

}

// GetConditionOk returns a tuple with the Condition field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetConditionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Condition, true
}

// SetCondition sets field value
func (o *ApplicationLoadBalancerHttpRuleCondition) SetCondition(v string) {


	o.Condition = &v

}

// HasCondition returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerHttpRuleCondition) HasCondition() bool {
	if o != nil && o.Condition != nil {
		return true
	}

	return false
}



// GetNegate returns the Negate field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetNegate() *bool {
	if o == nil {
		return nil
	}


	return o.Negate

}

// GetNegateOk returns a tuple with the Negate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetNegateOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}


	return o.Negate, true
}

// SetNegate sets field value
func (o *ApplicationLoadBalancerHttpRuleCondition) SetNegate(v bool) {


	o.Negate = &v

}

// HasNegate returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerHttpRuleCondition) HasNegate() bool {
	if o != nil && o.Negate != nil {
		return true
	}

	return false
}



// GetKey returns the Key field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetKey() *string {
	if o == nil {
		return nil
	}


	return o.Key

}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Key, true
}

// SetKey sets field value
func (o *ApplicationLoadBalancerHttpRuleCondition) SetKey(v string) {


	o.Key = &v

}

// HasKey returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerHttpRuleCondition) HasKey() bool {
	if o != nil && o.Key != nil {
		return true
	}

	return false
}



// GetValue returns the Value field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetValue() *string {
	if o == nil {
		return nil
	}


	return o.Value

}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerHttpRuleCondition) GetValueOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Value, true
}

// SetValue sets field value
func (o *ApplicationLoadBalancerHttpRuleCondition) SetValue(v string) {


	o.Value = &v

}

// HasValue returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerHttpRuleCondition) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}


func (o ApplicationLoadBalancerHttpRuleCondition) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	

	if o.Condition != nil {
		toSerialize["condition"] = o.Condition
	}
	

	if o.Negate != nil {
		toSerialize["negate"] = o.Negate
	}
	

	if o.Key != nil {
		toSerialize["key"] = o.Key
	}
	

	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	
	return json.Marshal(toSerialize)
}

type NullableApplicationLoadBalancerHttpRuleCondition struct {
	value *ApplicationLoadBalancerHttpRuleCondition
	isSet bool
}

func (v NullableApplicationLoadBalancerHttpRuleCondition) Get() *ApplicationLoadBalancerHttpRuleCondition {
	return v.value
}

func (v *NullableApplicationLoadBalancerHttpRuleCondition) Set(val *ApplicationLoadBalancerHttpRuleCondition) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationLoadBalancerHttpRuleCondition) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationLoadBalancerHttpRuleCondition) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationLoadBalancerHttpRuleCondition(val *ApplicationLoadBalancerHttpRuleCondition) *NullableApplicationLoadBalancerHttpRuleCondition {
	return &NullableApplicationLoadBalancerHttpRuleCondition{value: val, isSet: true}
}

func (v NullableApplicationLoadBalancerHttpRuleCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationLoadBalancerHttpRuleCondition) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


