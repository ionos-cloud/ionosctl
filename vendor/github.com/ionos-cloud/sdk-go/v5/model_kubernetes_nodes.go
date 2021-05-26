/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesNodes struct for KubernetesNodes
type KubernetesNodes struct {
	// Unique representation for Kubernetes Node Pool as a collection on a resource.
	Id *string `json:"id,omitempty"`
	// The type of resource within a collection
	Type *string `json:"type,omitempty"`
	// URL to the collection representation (absolute path)
	Href *string `json:"href,omitempty"`
	// Array of items in that collection
	Items *[]KubernetesNode `json:"items,omitempty"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodes) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodes) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *KubernetesNodes) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *KubernetesNodes) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodes) GetType() *string {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodes) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *KubernetesNodes) SetType(v string) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *KubernetesNodes) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodes) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodes) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *KubernetesNodes) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *KubernetesNodes) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []KubernetesNode will be returned
func (o *KubernetesNodes) GetItems() *[]KubernetesNode {
	if o == nil {
		return nil
	}


	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodes) GetItemsOk() (*[]KubernetesNode, bool) {
	if o == nil {
		return nil, false
	}


	return o.Items, true
}

// SetItems sets field value
func (o *KubernetesNodes) SetItems(v []KubernetesNode) {


	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *KubernetesNodes) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}


func (o KubernetesNodes) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}
	
	return json.Marshal(toSerialize)
}

type NullableKubernetesNodes struct {
	value *KubernetesNodes
	isSet bool
}

func (v NullableKubernetesNodes) Get() *KubernetesNodes {
	return v.value
}

func (v *NullableKubernetesNodes) Set(val *KubernetesNodes) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodes) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodes(val *KubernetesNodes) *NullableKubernetesNodes {
	return &NullableKubernetesNodes{value: val, isSet: true}
}

func (v NullableKubernetesNodes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


