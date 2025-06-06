/*
 * IONOS Cloud - API Gateway
 *
 * API Gateway is an application that acts as a \"front door\" for backend services and APIs, handling client requests and routing them to the appropriate backend.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apigateway

import (
	"encoding/json"
)

// checks if the GatewayReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GatewayReadListAllOf{}

// GatewayReadListAllOf struct for GatewayReadListAllOf
type GatewayReadListAllOf struct {
	// ID of the Gateway.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the Gateway.
	Href string `json:"href"`
	// The list of Gateway resources.
	Items []GatewayRead `json:"items,omitempty"`
}

// NewGatewayReadListAllOf instantiates a new GatewayReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGatewayReadListAllOf(id string, type_ string, href string) *GatewayReadListAllOf {
	this := GatewayReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewGatewayReadListAllOfWithDefaults instantiates a new GatewayReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGatewayReadListAllOfWithDefaults() *GatewayReadListAllOf {
	this := GatewayReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *GatewayReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GatewayReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GatewayReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *GatewayReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GatewayReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GatewayReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *GatewayReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *GatewayReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *GatewayReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *GatewayReadListAllOf) GetItems() []GatewayRead {
	if o == nil || IsNil(o.Items) {
		var ret []GatewayRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GatewayReadListAllOf) GetItemsOk() ([]GatewayRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *GatewayReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []GatewayRead and assigns it to the Items field.
func (o *GatewayReadListAllOf) SetItems(v []GatewayRead) {
	o.Items = v
}

func (o GatewayReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableGatewayReadListAllOf struct {
	value *GatewayReadListAllOf
	isSet bool
}

func (v NullableGatewayReadListAllOf) Get() *GatewayReadListAllOf {
	return v.value
}

func (v *NullableGatewayReadListAllOf) Set(val *GatewayReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableGatewayReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableGatewayReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGatewayReadListAllOf(val *GatewayReadListAllOf) *NullableGatewayReadListAllOf {
	return &NullableGatewayReadListAllOf{value: val, isSet: true}
}

func (v NullableGatewayReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGatewayReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
