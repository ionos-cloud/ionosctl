/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ClusterProperties Properties of a database cluster
type ClusterProperties struct {
	// The friendly name of your cluster.
	DisplayName *string `json:"displayName,omitempty"`
	// The PostgreSQL version of your cluster.
	PostgresVersion *string         `json:"postgresVersion,omitempty"`
	Location        *Location       `json:"location,omitempty"`
	BackupLocation  *BackupLocation `json:"backupLocation,omitempty"`
	// The total number of instances in the cluster (one master and n-1 standbys).
	Instances *int32 `json:"instances,omitempty"`
	// The amount of memory per instance in megabytes. Has to be a multiple of 256.
	Ram *int32 `json:"ram,omitempty"`
	// The number of CPU cores per instance.
	Cores *int32 `json:"cores,omitempty"`
	// The amount of storage per instance in megabytes.
	StorageSize         *int32               `json:"storageSize,omitempty"`
	StorageType         *StorageType         `json:"storageType,omitempty"`
	Connections         *[]Connection        `json:"connections,omitempty"`
	MaintenanceWindow   *MaintenanceWindow   `json:"maintenanceWindow,omitempty"`
	SynchronizationMode *SynchronizationMode `json:"synchronizationMode,omitempty"`
}

// GetDisplayName returns the DisplayName field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterProperties) GetDisplayName() *string {
	if o == nil {
		return nil
	}

	return o.DisplayName

}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.DisplayName, true
}

// SetDisplayName sets field value
func (o *ClusterProperties) SetDisplayName(v string) {

	o.DisplayName = &v

}

// HasDisplayName returns a boolean if a field has been set.
func (o *ClusterProperties) HasDisplayName() bool {
	if o != nil && o.DisplayName != nil {
		return true
	}

	return false
}

// GetPostgresVersion returns the PostgresVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterProperties) GetPostgresVersion() *string {
	if o == nil {
		return nil
	}

	return o.PostgresVersion

}

// GetPostgresVersionOk returns a tuple with the PostgresVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetPostgresVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.PostgresVersion, true
}

// SetPostgresVersion sets field value
func (o *ClusterProperties) SetPostgresVersion(v string) {

	o.PostgresVersion = &v

}

// HasPostgresVersion returns a boolean if a field has been set.
func (o *ClusterProperties) HasPostgresVersion() bool {
	if o != nil && o.PostgresVersion != nil {
		return true
	}

	return false
}

// GetLocation returns the Location field value
// If the value is explicit nil, the zero value for Location will be returned
func (o *ClusterProperties) GetLocation() *Location {
	if o == nil {
		return nil
	}

	return o.Location

}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetLocationOk() (*Location, bool) {
	if o == nil {
		return nil, false
	}

	return o.Location, true
}

// SetLocation sets field value
func (o *ClusterProperties) SetLocation(v Location) {

	o.Location = &v

}

// HasLocation returns a boolean if a field has been set.
func (o *ClusterProperties) HasLocation() bool {
	if o != nil && o.Location != nil {
		return true
	}

	return false
}

// GetBackupLocation returns the BackupLocation field value
// If the value is explicit nil, the zero value for BackupLocation will be returned
func (o *ClusterProperties) GetBackupLocation() *BackupLocation {
	if o == nil {
		return nil
	}

	return o.BackupLocation

}

// GetBackupLocationOk returns a tuple with the BackupLocation field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetBackupLocationOk() (*BackupLocation, bool) {
	if o == nil {
		return nil, false
	}

	return o.BackupLocation, true
}

// SetBackupLocation sets field value
func (o *ClusterProperties) SetBackupLocation(v BackupLocation) {

	o.BackupLocation = &v

}

// HasBackupLocation returns a boolean if a field has been set.
func (o *ClusterProperties) HasBackupLocation() bool {
	if o != nil && o.BackupLocation != nil {
		return true
	}

	return false
}

// GetInstances returns the Instances field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterProperties) GetInstances() *int32 {
	if o == nil {
		return nil
	}

	return o.Instances

}

// GetInstancesOk returns a tuple with the Instances field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetInstancesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Instances, true
}

// SetInstances sets field value
func (o *ClusterProperties) SetInstances(v int32) {

	o.Instances = &v

}

// HasInstances returns a boolean if a field has been set.
func (o *ClusterProperties) HasInstances() bool {
	if o != nil && o.Instances != nil {
		return true
	}

	return false
}

// GetRam returns the Ram field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterProperties) GetRam() *int32 {
	if o == nil {
		return nil
	}

	return o.Ram

}

// GetRamOk returns a tuple with the Ram field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetRamOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Ram, true
}

// SetRam sets field value
func (o *ClusterProperties) SetRam(v int32) {

	o.Ram = &v

}

// HasRam returns a boolean if a field has been set.
func (o *ClusterProperties) HasRam() bool {
	if o != nil && o.Ram != nil {
		return true
	}

	return false
}

// GetCores returns the Cores field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterProperties) GetCores() *int32 {
	if o == nil {
		return nil
	}

	return o.Cores

}

// GetCoresOk returns a tuple with the Cores field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetCoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Cores, true
}

// SetCores sets field value
func (o *ClusterProperties) SetCores(v int32) {

	o.Cores = &v

}

// HasCores returns a boolean if a field has been set.
func (o *ClusterProperties) HasCores() bool {
	if o != nil && o.Cores != nil {
		return true
	}

	return false
}

// GetStorageSize returns the StorageSize field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterProperties) GetStorageSize() *int32 {
	if o == nil {
		return nil
	}

	return o.StorageSize

}

// GetStorageSizeOk returns a tuple with the StorageSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetStorageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageSize, true
}

// SetStorageSize sets field value
func (o *ClusterProperties) SetStorageSize(v int32) {

	o.StorageSize = &v

}

// HasStorageSize returns a boolean if a field has been set.
func (o *ClusterProperties) HasStorageSize() bool {
	if o != nil && o.StorageSize != nil {
		return true
	}

	return false
}

// GetStorageType returns the StorageType field value
// If the value is explicit nil, the zero value for StorageType will be returned
func (o *ClusterProperties) GetStorageType() *StorageType {
	if o == nil {
		return nil
	}

	return o.StorageType

}

// GetStorageTypeOk returns a tuple with the StorageType field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageType, true
}

// SetStorageType sets field value
func (o *ClusterProperties) SetStorageType(v StorageType) {

	o.StorageType = &v

}

// HasStorageType returns a boolean if a field has been set.
func (o *ClusterProperties) HasStorageType() bool {
	if o != nil && o.StorageType != nil {
		return true
	}

	return false
}

// GetConnections returns the Connections field value
// If the value is explicit nil, the zero value for []Connection will be returned
func (o *ClusterProperties) GetConnections() *[]Connection {
	if o == nil {
		return nil
	}

	return o.Connections

}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetConnectionsOk() (*[]Connection, bool) {
	if o == nil {
		return nil, false
	}

	return o.Connections, true
}

// SetConnections sets field value
func (o *ClusterProperties) SetConnections(v []Connection) {

	o.Connections = &v

}

// HasConnections returns a boolean if a field has been set.
func (o *ClusterProperties) HasConnections() bool {
	if o != nil && o.Connections != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for MaintenanceWindow will be returned
func (o *ClusterProperties) GetMaintenanceWindow() *MaintenanceWindow {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *ClusterProperties) SetMaintenanceWindow(v MaintenanceWindow) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *ClusterProperties) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetSynchronizationMode returns the SynchronizationMode field value
// If the value is explicit nil, the zero value for SynchronizationMode will be returned
func (o *ClusterProperties) GetSynchronizationMode() *SynchronizationMode {
	if o == nil {
		return nil
	}

	return o.SynchronizationMode

}

// GetSynchronizationModeOk returns a tuple with the SynchronizationMode field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterProperties) GetSynchronizationModeOk() (*SynchronizationMode, bool) {
	if o == nil {
		return nil, false
	}

	return o.SynchronizationMode, true
}

// SetSynchronizationMode sets field value
func (o *ClusterProperties) SetSynchronizationMode(v SynchronizationMode) {

	o.SynchronizationMode = &v

}

// HasSynchronizationMode returns a boolean if a field has been set.
func (o *ClusterProperties) HasSynchronizationMode() bool {
	if o != nil && o.SynchronizationMode != nil {
		return true
	}

	return false
}

func (o ClusterProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.DisplayName != nil {
		toSerialize["displayName"] = o.DisplayName
	}

	if o.PostgresVersion != nil {
		toSerialize["postgresVersion"] = o.PostgresVersion
	}

	if o.Location != nil {
		toSerialize["location"] = o.Location
	}

	if o.BackupLocation != nil {
		toSerialize["backupLocation"] = o.BackupLocation
	}

	if o.Instances != nil {
		toSerialize["instances"] = o.Instances
	}

	if o.Ram != nil {
		toSerialize["ram"] = o.Ram
	}

	if o.Cores != nil {
		toSerialize["cores"] = o.Cores
	}

	if o.StorageSize != nil {
		toSerialize["storageSize"] = o.StorageSize
	}

	if o.StorageType != nil {
		toSerialize["storageType"] = o.StorageType
	}

	if o.Connections != nil {
		toSerialize["connections"] = o.Connections
	}

	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}

	if o.SynchronizationMode != nil {
		toSerialize["synchronizationMode"] = o.SynchronizationMode
	}

	return json.Marshal(toSerialize)
}

type NullableClusterProperties struct {
	value *ClusterProperties
	isSet bool
}

func (v NullableClusterProperties) Get() *ClusterProperties {
	return v.value
}

func (v *NullableClusterProperties) Set(val *ClusterProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterProperties(val *ClusterProperties) *NullableClusterProperties {
	return &NullableClusterProperties{value: val, isSet: true}
}

func (v NullableClusterProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
