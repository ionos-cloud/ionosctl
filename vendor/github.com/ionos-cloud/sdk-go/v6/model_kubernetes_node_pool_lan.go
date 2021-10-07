/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesNodePoolLan struct for KubernetesNodePoolLan
type KubernetesNodePoolLan struct {
	// The LAN ID of an existing LAN at the related datacenter
	Id *int32 `json:"id"`
	// Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP
	Dhcp *bool `json:"dhcp,omitempty"`
	// array of additional LANs attached to worker nodes
	Routes *[]KubernetesNodePoolLanRoutes `json:"routes,omitempty"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *KubernetesNodePoolLan) GetId() *int32 {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolLan) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *KubernetesNodePoolLan) SetId(v int32) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *KubernetesNodePoolLan) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetDhcp returns the Dhcp field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *KubernetesNodePoolLan) GetDhcp() *bool {
	if o == nil {
		return nil
	}


	return o.Dhcp

}

// GetDhcpOk returns a tuple with the Dhcp field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolLan) GetDhcpOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}


	return o.Dhcp, true
}

// SetDhcp sets field value
func (o *KubernetesNodePoolLan) SetDhcp(v bool) {


	o.Dhcp = &v

}

// HasDhcp returns a boolean if a field has been set.
func (o *KubernetesNodePoolLan) HasDhcp() bool {
	if o != nil && o.Dhcp != nil {
		return true
	}

	return false
}



// GetRoutes returns the Routes field value
// If the value is explicit nil, the zero value for []KubernetesNodePoolLanRoutes will be returned
func (o *KubernetesNodePoolLan) GetRoutes() *[]KubernetesNodePoolLanRoutes {
	if o == nil {
		return nil
	}


	return o.Routes

}

// GetRoutesOk returns a tuple with the Routes field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolLan) GetRoutesOk() (*[]KubernetesNodePoolLanRoutes, bool) {
	if o == nil {
		return nil, false
	}


	return o.Routes, true
}

// SetRoutes sets field value
func (o *KubernetesNodePoolLan) SetRoutes(v []KubernetesNodePoolLanRoutes) {


	o.Routes = &v

}

// HasRoutes returns a boolean if a field has been set.
func (o *KubernetesNodePoolLan) HasRoutes() bool {
	if o != nil && o.Routes != nil {
		return true
	}

	return false
}


func (o KubernetesNodePoolLan) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	

	if o.Dhcp != nil {
		toSerialize["dhcp"] = o.Dhcp
	}
	

	if o.Routes != nil {
		toSerialize["routes"] = o.Routes
	}
	
	return json.Marshal(toSerialize)
}

type NullableKubernetesNodePoolLan struct {
	value *KubernetesNodePoolLan
	isSet bool
}

func (v NullableKubernetesNodePoolLan) Get() *KubernetesNodePoolLan {
	return v.value
}

func (v *NullableKubernetesNodePoolLan) Set(val *KubernetesNodePoolLan) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodePoolLan) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodePoolLan) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodePoolLan(val *KubernetesNodePoolLan) *NullableKubernetesNodePoolLan {
	return &NullableKubernetesNodePoolLan{value: val, isSet: true}
}

func (v NullableKubernetesNodePoolLan) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodePoolLan) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


