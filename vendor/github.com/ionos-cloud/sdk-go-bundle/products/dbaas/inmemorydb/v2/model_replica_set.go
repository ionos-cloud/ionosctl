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

// checks if the ReplicaSet type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReplicaSet{}

// ReplicaSet Properties with all data needed to create a new In-Memory DB replication.
type ReplicaSet struct {
	// The human readable name of your replica set.
	DisplayName string `json:"displayName"`
	// The In-Memory DB version of your replica set.
	Version string `json:"version"`
	// The total number of replicas in the replica set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is >1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance.
	Replicas        int32           `json:"replicas"`
	Resources       Resources       `json:"resources"`
	PersistenceMode PersistenceMode `json:"persistenceMode"`
	EvictionPolicy  EvictionPolicy  `json:"evictionPolicy"`
	// The network connection for your replica set. Only one connection is allowed.
	Connections       []Connection       `json:"connections"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	Backup            *BackupProperties  `json:"backup,omitempty"`
	Credentials       User               `json:"credentials"`
	// The ID of a snapshot to restore the replica set from. If set, the replica set will be created from the snapshot.
	InitialSnapshotId *string `json:"initialSnapshotId,omitempty"`
}

// NewReplicaSet instantiates a new ReplicaSet object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicaSet(displayName string, version string, replicas int32, resources Resources, persistenceMode PersistenceMode, evictionPolicy EvictionPolicy, connections []Connection, credentials User) *ReplicaSet {
	this := ReplicaSet{}

	this.DisplayName = displayName
	this.Version = version
	this.Replicas = replicas
	this.Resources = resources
	this.PersistenceMode = persistenceMode
	this.EvictionPolicy = evictionPolicy
	this.Connections = connections
	this.Credentials = credentials

	return &this
}

// NewReplicaSetWithDefaults instantiates a new ReplicaSet object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicaSetWithDefaults() *ReplicaSet {
	this := ReplicaSet{}
	var persistenceMode PersistenceMode = PERSISTENCEMODE_NONE
	this.PersistenceMode = persistenceMode
	var evictionPolicy EvictionPolicy = EVICTIONPOLICY_ALLKEYS_LRU
	this.EvictionPolicy = evictionPolicy
	return &this
}

// GetDisplayName returns the DisplayName field value
func (o *ReplicaSet) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *ReplicaSet) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetVersion returns the Version field value
func (o *ReplicaSet) GetVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Version
}

// GetVersionOk returns a tuple with the Version field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Version, true
}

// SetVersion sets field value
func (o *ReplicaSet) SetVersion(v string) {
	o.Version = v
}

// GetReplicas returns the Replicas field value
func (o *ReplicaSet) GetReplicas() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Replicas
}

// GetReplicasOk returns a tuple with the Replicas field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetReplicasOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Replicas, true
}

// SetReplicas sets field value
func (o *ReplicaSet) SetReplicas(v int32) {
	o.Replicas = v
}

// GetResources returns the Resources field value
func (o *ReplicaSet) GetResources() Resources {
	if o == nil {
		var ret Resources
		return ret
	}

	return o.Resources
}

// GetResourcesOk returns a tuple with the Resources field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetResourcesOk() (*Resources, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Resources, true
}

// SetResources sets field value
func (o *ReplicaSet) SetResources(v Resources) {
	o.Resources = v
}

// GetPersistenceMode returns the PersistenceMode field value
func (o *ReplicaSet) GetPersistenceMode() PersistenceMode {
	if o == nil {
		var ret PersistenceMode
		return ret
	}

	return o.PersistenceMode
}

// GetPersistenceModeOk returns a tuple with the PersistenceMode field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetPersistenceModeOk() (*PersistenceMode, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PersistenceMode, true
}

// SetPersistenceMode sets field value
func (o *ReplicaSet) SetPersistenceMode(v PersistenceMode) {
	o.PersistenceMode = v
}

// GetEvictionPolicy returns the EvictionPolicy field value
func (o *ReplicaSet) GetEvictionPolicy() EvictionPolicy {
	if o == nil {
		var ret EvictionPolicy
		return ret
	}

	return o.EvictionPolicy
}

// GetEvictionPolicyOk returns a tuple with the EvictionPolicy field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetEvictionPolicyOk() (*EvictionPolicy, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EvictionPolicy, true
}

// SetEvictionPolicy sets field value
func (o *ReplicaSet) SetEvictionPolicy(v EvictionPolicy) {
	o.EvictionPolicy = v
}

// GetConnections returns the Connections field value
func (o *ReplicaSet) GetConnections() []Connection {
	if o == nil {
		var ret []Connection
		return ret
	}

	return o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetConnectionsOk() ([]Connection, bool) {
	if o == nil {
		return nil, false
	}
	return o.Connections, true
}

// SetConnections sets field value
func (o *ReplicaSet) SetConnections(v []Connection) {
	o.Connections = v
}

// GetMaintenanceWindow returns the MaintenanceWindow field value if set, zero value otherwise.
func (o *ReplicaSet) GetMaintenanceWindow() MaintenanceWindow {
	if o == nil || IsNil(o.MaintenanceWindow) {
		var ret MaintenanceWindow
		return ret
	}
	return *o.MaintenanceWindow
}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil || IsNil(o.MaintenanceWindow) {
		return nil, false
	}
	return o.MaintenanceWindow, true
}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *ReplicaSet) HasMaintenanceWindow() bool {
	if o != nil && !IsNil(o.MaintenanceWindow) {
		return true
	}

	return false
}

// SetMaintenanceWindow gets a reference to the given MaintenanceWindow and assigns it to the MaintenanceWindow field.
func (o *ReplicaSet) SetMaintenanceWindow(v MaintenanceWindow) {
	o.MaintenanceWindow = &v
}

// GetBackup returns the Backup field value if set, zero value otherwise.
func (o *ReplicaSet) GetBackup() BackupProperties {
	if o == nil || IsNil(o.Backup) {
		var ret BackupProperties
		return ret
	}
	return *o.Backup
}

// GetBackupOk returns a tuple with the Backup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetBackupOk() (*BackupProperties, bool) {
	if o == nil || IsNil(o.Backup) {
		return nil, false
	}
	return o.Backup, true
}

// HasBackup returns a boolean if a field has been set.
func (o *ReplicaSet) HasBackup() bool {
	if o != nil && !IsNil(o.Backup) {
		return true
	}

	return false
}

// SetBackup gets a reference to the given BackupProperties and assigns it to the Backup field.
func (o *ReplicaSet) SetBackup(v BackupProperties) {
	o.Backup = &v
}

// GetCredentials returns the Credentials field value
func (o *ReplicaSet) GetCredentials() User {
	if o == nil {
		var ret User
		return ret
	}

	return o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetCredentialsOk() (*User, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Credentials, true
}

// SetCredentials sets field value
func (o *ReplicaSet) SetCredentials(v User) {
	o.Credentials = v
}

// GetInitialSnapshotId returns the InitialSnapshotId field value if set, zero value otherwise.
func (o *ReplicaSet) GetInitialSnapshotId() string {
	if o == nil || IsNil(o.InitialSnapshotId) {
		var ret string
		return ret
	}
	return *o.InitialSnapshotId
}

// GetInitialSnapshotIdOk returns a tuple with the InitialSnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicaSet) GetInitialSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.InitialSnapshotId) {
		return nil, false
	}
	return o.InitialSnapshotId, true
}

// HasInitialSnapshotId returns a boolean if a field has been set.
func (o *ReplicaSet) HasInitialSnapshotId() bool {
	if o != nil && !IsNil(o.InitialSnapshotId) {
		return true
	}

	return false
}

// SetInitialSnapshotId gets a reference to the given string and assigns it to the InitialSnapshotId field.
func (o *ReplicaSet) SetInitialSnapshotId(v string) {
	o.InitialSnapshotId = &v
}

func (o ReplicaSet) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ReplicaSet) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["displayName"] = o.DisplayName
	toSerialize["version"] = o.Version
	toSerialize["replicas"] = o.Replicas
	toSerialize["resources"] = o.Resources
	toSerialize["persistenceMode"] = o.PersistenceMode
	toSerialize["evictionPolicy"] = o.EvictionPolicy
	toSerialize["connections"] = o.Connections
	if !IsNil(o.MaintenanceWindow) {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}
	if !IsNil(o.Backup) {
		toSerialize["backup"] = o.Backup
	}
	toSerialize["credentials"] = o.Credentials
	if !IsNil(o.InitialSnapshotId) {
		toSerialize["initialSnapshotId"] = o.InitialSnapshotId
	}
	return toSerialize, nil
}

type NullableReplicaSet struct {
	value *ReplicaSet
	isSet bool
}

func (v NullableReplicaSet) Get() *ReplicaSet {
	return v.value
}

func (v *NullableReplicaSet) Set(val *ReplicaSet) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicaSet) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicaSet) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicaSet(val *ReplicaSet) *NullableReplicaSet {
	return &NullableReplicaSet{value: val, isSet: true}
}

func (v NullableReplicaSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicaSet) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
