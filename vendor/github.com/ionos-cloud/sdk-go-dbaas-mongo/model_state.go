/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.  MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
)

// State The current status reported by the cluster. * **AVAILABLE** Resources for this cluster exist and are healthy. * **BUSY** Resources for this cluster are being created or updated. * **DESTROYING** Delete cluster command was issued, the cluster is being deleted. * **FAILED** Failed to get the cluster status. * **UNKNOWN** The state is unknown.
type State string

// List of State
const (
	STATE_AVAILABLE  State = "AVAILABLE"
	STATE_BUSY       State = "BUSY"
	STATE_DESTROYING State = "DESTROYING"
	STATE_FAILED     State = "FAILED"
	STATE_UNKNOWN    State = "UNKNOWN"
)

func (v *State) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := State(value)
	for _, existing := range []State{"AVAILABLE", "BUSY", "DESTROYING", "FAILED", "UNKNOWN"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid State", value)
}

// Ptr returns reference to State value
func (v State) Ptr() *State {
	return &v
}

type NullableState struct {
	value *State
	isSet bool
}

func (v NullableState) Get() *State {
	return v.value
}

func (v *NullableState) Set(val *State) {
	v.value = val
	v.isSet = true
}

func (v NullableState) IsSet() bool {
	return v.isSet
}

func (v *NullableState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableState(val *State) *NullableState {
	return &NullableState{value: val, isSet: true}
}

func (v NullableState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
