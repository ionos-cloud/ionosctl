/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.15.4
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// DnssecKey struct for DnssecKey
type DnssecKey struct {
	KeyTag *int32 `json:"keyTag,omitempty"`
	// A string that denotes the digest algorithm. This value must conform to the guidelines in [RFC-8624 Section 3.3](https://datatracker.ietf.org/doc/html/rfc8624#section-3.3).
	DigestAlgorithmMnemonic *string  `json:"digestAlgorithmMnemonic,omitempty"`
	Digest                  *string  `json:"digest,omitempty"`
	KeyData                 *KeyData `json:"keyData,omitempty"`
	// Represents the composed value of the The RDATA for a DNSKEY. The format should be the following: Flags | Protocol | Algorithm | Public Key The values must conform to the guidelines in [RFC-4034 Section 2.1](https://www.rfc-editor.org/rfc/rfc4034#section-2.1).
	ComposedKeyData *string `json:"composedKeyData,omitempty"`
}

// NewDnssecKey instantiates a new DnssecKey object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnssecKey() *DnssecKey {
	this := DnssecKey{}

	return &this
}

// NewDnssecKeyWithDefaults instantiates a new DnssecKey object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnssecKeyWithDefaults() *DnssecKey {
	this := DnssecKey{}
	return &this
}

// GetKeyTag returns the KeyTag field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *DnssecKey) GetKeyTag() *int32 {
	if o == nil {
		return nil
	}

	return o.KeyTag

}

// GetKeyTagOk returns a tuple with the KeyTag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnssecKey) GetKeyTagOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.KeyTag, true
}

// SetKeyTag sets field value
func (o *DnssecKey) SetKeyTag(v int32) {

	o.KeyTag = &v

}

// HasKeyTag returns a boolean if a field has been set.
func (o *DnssecKey) HasKeyTag() bool {
	if o != nil && o.KeyTag != nil {
		return true
	}

	return false
}

// GetDigestAlgorithmMnemonic returns the DigestAlgorithmMnemonic field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DnssecKey) GetDigestAlgorithmMnemonic() *string {
	if o == nil {
		return nil
	}

	return o.DigestAlgorithmMnemonic

}

// GetDigestAlgorithmMnemonicOk returns a tuple with the DigestAlgorithmMnemonic field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnssecKey) GetDigestAlgorithmMnemonicOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.DigestAlgorithmMnemonic, true
}

// SetDigestAlgorithmMnemonic sets field value
func (o *DnssecKey) SetDigestAlgorithmMnemonic(v string) {

	o.DigestAlgorithmMnemonic = &v

}

// HasDigestAlgorithmMnemonic returns a boolean if a field has been set.
func (o *DnssecKey) HasDigestAlgorithmMnemonic() bool {
	if o != nil && o.DigestAlgorithmMnemonic != nil {
		return true
	}

	return false
}

// GetDigest returns the Digest field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DnssecKey) GetDigest() *string {
	if o == nil {
		return nil
	}

	return o.Digest

}

// GetDigestOk returns a tuple with the Digest field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnssecKey) GetDigestOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Digest, true
}

// SetDigest sets field value
func (o *DnssecKey) SetDigest(v string) {

	o.Digest = &v

}

// HasDigest returns a boolean if a field has been set.
func (o *DnssecKey) HasDigest() bool {
	if o != nil && o.Digest != nil {
		return true
	}

	return false
}

// GetKeyData returns the KeyData field value
// If the value is explicit nil, the zero value for KeyData will be returned
func (o *DnssecKey) GetKeyData() *KeyData {
	if o == nil {
		return nil
	}

	return o.KeyData

}

// GetKeyDataOk returns a tuple with the KeyData field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnssecKey) GetKeyDataOk() (*KeyData, bool) {
	if o == nil {
		return nil, false
	}

	return o.KeyData, true
}

// SetKeyData sets field value
func (o *DnssecKey) SetKeyData(v KeyData) {

	o.KeyData = &v

}

// HasKeyData returns a boolean if a field has been set.
func (o *DnssecKey) HasKeyData() bool {
	if o != nil && o.KeyData != nil {
		return true
	}

	return false
}

// GetComposedKeyData returns the ComposedKeyData field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DnssecKey) GetComposedKeyData() *string {
	if o == nil {
		return nil
	}

	return o.ComposedKeyData

}

// GetComposedKeyDataOk returns a tuple with the ComposedKeyData field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnssecKey) GetComposedKeyDataOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ComposedKeyData, true
}

// SetComposedKeyData sets field value
func (o *DnssecKey) SetComposedKeyData(v string) {

	o.ComposedKeyData = &v

}

// HasComposedKeyData returns a boolean if a field has been set.
func (o *DnssecKey) HasComposedKeyData() bool {
	if o != nil && o.ComposedKeyData != nil {
		return true
	}

	return false
}

func (o DnssecKey) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.KeyTag != nil {
		toSerialize["keyTag"] = o.KeyTag
	}

	if o.DigestAlgorithmMnemonic != nil {
		toSerialize["digestAlgorithmMnemonic"] = o.DigestAlgorithmMnemonic
	}

	if o.Digest != nil {
		toSerialize["digest"] = o.Digest
	}

	if o.KeyData != nil {
		toSerialize["keyData"] = o.KeyData
	}

	if o.ComposedKeyData != nil {
		toSerialize["composedKeyData"] = o.ComposedKeyData
	}

	return json.Marshal(toSerialize)
}

type NullableDnssecKey struct {
	value *DnssecKey
	isSet bool
}

func (v NullableDnssecKey) Get() *DnssecKey {
	return v.value
}

func (v *NullableDnssecKey) Set(val *DnssecKey) {
	v.value = val
	v.isSet = true
}

func (v NullableDnssecKey) IsSet() bool {
	return v.isSet
}

func (v *NullableDnssecKey) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnssecKey(val *DnssecKey) *NullableDnssecKey {
	return &NullableDnssecKey{value: val, isSet: true}
}

func (v NullableDnssecKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnssecKey) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
