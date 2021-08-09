/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.2
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// TargetGroupProperties struct for TargetGroupProperties
type TargetGroupProperties struct {
	// A name of that Target Group
	Name *string `json:"name"`
	// Algorithm for the balancing.
	Algorithm *string `json:"algorithm"`
	// Protocol of the balancing.
	Protocol *string `json:"protocol"`
	// Array of items in that collection
	Targets *[]TargetGroupTarget `json:"targets,omitempty"`
	HealthCheck *TargetGroupHealthCheck `json:"healthCheck,omitempty"`
	HttpHealthCheck *TargetGroupHttpHealthCheck `json:"httpHealthCheck,omitempty"`
}



// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TargetGroupProperties) GetName() *string {
	if o == nil {
		return nil
	}


	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Name, true
}

// SetName sets field value
func (o *TargetGroupProperties) SetName(v string) {


	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}



// GetAlgorithm returns the Algorithm field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TargetGroupProperties) GetAlgorithm() *string {
	if o == nil {
		return nil
	}


	return o.Algorithm

}

// GetAlgorithmOk returns a tuple with the Algorithm field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetAlgorithmOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Algorithm, true
}

// SetAlgorithm sets field value
func (o *TargetGroupProperties) SetAlgorithm(v string) {


	o.Algorithm = &v

}

// HasAlgorithm returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasAlgorithm() bool {
	if o != nil && o.Algorithm != nil {
		return true
	}

	return false
}



// GetProtocol returns the Protocol field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TargetGroupProperties) GetProtocol() *string {
	if o == nil {
		return nil
	}


	return o.Protocol

}

// GetProtocolOk returns a tuple with the Protocol field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetProtocolOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Protocol, true
}

// SetProtocol sets field value
func (o *TargetGroupProperties) SetProtocol(v string) {


	o.Protocol = &v

}

// HasProtocol returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasProtocol() bool {
	if o != nil && o.Protocol != nil {
		return true
	}

	return false
}



// GetTargets returns the Targets field value
// If the value is explicit nil, the zero value for []TargetGroupTarget will be returned
func (o *TargetGroupProperties) GetTargets() *[]TargetGroupTarget {
	if o == nil {
		return nil
	}


	return o.Targets

}

// GetTargetsOk returns a tuple with the Targets field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetTargetsOk() (*[]TargetGroupTarget, bool) {
	if o == nil {
		return nil, false
	}


	return o.Targets, true
}

// SetTargets sets field value
func (o *TargetGroupProperties) SetTargets(v []TargetGroupTarget) {


	o.Targets = &v

}

// HasTargets returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasTargets() bool {
	if o != nil && o.Targets != nil {
		return true
	}

	return false
}



// GetHealthCheck returns the HealthCheck field value
// If the value is explicit nil, the zero value for TargetGroupHealthCheck will be returned
func (o *TargetGroupProperties) GetHealthCheck() *TargetGroupHealthCheck {
	if o == nil {
		return nil
	}


	return o.HealthCheck

}

// GetHealthCheckOk returns a tuple with the HealthCheck field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetHealthCheckOk() (*TargetGroupHealthCheck, bool) {
	if o == nil {
		return nil, false
	}


	return o.HealthCheck, true
}

// SetHealthCheck sets field value
func (o *TargetGroupProperties) SetHealthCheck(v TargetGroupHealthCheck) {


	o.HealthCheck = &v

}

// HasHealthCheck returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasHealthCheck() bool {
	if o != nil && o.HealthCheck != nil {
		return true
	}

	return false
}



// GetHttpHealthCheck returns the HttpHealthCheck field value
// If the value is explicit nil, the zero value for TargetGroupHttpHealthCheck will be returned
func (o *TargetGroupProperties) GetHttpHealthCheck() *TargetGroupHttpHealthCheck {
	if o == nil {
		return nil
	}


	return o.HttpHealthCheck

}

// GetHttpHealthCheckOk returns a tuple with the HttpHealthCheck field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupProperties) GetHttpHealthCheckOk() (*TargetGroupHttpHealthCheck, bool) {
	if o == nil {
		return nil, false
	}


	return o.HttpHealthCheck, true
}

// SetHttpHealthCheck sets field value
func (o *TargetGroupProperties) SetHttpHealthCheck(v TargetGroupHttpHealthCheck) {


	o.HttpHealthCheck = &v

}

// HasHttpHealthCheck returns a boolean if a field has been set.
func (o *TargetGroupProperties) HasHttpHealthCheck() bool {
	if o != nil && o.HttpHealthCheck != nil {
		return true
	}

	return false
}


func (o TargetGroupProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	

	if o.Algorithm != nil {
		toSerialize["algorithm"] = o.Algorithm
	}
	

	if o.Protocol != nil {
		toSerialize["protocol"] = o.Protocol
	}
	

	if o.Targets != nil {
		toSerialize["targets"] = o.Targets
	}
	

	if o.HealthCheck != nil {
		toSerialize["healthCheck"] = o.HealthCheck
	}
	

	if o.HttpHealthCheck != nil {
		toSerialize["httpHealthCheck"] = o.HttpHealthCheck
	}
	
	return json.Marshal(toSerialize)
}

type NullableTargetGroupProperties struct {
	value *TargetGroupProperties
	isSet bool
}

func (v NullableTargetGroupProperties) Get() *TargetGroupProperties {
	return v.value
}

func (v *NullableTargetGroupProperties) Set(val *TargetGroupProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableTargetGroupProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableTargetGroupProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTargetGroupProperties(val *TargetGroupProperties) *NullableTargetGroupProperties {
	return &NullableTargetGroupProperties{value: val, isSet: true}
}

func (v NullableTargetGroupProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTargetGroupProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


