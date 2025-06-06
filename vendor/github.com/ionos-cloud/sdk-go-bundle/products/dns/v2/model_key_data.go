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

// checks if the KeyData type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KeyData{}

// KeyData Represents the separate components of the RDATA for a DNSKEY. The values must conform to the guidelines in [RFC-4034 Section 2.1](https://www.rfc-editor.org/rfc/rfc4034#section-2.1).
type KeyData struct {
	// Represents the key's metadata and usage information.
	Flags *int32 `json:"flags,omitempty"`
	// Represents the public key data in Base64 encoding.
	PubKey *string `json:"pubKey,omitempty"`
}

// NewKeyData instantiates a new KeyData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeyData() *KeyData {
	this := KeyData{}

	return &this
}

// NewKeyDataWithDefaults instantiates a new KeyData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeyDataWithDefaults() *KeyData {
	this := KeyData{}
	return &this
}

// GetFlags returns the Flags field value if set, zero value otherwise.
func (o *KeyData) GetFlags() int32 {
	if o == nil || IsNil(o.Flags) {
		var ret int32
		return ret
	}
	return *o.Flags
}

// GetFlagsOk returns a tuple with the Flags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeyData) GetFlagsOk() (*int32, bool) {
	if o == nil || IsNil(o.Flags) {
		return nil, false
	}
	return o.Flags, true
}

// HasFlags returns a boolean if a field has been set.
func (o *KeyData) HasFlags() bool {
	if o != nil && !IsNil(o.Flags) {
		return true
	}

	return false
}

// SetFlags gets a reference to the given int32 and assigns it to the Flags field.
func (o *KeyData) SetFlags(v int32) {
	o.Flags = &v
}

// GetPubKey returns the PubKey field value if set, zero value otherwise.
func (o *KeyData) GetPubKey() string {
	if o == nil || IsNil(o.PubKey) {
		var ret string
		return ret
	}
	return *o.PubKey
}

// GetPubKeyOk returns a tuple with the PubKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeyData) GetPubKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PubKey) {
		return nil, false
	}
	return o.PubKey, true
}

// HasPubKey returns a boolean if a field has been set.
func (o *KeyData) HasPubKey() bool {
	if o != nil && !IsNil(o.PubKey) {
		return true
	}

	return false
}

// SetPubKey gets a reference to the given string and assigns it to the PubKey field.
func (o *KeyData) SetPubKey(v string) {
	o.PubKey = &v
}

func (o KeyData) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Flags) {
		toSerialize["flags"] = o.Flags
	}
	if !IsNil(o.PubKey) {
		toSerialize["pubKey"] = o.PubKey
	}
	return toSerialize, nil
}

type NullableKeyData struct {
	value *KeyData
	isSet bool
}

func (v NullableKeyData) Get() *KeyData {
	return v.value
}

func (v *NullableKeyData) Set(val *KeyData) {
	v.value = val
	v.isSet = true
}

func (v NullableKeyData) IsSet() bool {
	return v.isSet
}

func (v *NullableKeyData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeyData(val *KeyData) *NullableKeyData {
	return &NullableKeyData{value: val, isSet: true}
}

func (v NullableKeyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeyData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
