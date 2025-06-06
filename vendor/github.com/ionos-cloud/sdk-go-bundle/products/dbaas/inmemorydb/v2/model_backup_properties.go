/*
 * In-Memory DB API
 *
 * API description for the IONOS In-Memory DB
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package inmemorydb

import (
	"encoding/json"
)

// checks if the BackupProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BackupProperties{}

// BackupProperties Properties configuring the backup of the replicaset.
type BackupProperties struct {
	// The S3 location where the backups will be stored.
	Location *string `json:"location,omitempty"`
}

// NewBackupProperties instantiates a new BackupProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupProperties() *BackupProperties {
	this := BackupProperties{}

	return &this
}

// NewBackupPropertiesWithDefaults instantiates a new BackupProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupPropertiesWithDefaults() *BackupProperties {
	this := BackupProperties{}
	return &this
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *BackupProperties) GetLocation() string {
	if o == nil || IsNil(o.Location) {
		var ret string
		return ret
	}
	return *o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupProperties) GetLocationOk() (*string, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *BackupProperties) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given string and assigns it to the Location field.
func (o *BackupProperties) SetLocation(v string) {
	o.Location = &v
}

func (o BackupProperties) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BackupProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Location) {
		toSerialize["location"] = o.Location
	}
	return toSerialize, nil
}

type NullableBackupProperties struct {
	value *BackupProperties
	isSet bool
}

func (v NullableBackupProperties) Get() *BackupProperties {
	return v.value
}

func (v *NullableBackupProperties) Set(val *BackupProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableBackupProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableBackupProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBackupProperties(val *BackupProperties) *NullableBackupProperties {
	return &NullableBackupProperties{value: val, isSet: true}
}

func (v NullableBackupProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBackupProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
