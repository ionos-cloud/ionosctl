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

// ApplicationLoadBalancerProperties struct for ApplicationLoadBalancerProperties
type ApplicationLoadBalancerProperties struct {
	// A name of that Application Load Balancer
	Name *string `json:"name"`
	// Id of the listening LAN. (inbound)
	ListenerLan *int32 `json:"listenerLan"`
	// Collection of IP addresses of the Application Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
	Ips *[]string `json:"ips,omitempty"`
	// Id of the balanced private target LAN. (outbound)
	TargetLan *int32 `json:"targetLan"`
	// Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.
	LbPrivateIps *[]string `json:"lbPrivateIps,omitempty"`
}



// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ApplicationLoadBalancerProperties) GetName() *string {
	if o == nil {
		return nil
	}


	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Name, true
}

// SetName sets field value
func (o *ApplicationLoadBalancerProperties) SetName(v string) {


	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}



// GetListenerLan returns the ListenerLan field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ApplicationLoadBalancerProperties) GetListenerLan() *int32 {
	if o == nil {
		return nil
	}


	return o.ListenerLan

}

// GetListenerLanOk returns a tuple with the ListenerLan field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerProperties) GetListenerLanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.ListenerLan, true
}

// SetListenerLan sets field value
func (o *ApplicationLoadBalancerProperties) SetListenerLan(v int32) {


	o.ListenerLan = &v

}

// HasListenerLan returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerProperties) HasListenerLan() bool {
	if o != nil && o.ListenerLan != nil {
		return true
	}

	return false
}



// GetIps returns the Ips field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *ApplicationLoadBalancerProperties) GetIps() *[]string {
	if o == nil {
		return nil
	}


	return o.Ips

}

// GetIpsOk returns a tuple with the Ips field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerProperties) GetIpsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Ips, true
}

// SetIps sets field value
func (o *ApplicationLoadBalancerProperties) SetIps(v []string) {


	o.Ips = &v

}

// HasIps returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerProperties) HasIps() bool {
	if o != nil && o.Ips != nil {
		return true
	}

	return false
}



// GetTargetLan returns the TargetLan field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ApplicationLoadBalancerProperties) GetTargetLan() *int32 {
	if o == nil {
		return nil
	}


	return o.TargetLan

}

// GetTargetLanOk returns a tuple with the TargetLan field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerProperties) GetTargetLanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.TargetLan, true
}

// SetTargetLan sets field value
func (o *ApplicationLoadBalancerProperties) SetTargetLan(v int32) {


	o.TargetLan = &v

}

// HasTargetLan returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerProperties) HasTargetLan() bool {
	if o != nil && o.TargetLan != nil {
		return true
	}

	return false
}



// GetLbPrivateIps returns the LbPrivateIps field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *ApplicationLoadBalancerProperties) GetLbPrivateIps() *[]string {
	if o == nil {
		return nil
	}


	return o.LbPrivateIps

}

// GetLbPrivateIpsOk returns a tuple with the LbPrivateIps field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApplicationLoadBalancerProperties) GetLbPrivateIpsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LbPrivateIps, true
}

// SetLbPrivateIps sets field value
func (o *ApplicationLoadBalancerProperties) SetLbPrivateIps(v []string) {


	o.LbPrivateIps = &v

}

// HasLbPrivateIps returns a boolean if a field has been set.
func (o *ApplicationLoadBalancerProperties) HasLbPrivateIps() bool {
	if o != nil && o.LbPrivateIps != nil {
		return true
	}

	return false
}


func (o ApplicationLoadBalancerProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	

	if o.ListenerLan != nil {
		toSerialize["listenerLan"] = o.ListenerLan
	}
	

	if o.Ips != nil {
		toSerialize["ips"] = o.Ips
	}
	

	if o.TargetLan != nil {
		toSerialize["targetLan"] = o.TargetLan
	}
	

	if o.LbPrivateIps != nil {
		toSerialize["lbPrivateIps"] = o.LbPrivateIps
	}
	
	return json.Marshal(toSerialize)
}

type NullableApplicationLoadBalancerProperties struct {
	value *ApplicationLoadBalancerProperties
	isSet bool
}

func (v NullableApplicationLoadBalancerProperties) Get() *ApplicationLoadBalancerProperties {
	return v.value
}

func (v *NullableApplicationLoadBalancerProperties) Set(val *ApplicationLoadBalancerProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableApplicationLoadBalancerProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableApplicationLoadBalancerProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApplicationLoadBalancerProperties(val *ApplicationLoadBalancerProperties) *NullableApplicationLoadBalancerProperties {
	return &NullableApplicationLoadBalancerProperties{value: val, isSet: true}
}

func (v NullableApplicationLoadBalancerProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApplicationLoadBalancerProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


