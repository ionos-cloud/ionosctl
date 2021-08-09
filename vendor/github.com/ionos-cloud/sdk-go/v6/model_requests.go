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

// Requests struct for Requests
type Requests struct {
	// The resource's unique identifier
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	// Array of items in that collection
	Items *[]Request `json:"items,omitempty"`
	// the offset specified in the request (or, if none was specified, the default offset of 0)
	Offset *float32 `json:"offset"`
	// the limit specified in the request (or, if none was specified, the default limit of 0)
	Limit *float32 `json:"limit"`
	Links *PaginationLinks `json:"_links"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Requests) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *Requests) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *Requests) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for Type will be returned
func (o *Requests) GetType() *Type {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *Requests) SetType(v Type) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Requests) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Requests) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *Requests) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *Requests) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []Request will be returned
func (o *Requests) GetItems() *[]Request {
	if o == nil {
		return nil
	}


	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetItemsOk() (*[]Request, bool) {
	if o == nil {
		return nil, false
	}


	return o.Items, true
}

// SetItems sets field value
func (o *Requests) SetItems(v []Request) {


	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *Requests) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}



// GetOffset returns the Offset field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Requests) GetOffset() *float32 {
	if o == nil {
		return nil
	}


	return o.Offset

}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetOffsetOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Offset, true
}

// SetOffset sets field value
func (o *Requests) SetOffset(v float32) {


	o.Offset = &v

}

// HasOffset returns a boolean if a field has been set.
func (o *Requests) HasOffset() bool {
	if o != nil && o.Offset != nil {
		return true
	}

	return false
}



// GetLimit returns the Limit field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Requests) GetLimit() *float32 {
	if o == nil {
		return nil
	}


	return o.Limit

}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetLimitOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Limit, true
}

// SetLimit sets field value
func (o *Requests) SetLimit(v float32) {


	o.Limit = &v

}

// HasLimit returns a boolean if a field has been set.
func (o *Requests) HasLimit() bool {
	if o != nil && o.Limit != nil {
		return true
	}

	return false
}



// GetLinks returns the Links field value
// If the value is explicit nil, the zero value for PaginationLinks will be returned
func (o *Requests) GetLinks() *PaginationLinks {
	if o == nil {
		return nil
	}


	return o.Links

}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Requests) GetLinksOk() (*PaginationLinks, bool) {
	if o == nil {
		return nil, false
	}


	return o.Links, true
}

// SetLinks sets field value
func (o *Requests) SetLinks(v PaginationLinks) {


	o.Links = &v

}

// HasLinks returns a boolean if a field has been set.
func (o *Requests) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}


func (o Requests) MarshalJSON() ([]byte, error) {
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
	

	if o.Offset != nil {
		toSerialize["offset"] = o.Offset
	}
	

	if o.Limit != nil {
		toSerialize["limit"] = o.Limit
	}
	

	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}
	
	return json.Marshal(toSerialize)
}

type NullableRequests struct {
	value *Requests
	isSet bool
}

func (v NullableRequests) Get() *Requests {
	return v.value
}

func (v *NullableRequests) Set(val *Requests) {
	v.value = val
	v.isSet = true
}

func (v NullableRequests) IsSet() bool {
	return v.isSet
}

func (v *NullableRequests) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequests(val *Requests) *NullableRequests {
	return &NullableRequests{value: val, isSet: true}
}

func (v NullableRequests) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequests) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


