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

// RequestStatusMetadata struct for RequestStatusMetadata
type RequestStatusMetadata struct {
	Status *string `json:"status,omitempty"`
	Message *string `json:"message,omitempty"`
	// Resource's Entity Tag as defined in http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.11 . Entity Tag is also added as an 'ETag response header to requests which don't use 'depth' parameter. 
	Etag *string `json:"etag,omitempty"`
	Targets *[]RequestTarget `json:"targets,omitempty"`
}



// GetStatus returns the Status field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RequestStatusMetadata) GetStatus() *string {
	if o == nil {
		return nil
	}


	return o.Status

}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestStatusMetadata) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Status, true
}

// SetStatus sets field value
func (o *RequestStatusMetadata) SetStatus(v string) {


	o.Status = &v

}

// HasStatus returns a boolean if a field has been set.
func (o *RequestStatusMetadata) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}



// GetMessage returns the Message field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RequestStatusMetadata) GetMessage() *string {
	if o == nil {
		return nil
	}


	return o.Message

}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestStatusMetadata) GetMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Message, true
}

// SetMessage sets field value
func (o *RequestStatusMetadata) SetMessage(v string) {


	o.Message = &v

}

// HasMessage returns a boolean if a field has been set.
func (o *RequestStatusMetadata) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}



// GetEtag returns the Etag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RequestStatusMetadata) GetEtag() *string {
	if o == nil {
		return nil
	}


	return o.Etag

}

// GetEtagOk returns a tuple with the Etag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestStatusMetadata) GetEtagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Etag, true
}

// SetEtag sets field value
func (o *RequestStatusMetadata) SetEtag(v string) {


	o.Etag = &v

}

// HasEtag returns a boolean if a field has been set.
func (o *RequestStatusMetadata) HasEtag() bool {
	if o != nil && o.Etag != nil {
		return true
	}

	return false
}



// GetTargets returns the Targets field value
// If the value is explicit nil, the zero value for []RequestTarget will be returned
func (o *RequestStatusMetadata) GetTargets() *[]RequestTarget {
	if o == nil {
		return nil
	}


	return o.Targets

}

// GetTargetsOk returns a tuple with the Targets field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestStatusMetadata) GetTargetsOk() (*[]RequestTarget, bool) {
	if o == nil {
		return nil, false
	}


	return o.Targets, true
}

// SetTargets sets field value
func (o *RequestStatusMetadata) SetTargets(v []RequestTarget) {


	o.Targets = &v

}

// HasTargets returns a boolean if a field has been set.
func (o *RequestStatusMetadata) HasTargets() bool {
	if o != nil && o.Targets != nil {
		return true
	}

	return false
}


func (o RequestStatusMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	

	if o.Message != nil {
		toSerialize["message"] = o.Message
	}
	

	if o.Etag != nil {
		toSerialize["etag"] = o.Etag
	}
	

	if o.Targets != nil {
		toSerialize["targets"] = o.Targets
	}
	
	return json.Marshal(toSerialize)
}

type NullableRequestStatusMetadata struct {
	value *RequestStatusMetadata
	isSet bool
}

func (v NullableRequestStatusMetadata) Get() *RequestStatusMetadata {
	return v.value
}

func (v *NullableRequestStatusMetadata) Set(val *RequestStatusMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestStatusMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestStatusMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestStatusMetadata(val *RequestStatusMetadata) *NullableRequestStatusMetadata {
	return &NullableRequestStatusMetadata{value: val, isSet: true}
}

func (v NullableRequestStatusMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestStatusMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


