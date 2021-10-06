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
)

// Cluster A database cluster
type Cluster struct {
	// The unique ID of the resource.
	Id *string `json:"id,omitempty"`
	// Deprecated: backup is always enabled. Enables automatic backups of your cluster. 
	BackupEnabled *bool `json:"backup_enabled,omitempty"`
	// The current status reported back by the cluster.
	LifecycleStatus *string `json:"lifecycle_status,omitempty"`
	// The friendly name of your cluster.
	DisplayName *string `json:"display_name,omitempty"`
	// The PostgreSQL version of your cluster.
	PostgresVersion *string `json:"postgres_version,omitempty"`
	// The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)
	Location *string `json:"location,omitempty"`
	// The number of replicas in your cluster.
	Replicas *float32 `json:"replicas,omitempty"`
	// The amount of memory per replica.
	RamSize *string `json:"ram_size,omitempty"`
	// The number of CPU cores per replica.
	CpuCoreCount *float32 `json:"cpu_core_count,omitempty"`
	// The amount of storage per replica.
	StorageSize *string `json:"storage_size,omitempty"`
	StorageType *StorageType `json:"storage_type,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	VdcConnections *[]VDCConnection `json:"vdc_connections,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *Cluster) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *Cluster) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetBackupEnabled returns the BackupEnabled field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *Cluster) GetBackupEnabled() *bool {
	if o == nil {
		return nil
	}


	return o.BackupEnabled

}

// GetBackupEnabledOk returns a tuple with the BackupEnabled field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetBackupEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}


	return o.BackupEnabled, true
}

// SetBackupEnabled sets field value
func (o *Cluster) SetBackupEnabled(v bool) {


	o.BackupEnabled = &v

}

// HasBackupEnabled returns a boolean if a field has been set.
func (o *Cluster) HasBackupEnabled() bool {
	if o != nil && o.BackupEnabled != nil {
		return true
	}

	return false
}



// GetLifecycleStatus returns the LifecycleStatus field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetLifecycleStatus() *string {
	if o == nil {
		return nil
	}


	return o.LifecycleStatus

}

// GetLifecycleStatusOk returns a tuple with the LifecycleStatus field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetLifecycleStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LifecycleStatus, true
}

// SetLifecycleStatus sets field value
func (o *Cluster) SetLifecycleStatus(v string) {


	o.LifecycleStatus = &v

}

// HasLifecycleStatus returns a boolean if a field has been set.
func (o *Cluster) HasLifecycleStatus() bool {
	if o != nil && o.LifecycleStatus != nil {
		return true
	}

	return false
}



// GetDisplayName returns the DisplayName field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetDisplayName() *string {
	if o == nil {
		return nil
	}


	return o.DisplayName

}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.DisplayName, true
}

// SetDisplayName sets field value
func (o *Cluster) SetDisplayName(v string) {


	o.DisplayName = &v

}

// HasDisplayName returns a boolean if a field has been set.
func (o *Cluster) HasDisplayName() bool {
	if o != nil && o.DisplayName != nil {
		return true
	}

	return false
}



// GetPostgresVersion returns the PostgresVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetPostgresVersion() *string {
	if o == nil {
		return nil
	}


	return o.PostgresVersion

}

// GetPostgresVersionOk returns a tuple with the PostgresVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetPostgresVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.PostgresVersion, true
}

// SetPostgresVersion sets field value
func (o *Cluster) SetPostgresVersion(v string) {


	o.PostgresVersion = &v

}

// HasPostgresVersion returns a boolean if a field has been set.
func (o *Cluster) HasPostgresVersion() bool {
	if o != nil && o.PostgresVersion != nil {
		return true
	}

	return false
}



// GetLocation returns the Location field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetLocation() *string {
	if o == nil {
		return nil
	}


	return o.Location

}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Location, true
}

// SetLocation sets field value
func (o *Cluster) SetLocation(v string) {


	o.Location = &v

}

// HasLocation returns a boolean if a field has been set.
func (o *Cluster) HasLocation() bool {
	if o != nil && o.Location != nil {
		return true
	}

	return false
}



// GetReplicas returns the Replicas field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Cluster) GetReplicas() *float32 {
	if o == nil {
		return nil
	}


	return o.Replicas

}

// GetReplicasOk returns a tuple with the Replicas field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetReplicasOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Replicas, true
}

// SetReplicas sets field value
func (o *Cluster) SetReplicas(v float32) {


	o.Replicas = &v

}

// HasReplicas returns a boolean if a field has been set.
func (o *Cluster) HasReplicas() bool {
	if o != nil && o.Replicas != nil {
		return true
	}

	return false
}



// GetRamSize returns the RamSize field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetRamSize() *string {
	if o == nil {
		return nil
	}


	return o.RamSize

}

// GetRamSizeOk returns a tuple with the RamSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetRamSizeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.RamSize, true
}

// SetRamSize sets field value
func (o *Cluster) SetRamSize(v string) {


	o.RamSize = &v

}

// HasRamSize returns a boolean if a field has been set.
func (o *Cluster) HasRamSize() bool {
	if o != nil && o.RamSize != nil {
		return true
	}

	return false
}



// GetCpuCoreCount returns the CpuCoreCount field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Cluster) GetCpuCoreCount() *float32 {
	if o == nil {
		return nil
	}


	return o.CpuCoreCount

}

// GetCpuCoreCountOk returns a tuple with the CpuCoreCount field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetCpuCoreCountOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.CpuCoreCount, true
}

// SetCpuCoreCount sets field value
func (o *Cluster) SetCpuCoreCount(v float32) {


	o.CpuCoreCount = &v

}

// HasCpuCoreCount returns a boolean if a field has been set.
func (o *Cluster) HasCpuCoreCount() bool {
	if o != nil && o.CpuCoreCount != nil {
		return true
	}

	return false
}



// GetStorageSize returns the StorageSize field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetStorageSize() *string {
	if o == nil {
		return nil
	}


	return o.StorageSize

}

// GetStorageSizeOk returns a tuple with the StorageSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetStorageSizeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.StorageSize, true
}

// SetStorageSize sets field value
func (o *Cluster) SetStorageSize(v string) {


	o.StorageSize = &v

}

// HasStorageSize returns a boolean if a field has been set.
func (o *Cluster) HasStorageSize() bool {
	if o != nil && o.StorageSize != nil {
		return true
	}

	return false
}



// GetStorageType returns the StorageType field value
// If the value is explicit nil, the zero value for StorageType will be returned
func (o *Cluster) GetStorageType() *StorageType {
	if o == nil {
		return nil
	}


	return o.StorageType

}

// GetStorageTypeOk returns a tuple with the StorageType field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil {
		return nil, false
	}


	return o.StorageType, true
}

// SetStorageType sets field value
func (o *Cluster) SetStorageType(v StorageType) {


	o.StorageType = &v

}

// HasStorageType returns a boolean if a field has been set.
func (o *Cluster) HasStorageType() bool {
	if o != nil && o.StorageType != nil {
		return true
	}

	return false
}



// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for Metadata will be returned
func (o *Cluster) GetMetadata() *Metadata {
	if o == nil {
		return nil
	}


	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetMetadataOk() (*Metadata, bool) {
	if o == nil {
		return nil, false
	}


	return o.Metadata, true
}

// SetMetadata sets field value
func (o *Cluster) SetMetadata(v Metadata) {


	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *Cluster) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}



// GetVdcConnections returns the VdcConnections field value
// If the value is explicit nil, the zero value for []VDCConnection will be returned
func (o *Cluster) GetVdcConnections() *[]VDCConnection {
	if o == nil {
		return nil
	}


	return o.VdcConnections

}

// GetVdcConnectionsOk returns a tuple with the VdcConnections field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetVdcConnectionsOk() (*[]VDCConnection, bool) {
	if o == nil {
		return nil, false
	}


	return o.VdcConnections, true
}

// SetVdcConnections sets field value
func (o *Cluster) SetVdcConnections(v []VDCConnection) {


	o.VdcConnections = &v

}

// HasVdcConnections returns a boolean if a field has been set.
func (o *Cluster) HasVdcConnections() bool {
	if o != nil && o.VdcConnections != nil {
		return true
	}

	return false
}



// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for MaintenanceWindow will be returned
func (o *Cluster) GetMaintenanceWindow() *MaintenanceWindow {
	if o == nil {
		return nil
	}


	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}


	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *Cluster) SetMaintenanceWindow(v MaintenanceWindow) {


	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *Cluster) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}


func (o Cluster) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	

	if o.BackupEnabled != nil {
		toSerialize["backup_enabled"] = o.BackupEnabled
	}
	

	if o.LifecycleStatus != nil {
		toSerialize["lifecycle_status"] = o.LifecycleStatus
	}
	

	if o.DisplayName != nil {
		toSerialize["display_name"] = o.DisplayName
	}
	

	if o.PostgresVersion != nil {
		toSerialize["postgres_version"] = o.PostgresVersion
	}
	

	if o.Location != nil {
		toSerialize["location"] = o.Location
	}
	

	if o.Replicas != nil {
		toSerialize["replicas"] = o.Replicas
	}
	

	if o.RamSize != nil {
		toSerialize["ram_size"] = o.RamSize
	}
	

	if o.CpuCoreCount != nil {
		toSerialize["cpu_core_count"] = o.CpuCoreCount
	}
	

	if o.StorageSize != nil {
		toSerialize["storage_size"] = o.StorageSize
	}
	

	if o.StorageType != nil {
		toSerialize["storage_type"] = o.StorageType
	}
	

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	

	if o.VdcConnections != nil {
		toSerialize["vdc_connections"] = o.VdcConnections
	}
	

	if o.MaintenanceWindow != nil {
		toSerialize["maintenance_window"] = o.MaintenanceWindow
	}
	
	return json.Marshal(toSerialize)
}

type NullableCluster struct {
	value *Cluster
	isSet bool
}

func (v NullableCluster) Get() *Cluster {
	return v.value
}

func (v *NullableCluster) Set(val *Cluster) {
	v.value = val
	v.isSet = true
}

func (v NullableCluster) IsSet() bool {
	return v.isSet
}

func (v *NullableCluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCluster(val *Cluster) *NullableCluster {
	return &NullableCluster{value: val, isSet: true}
}

func (v NullableCluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


