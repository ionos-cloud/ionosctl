/*
 * IONOS DBaaS PostgreSQL REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional PostgreSQL database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package psql

import (
	"encoding/json"

	"fmt"
)

// PoolMode Represents different modes of connection pooling for the connection pooler.
type PoolMode string

// List of PoolMode
const (
	POOLMODE_TRANSACTION PoolMode = "transaction"
	POOLMODE_SESSION     PoolMode = "session"
)

func (v *PoolMode) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := PoolMode(value)
	for _, existing := range []PoolMode{"transaction", "session"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid PoolMode", value)
}

// Ptr returns reference to PoolMode value
func (v PoolMode) Ptr() *PoolMode {
	return &v
}

type NullablePoolMode struct {
	value *PoolMode
	isSet bool
}

func (v NullablePoolMode) Get() *PoolMode {
	return v.value
}

func (v *NullablePoolMode) Set(val *PoolMode) {
	v.value = val
	v.isSet = true
}

func (v NullablePoolMode) IsSet() bool {
	return v.isSet
}

func (v *NullablePoolMode) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePoolMode(val *PoolMode) *NullablePoolMode {
	return &NullablePoolMode{value: val, isSet: true}
}

func (v NullablePoolMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePoolMode) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
