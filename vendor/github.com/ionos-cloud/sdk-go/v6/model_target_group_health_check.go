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

// TargetGroupHealthCheck struct for TargetGroupHealthCheck
type TargetGroupHealthCheck struct {
	// It specifies the time (in milliseconds) for a target VM in this pool to answer the check. If a target VM has CheckInterval set and CheckTimeout is set too, then the smaller value of the two is used after the TCP connection is established.
	CheckTimeout *int32 `json:"checkTimeout,omitempty"`
	// It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed. If unset, the default of 5 seconds will be used.
	ConnectTimeout *int32 `json:"connectTimeout,omitempty"`
	// TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side. If unset, the default of 50 seconds will be used.
	TargetTimeout *int32 `json:"targetTimeout,omitempty"`
	// Retries specifies the number of retries to perform on a target VM after a connection failure. If unset, the default value of 3 will be used. (valid range: [0, 65535])
	Retries *int32 `json:"retries,omitempty"`
}



// GetCheckTimeout returns the CheckTimeout field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroupHealthCheck) GetCheckTimeout() *int32 {
	if o == nil {
		return nil
	}


	return o.CheckTimeout

}

// GetCheckTimeoutOk returns a tuple with the CheckTimeout field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupHealthCheck) GetCheckTimeoutOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.CheckTimeout, true
}

// SetCheckTimeout sets field value
func (o *TargetGroupHealthCheck) SetCheckTimeout(v int32) {


	o.CheckTimeout = &v

}

// HasCheckTimeout returns a boolean if a field has been set.
func (o *TargetGroupHealthCheck) HasCheckTimeout() bool {
	if o != nil && o.CheckTimeout != nil {
		return true
	}

	return false
}



// GetConnectTimeout returns the ConnectTimeout field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroupHealthCheck) GetConnectTimeout() *int32 {
	if o == nil {
		return nil
	}


	return o.ConnectTimeout

}

// GetConnectTimeoutOk returns a tuple with the ConnectTimeout field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupHealthCheck) GetConnectTimeoutOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.ConnectTimeout, true
}

// SetConnectTimeout sets field value
func (o *TargetGroupHealthCheck) SetConnectTimeout(v int32) {


	o.ConnectTimeout = &v

}

// HasConnectTimeout returns a boolean if a field has been set.
func (o *TargetGroupHealthCheck) HasConnectTimeout() bool {
	if o != nil && o.ConnectTimeout != nil {
		return true
	}

	return false
}



// GetTargetTimeout returns the TargetTimeout field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroupHealthCheck) GetTargetTimeout() *int32 {
	if o == nil {
		return nil
	}


	return o.TargetTimeout

}

// GetTargetTimeoutOk returns a tuple with the TargetTimeout field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupHealthCheck) GetTargetTimeoutOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.TargetTimeout, true
}

// SetTargetTimeout sets field value
func (o *TargetGroupHealthCheck) SetTargetTimeout(v int32) {


	o.TargetTimeout = &v

}

// HasTargetTimeout returns a boolean if a field has been set.
func (o *TargetGroupHealthCheck) HasTargetTimeout() bool {
	if o != nil && o.TargetTimeout != nil {
		return true
	}

	return false
}



// GetRetries returns the Retries field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *TargetGroupHealthCheck) GetRetries() *int32 {
	if o == nil {
		return nil
	}


	return o.Retries

}

// GetRetriesOk returns a tuple with the Retries field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TargetGroupHealthCheck) GetRetriesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Retries, true
}

// SetRetries sets field value
func (o *TargetGroupHealthCheck) SetRetries(v int32) {


	o.Retries = &v

}

// HasRetries returns a boolean if a field has been set.
func (o *TargetGroupHealthCheck) HasRetries() bool {
	if o != nil && o.Retries != nil {
		return true
	}

	return false
}


func (o TargetGroupHealthCheck) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.CheckTimeout != nil {
		toSerialize["checkTimeout"] = o.CheckTimeout
	}
	

	if o.ConnectTimeout != nil {
		toSerialize["connectTimeout"] = o.ConnectTimeout
	}
	

	if o.TargetTimeout != nil {
		toSerialize["targetTimeout"] = o.TargetTimeout
	}
	

	if o.Retries != nil {
		toSerialize["retries"] = o.Retries
	}
	
	return json.Marshal(toSerialize)
}

type NullableTargetGroupHealthCheck struct {
	value *TargetGroupHealthCheck
	isSet bool
}

func (v NullableTargetGroupHealthCheck) Get() *TargetGroupHealthCheck {
	return v.value
}

func (v *NullableTargetGroupHealthCheck) Set(val *TargetGroupHealthCheck) {
	v.value = val
	v.isSet = true
}

func (v NullableTargetGroupHealthCheck) IsSet() bool {
	return v.isSet
}

func (v *NullableTargetGroupHealthCheck) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTargetGroupHealthCheck(val *TargetGroupHealthCheck) *NullableTargetGroupHealthCheck {
	return &NullableTargetGroupHealthCheck{value: val, isSet: true}
}

func (v NullableTargetGroupHealthCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTargetGroupHealthCheck) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


