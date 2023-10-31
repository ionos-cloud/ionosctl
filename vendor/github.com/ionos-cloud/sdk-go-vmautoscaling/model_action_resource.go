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

// ActionResource struct for ActionResource
type ActionResource struct {
	// The unique resource identifier.
	Id *string `json:"id"`
	// The resource type.
	Type *string `json:"type,omitempty"`
	// The absolute URL to the resource's representation.
	Href *string `json:"href,omitempty"`
}

// NewActionResource instantiates a new ActionResource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewActionResource(id string) *ActionResource {
	this := ActionResource{}

	this.Id = &id

	return &this
}

// NewActionResourceWithDefaults instantiates a new ActionResource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewActionResourceWithDefaults() *ActionResource {
	this := ActionResource{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ActionResource) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ActionResource) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *ActionResource) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ActionResource) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ActionResource) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ActionResource) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *ActionResource) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ActionResource) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ActionResource) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ActionResource) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *ActionResource) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *ActionResource) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

func (o ActionResource) MarshalJSON() ([]byte, error) {
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

	return json.Marshal(toSerialize)
}

type NullableActionResource struct {
	value *ActionResource
	isSet bool
}

func (v NullableActionResource) Get() *ActionResource {
	return v.value
}

func (v *NullableActionResource) Set(val *ActionResource) {
	v.value = val
	v.isSet = true
}

func (v NullableActionResource) IsSet() bool {
	return v.isSet
}

func (v *NullableActionResource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableActionResource(val *ActionResource) *NullableActionResource {
	return &NullableActionResource{value: val, isSet: true}
}

func (v NullableActionResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableActionResource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}