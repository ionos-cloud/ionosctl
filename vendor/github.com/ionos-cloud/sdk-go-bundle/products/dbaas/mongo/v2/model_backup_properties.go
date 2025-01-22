/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.  MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mongo

import (
	"encoding/json"
)

// checks if the BackupProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BackupProperties{}

// BackupProperties Backup related properties.
type BackupProperties struct {
	// Number of hours between snapshots.
	SnapshotIntervalHours *int32 `json:"snapshotIntervalHours,omitempty"`
	// Number of hours in the past for which a point-in-time snapshot can be created.
	PointInTimeWindowHours *int32                     `json:"pointInTimeWindowHours,omitempty"`
	BackupRetention        *BackupRetentionProperties `json:"backupRetention,omitempty"`
	// The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster.
	Location *string `json:"location,omitempty"`
}

// NewBackupProperties instantiates a new BackupProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupProperties() *BackupProperties {
	this := BackupProperties{}

	var snapshotIntervalHours int32 = 24
	this.SnapshotIntervalHours = &snapshotIntervalHours

	return &this
}

// NewBackupPropertiesWithDefaults instantiates a new BackupProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupPropertiesWithDefaults() *BackupProperties {
	this := BackupProperties{}
	var snapshotIntervalHours int32 = 24
	this.SnapshotIntervalHours = &snapshotIntervalHours
	return &this
}

// GetSnapshotIntervalHours returns the SnapshotIntervalHours field value if set, zero value otherwise.
func (o *BackupProperties) GetSnapshotIntervalHours() int32 {
	if o == nil || IsNil(o.SnapshotIntervalHours) {
		var ret int32
		return ret
	}
	return *o.SnapshotIntervalHours
}

// GetSnapshotIntervalHoursOk returns a tuple with the SnapshotIntervalHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupProperties) GetSnapshotIntervalHoursOk() (*int32, bool) {
	if o == nil || IsNil(o.SnapshotIntervalHours) {
		return nil, false
	}
	return o.SnapshotIntervalHours, true
}

// HasSnapshotIntervalHours returns a boolean if a field has been set.
func (o *BackupProperties) HasSnapshotIntervalHours() bool {
	if o != nil && !IsNil(o.SnapshotIntervalHours) {
		return true
	}

	return false
}

// SetSnapshotIntervalHours gets a reference to the given int32 and assigns it to the SnapshotIntervalHours field.
func (o *BackupProperties) SetSnapshotIntervalHours(v int32) {
	o.SnapshotIntervalHours = &v
}

// GetPointInTimeWindowHours returns the PointInTimeWindowHours field value if set, zero value otherwise.
func (o *BackupProperties) GetPointInTimeWindowHours() int32 {
	if o == nil || IsNil(o.PointInTimeWindowHours) {
		var ret int32
		return ret
	}
	return *o.PointInTimeWindowHours
}

// GetPointInTimeWindowHoursOk returns a tuple with the PointInTimeWindowHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupProperties) GetPointInTimeWindowHoursOk() (*int32, bool) {
	if o == nil || IsNil(o.PointInTimeWindowHours) {
		return nil, false
	}
	return o.PointInTimeWindowHours, true
}

// HasPointInTimeWindowHours returns a boolean if a field has been set.
func (o *BackupProperties) HasPointInTimeWindowHours() bool {
	if o != nil && !IsNil(o.PointInTimeWindowHours) {
		return true
	}

	return false
}

// SetPointInTimeWindowHours gets a reference to the given int32 and assigns it to the PointInTimeWindowHours field.
func (o *BackupProperties) SetPointInTimeWindowHours(v int32) {
	o.PointInTimeWindowHours = &v
}

// GetBackupRetention returns the BackupRetention field value if set, zero value otherwise.
func (o *BackupProperties) GetBackupRetention() BackupRetentionProperties {
	if o == nil || IsNil(o.BackupRetention) {
		var ret BackupRetentionProperties
		return ret
	}
	return *o.BackupRetention
}

// GetBackupRetentionOk returns a tuple with the BackupRetention field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupProperties) GetBackupRetentionOk() (*BackupRetentionProperties, bool) {
	if o == nil || IsNil(o.BackupRetention) {
		return nil, false
	}
	return o.BackupRetention, true
}

// HasBackupRetention returns a boolean if a field has been set.
func (o *BackupProperties) HasBackupRetention() bool {
	if o != nil && !IsNil(o.BackupRetention) {
		return true
	}

	return false
}

// SetBackupRetention gets a reference to the given BackupRetentionProperties and assigns it to the BackupRetention field.
func (o *BackupProperties) SetBackupRetention(v BackupRetentionProperties) {
	o.BackupRetention = &v
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

func (o BackupProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.SnapshotIntervalHours) {
		toSerialize["snapshotIntervalHours"] = o.SnapshotIntervalHours
	}
	if !IsNil(o.PointInTimeWindowHours) {
		toSerialize["pointInTimeWindowHours"] = o.PointInTimeWindowHours
	}
	if !IsNil(o.BackupRetention) {
		toSerialize["backupRetention"] = o.BackupRetention
	}
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
