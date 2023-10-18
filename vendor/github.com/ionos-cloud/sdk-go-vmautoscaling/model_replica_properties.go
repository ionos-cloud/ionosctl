/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ReplicaProperties struct for ReplicaProperties
type ReplicaProperties struct {
	AvailabilityZone *AvailabilityZone `json:"availabilityZone,omitempty"`
	// The total number of cores for the VMs.
	Cores     *int32     `json:"cores"`
	CpuFamily *CpuFamily `json:"cpuFamily,omitempty"`
	// The list of NICs associated with this replica.
	Nics *[]ReplicaNic `json:"nics,omitempty"`
	// The size of the memory for the VMs in MB. The size must be in multiples of 256 MB, with a minimum of 256 MB; if you set 'ramHotPlug=TRUE', you must use at least 1024 MB. If you set the RAM size to more than 240 GB, 'ramHotPlug=FALSE' is fixed.
	Ram *int32 `json:"ram"`
}

// NewReplicaProperties instantiates a new ReplicaProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicaProperties(cores int32, ram int32) *ReplicaProperties {
	this := ReplicaProperties{}

	this.Cores = &cores
	this.Ram = &ram

	return &this
}

// NewReplicaPropertiesWithDefaults instantiates a new ReplicaProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicaPropertiesWithDefaults() *ReplicaProperties {
	this := ReplicaProperties{}
	return &this
}

// GetAvailabilityZone returns the AvailabilityZone field value
// If the value is explicit nil, the zero value for AvailabilityZone will be returned
func (o *ReplicaProperties) GetAvailabilityZone() *AvailabilityZone {
	if o == nil {
		return nil
	}

	return o.AvailabilityZone

}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaProperties) GetAvailabilityZoneOk() (*AvailabilityZone, bool) {
	if o == nil {
		return nil, false
	}

	return o.AvailabilityZone, true
}

// SetAvailabilityZone sets field value
func (o *ReplicaProperties) SetAvailabilityZone(v AvailabilityZone) {

	o.AvailabilityZone = &v

}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *ReplicaProperties) HasAvailabilityZone() bool {
	if o != nil && o.AvailabilityZone != nil {
		return true
	}

	return false
}

// GetCores returns the Cores field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ReplicaProperties) GetCores() *int32 {
	if o == nil {
		return nil
	}

	return o.Cores

}

// GetCoresOk returns a tuple with the Cores field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaProperties) GetCoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Cores, true
}

// SetCores sets field value
func (o *ReplicaProperties) SetCores(v int32) {

	o.Cores = &v

}

// HasCores returns a boolean if a field has been set.
func (o *ReplicaProperties) HasCores() bool {
	if o != nil && o.Cores != nil {
		return true
	}

	return false
}

// GetCpuFamily returns the CpuFamily field value
// If the value is explicit nil, the zero value for CpuFamily will be returned
func (o *ReplicaProperties) GetCpuFamily() *CpuFamily {
	if o == nil {
		return nil
	}

	return o.CpuFamily

}

// GetCpuFamilyOk returns a tuple with the CpuFamily field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaProperties) GetCpuFamilyOk() (*CpuFamily, bool) {
	if o == nil {
		return nil, false
	}

	return o.CpuFamily, true
}

// SetCpuFamily sets field value
func (o *ReplicaProperties) SetCpuFamily(v CpuFamily) {

	o.CpuFamily = &v

}

// HasCpuFamily returns a boolean if a field has been set.
func (o *ReplicaProperties) HasCpuFamily() bool {
	if o != nil && o.CpuFamily != nil {
		return true
	}

	return false
}

// GetNics returns the Nics field value
// If the value is explicit nil, the zero value for []ReplicaNic will be returned
func (o *ReplicaProperties) GetNics() *[]ReplicaNic {
	if o == nil {
		return nil
	}

	return o.Nics

}

// GetNicsOk returns a tuple with the Nics field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaProperties) GetNicsOk() (*[]ReplicaNic, bool) {
	if o == nil {
		return nil, false
	}

	return o.Nics, true
}

// SetNics sets field value
func (o *ReplicaProperties) SetNics(v []ReplicaNic) {

	o.Nics = &v

}

// HasNics returns a boolean if a field has been set.
func (o *ReplicaProperties) HasNics() bool {
	if o != nil && o.Nics != nil {
		return true
	}

	return false
}

// GetRam returns the Ram field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ReplicaProperties) GetRam() *int32 {
	if o == nil {
		return nil
	}

	return o.Ram

}

// GetRamOk returns a tuple with the Ram field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaProperties) GetRamOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Ram, true
}

// SetRam sets field value
func (o *ReplicaProperties) SetRam(v int32) {

	o.Ram = &v

}

// HasRam returns a boolean if a field has been set.
func (o *ReplicaProperties) HasRam() bool {
	if o != nil && o.Ram != nil {
		return true
	}

	return false
}

func (o ReplicaProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["availabilityZone"] = o.AvailabilityZone

	if o.Cores != nil {
		toSerialize["cores"] = o.Cores
	}

	if o.CpuFamily != nil {
		toSerialize["cpuFamily"] = o.CpuFamily
	}

	if o.Nics != nil {
		toSerialize["nics"] = o.Nics
	}

	if o.Ram != nil {
		toSerialize["ram"] = o.Ram
	}

	return json.Marshal(toSerialize)
}

type NullableReplicaProperties struct {
	value *ReplicaProperties
	isSet bool
}

func (v NullableReplicaProperties) Get() *ReplicaProperties {
	return v.value
}

func (v *NullableReplicaProperties) Set(val *ReplicaProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicaProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicaProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicaProperties(val *ReplicaProperties) *NullableReplicaProperties {
	return &NullableReplicaProperties{value: val, isSet: true}
}

func (v NullableReplicaProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicaProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
