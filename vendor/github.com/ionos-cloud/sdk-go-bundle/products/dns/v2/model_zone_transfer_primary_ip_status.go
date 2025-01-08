/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.16.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the ZoneTransferPrimaryIpStatus type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ZoneTransferPrimaryIpStatus{}

// ZoneTransferPrimaryIpStatus struct for ZoneTransferPrimaryIpStatus
type ZoneTransferPrimaryIpStatus struct {
	// one single IP from the primaryIps field for secondary zones
	PrimaryIp string `json:"primaryIp"`
	// Human readable status of the zone transfer status for the IP
	Status string `json:"status"`
	// Human readable explanation of the error when status is not ok
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// NewZoneTransferPrimaryIpStatus instantiates a new ZoneTransferPrimaryIpStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewZoneTransferPrimaryIpStatus(primaryIp string, status string) *ZoneTransferPrimaryIpStatus {
	this := ZoneTransferPrimaryIpStatus{}

	this.PrimaryIp = primaryIp
	this.Status = status

	return &this
}

// NewZoneTransferPrimaryIpStatusWithDefaults instantiates a new ZoneTransferPrimaryIpStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewZoneTransferPrimaryIpStatusWithDefaults() *ZoneTransferPrimaryIpStatus {
	this := ZoneTransferPrimaryIpStatus{}
	return &this
}

// GetPrimaryIp returns the PrimaryIp field value
func (o *ZoneTransferPrimaryIpStatus) GetPrimaryIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PrimaryIp
}

// GetPrimaryIpOk returns a tuple with the PrimaryIp field value
// and a boolean to check if the value has been set.
func (o *ZoneTransferPrimaryIpStatus) GetPrimaryIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PrimaryIp, true
}

// SetPrimaryIp sets field value
func (o *ZoneTransferPrimaryIpStatus) SetPrimaryIp(v string) {
	o.PrimaryIp = v
}

// GetStatus returns the Status field value
func (o *ZoneTransferPrimaryIpStatus) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *ZoneTransferPrimaryIpStatus) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *ZoneTransferPrimaryIpStatus) SetStatus(v string) {
	o.Status = v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise.
func (o *ZoneTransferPrimaryIpStatus) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ZoneTransferPrimaryIpStatus) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}
	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *ZoneTransferPrimaryIpStatus) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *ZoneTransferPrimaryIpStatus) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

func (o ZoneTransferPrimaryIpStatus) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["primaryIp"] = o.PrimaryIp
	toSerialize["status"] = o.Status
	if !IsNil(o.ErrorMessage) {
		toSerialize["errorMessage"] = o.ErrorMessage
	}
	return toSerialize, nil
}

type NullableZoneTransferPrimaryIpStatus struct {
	value *ZoneTransferPrimaryIpStatus
	isSet bool
}

func (v NullableZoneTransferPrimaryIpStatus) Get() *ZoneTransferPrimaryIpStatus {
	return v.value
}

func (v *NullableZoneTransferPrimaryIpStatus) Set(val *ZoneTransferPrimaryIpStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableZoneTransferPrimaryIpStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableZoneTransferPrimaryIpStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableZoneTransferPrimaryIpStatus(val *ZoneTransferPrimaryIpStatus) *NullableZoneTransferPrimaryIpStatus {
	return &NullableZoneTransferPrimaryIpStatus{value: val, isSet: true}
}

func (v NullableZoneTransferPrimaryIpStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableZoneTransferPrimaryIpStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
