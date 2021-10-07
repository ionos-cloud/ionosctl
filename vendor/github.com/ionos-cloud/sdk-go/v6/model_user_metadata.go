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
	"time"
)

// UserMetadata struct for UserMetadata
type UserMetadata struct {
	// Resource's Entity Tag as defined in http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.11 . Entity Tag is also added as an 'ETag response header to requests which don't use 'depth' parameter. 
	Etag *string `json:"etag,omitempty"`
	// time of creation of the user
	CreatedDate *IonosTime
	// time of last login by the user
	LastLogin *IonosTime
}



// GetEtag returns the Etag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *UserMetadata) GetEtag() *string {
	if o == nil {
		return nil
	}


	return o.Etag

}

// GetEtagOk returns a tuple with the Etag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserMetadata) GetEtagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Etag, true
}

// SetEtag sets field value
func (o *UserMetadata) SetEtag(v string) {


	o.Etag = &v

}

// HasEtag returns a boolean if a field has been set.
func (o *UserMetadata) HasEtag() bool {
	if o != nil && o.Etag != nil {
		return true
	}

	return false
}



// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *UserMetadata) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time


}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserMetadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *UserMetadata) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}


}

// HasCreatedDate returns a boolean if a field has been set.
func (o *UserMetadata) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}



// GetLastLogin returns the LastLogin field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *UserMetadata) GetLastLogin() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastLogin == nil {
		return nil
	}
	return &o.LastLogin.Time


}

// GetLastLoginOk returns a tuple with the LastLogin field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserMetadata) GetLastLoginOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastLogin == nil {
		return nil, false
	}
	return &o.LastLogin.Time, true

}

// SetLastLogin sets field value
func (o *UserMetadata) SetLastLogin(v time.Time) {

	o.LastLogin = &IonosTime{v}


}

// HasLastLogin returns a boolean if a field has been set.
func (o *UserMetadata) HasLastLogin() bool {
	if o != nil && o.LastLogin != nil {
		return true
	}

	return false
}


func (o UserMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Etag != nil {
		toSerialize["etag"] = o.Etag
	}
	

	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}
	

	if o.LastLogin != nil {
		toSerialize["lastLogin"] = o.LastLogin
	}
	
	return json.Marshal(toSerialize)
}

type NullableUserMetadata struct {
	value *UserMetadata
	isSet bool
}

func (v NullableUserMetadata) Get() *UserMetadata {
	return v.value
}

func (v *NullableUserMetadata) Set(val *UserMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableUserMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableUserMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserMetadata(val *UserMetadata) *NullableUserMetadata {
	return &NullableUserMetadata{value: val, isSet: true}
}

func (v NullableUserMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


