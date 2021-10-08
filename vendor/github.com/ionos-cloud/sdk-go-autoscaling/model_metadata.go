/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional databse clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1-SDK.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"time"
)

// Metadata Metadata of the resource
type Metadata struct {
	// The ISO 8601 creation timestamp.
	CreatedDate *IonosTime
	CreatedBy *string `json:"created_by,omitempty"`
	CreatedByUserId *string `json:"created_by_user_id,omitempty"`
	// The ISO 8601 modified timestamp.
	LastModifiedDate *IonosTime
	LastModifiedBy *string `json:"last_modified_by,omitempty"`
	LastModifiedByUserId *string `json:"last_modified_by_user_id,omitempty"`
}



// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *Metadata) GetCreatedDate() *time.Time {
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
func (o *Metadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *Metadata) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}


}

// HasCreatedDate returns a boolean if a field has been set.
func (o *Metadata) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}



// GetCreatedBy returns the CreatedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCreatedBy() *string {
	if o == nil {
		return nil
	}


	return o.CreatedBy

}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *Metadata) SetCreatedBy(v string) {


	o.CreatedBy = &v

}

// HasCreatedBy returns a boolean if a field has been set.
func (o *Metadata) HasCreatedBy() bool {
	if o != nil && o.CreatedBy != nil {
		return true
	}

	return false
}



// GetCreatedByUserId returns the CreatedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCreatedByUserId() *string {
	if o == nil {
		return nil
	}


	return o.CreatedByUserId

}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.CreatedByUserId, true
}

// SetCreatedByUserId sets field value
func (o *Metadata) SetCreatedByUserId(v string) {


	o.CreatedByUserId = &v

}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasCreatedByUserId() bool {
	if o != nil && o.CreatedByUserId != nil {
		return true
	}

	return false
}



// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *Metadata) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time


}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *Metadata) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}


}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}



// GetLastModifiedBy returns the LastModifiedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetLastModifiedBy() *string {
	if o == nil {
		return nil
	}


	return o.LastModifiedBy

}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LastModifiedBy, true
}

// SetLastModifiedBy sets field value
func (o *Metadata) SetLastModifiedBy(v string) {


	o.LastModifiedBy = &v

}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedBy() bool {
	if o != nil && o.LastModifiedBy != nil {
		return true
	}

	return false
}



// GetLastModifiedByUserId returns the LastModifiedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetLastModifiedByUserId() *string {
	if o == nil {
		return nil
	}


	return o.LastModifiedByUserId

}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LastModifiedByUserId, true
}

// SetLastModifiedByUserId sets field value
func (o *Metadata) SetLastModifiedByUserId(v string) {


	o.LastModifiedByUserId = &v

}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedByUserId() bool {
	if o != nil && o.LastModifiedByUserId != nil {
		return true
	}

	return false
}


func (o Metadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.CreatedDate != nil {
		toSerialize["created_date"] = o.CreatedDate
	}
	

	if o.CreatedBy != nil {
		toSerialize["created_by"] = o.CreatedBy
	}
	

	if o.CreatedByUserId != nil {
		toSerialize["created_by_user_id"] = o.CreatedByUserId
	}
	

	if o.LastModifiedDate != nil {
		toSerialize["last_modified_date"] = o.LastModifiedDate
	}
	

	if o.LastModifiedBy != nil {
		toSerialize["last_modified_by"] = o.LastModifiedBy
	}
	

	if o.LastModifiedByUserId != nil {
		toSerialize["last_modified_by_user_id"] = o.LastModifiedByUserId
	}
	
	return json.Marshal(toSerialize)
}

type NullableMetadata struct {
	value *Metadata
	isSet bool
}

func (v NullableMetadata) Get() *Metadata {
	return v.value
}

func (v *NullableMetadata) Set(val *Metadata) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadata(val *Metadata) *NullableMetadata {
	return &NullableMetadata{value: val, isSet: true}
}

func (v NullableMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


