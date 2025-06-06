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

// checks if the CreateClusterProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateClusterProperties{}

// CreateClusterProperties The properties with all data needed to create a new MongoDB cluster.
type CreateClusterProperties struct {
	// The cluster type, either `replicaset` or `sharded-cluster`.
	Type *string `json:"type,omitempty"`
	// The unique ID of the template, which specifies the number of cores, storage size, and memory. You cannot downgrade to a smaller template or minor edition (e.g. from business to playground). To get a list of all templates to confirm the changes use the /templates endpoint.
	TemplateID *string `json:"templateID,omitempty"`
	// The MongoDB version of your cluster.
	MongoDBVersion string `json:"mongoDBVersion"`
	// The total number of instances in the cluster (one primary and n-1 secondaries).
	Instances int32 `json:"instances"`
	// The total number of shards in the cluster.
	Shards      *int32       `json:"shards,omitempty"`
	Connections []Connection `json:"connections"`
	// The physical location where the cluster will be created. This is the location where all your instances will be located. This property is immutable.
	Location string            `json:"location"`
	Backup   *BackupProperties `json:"backup,omitempty"`
	// The name of your cluster.
	DisplayName       string                 `json:"displayName"`
	MaintenanceWindow *MaintenanceWindow     `json:"maintenanceWindow,omitempty"`
	BiConnector       *BiConnectorProperties `json:"biConnector,omitempty"`
	FromBackup        *CreateRestoreRequest  `json:"fromBackup,omitempty"`
	// The cluster edition.
	Edition *string `json:"edition,omitempty"`
	// The number of CPU cores per instance.
	Cores *int32 `json:"cores,omitempty"`
	// The amount of memory per instance in megabytes. Has to be a multiple of 1024.
	Ram *int32 `json:"ram,omitempty"`
	// The amount of storage per instance in megabytes.
	StorageSize *int32       `json:"storageSize,omitempty"`
	StorageType *StorageType `json:"storageType,omitempty"`
}

// NewCreateClusterProperties instantiates a new CreateClusterProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateClusterProperties(mongoDBVersion string, instances int32, connections []Connection, location string, displayName string) *CreateClusterProperties {
	this := CreateClusterProperties{}

	this.MongoDBVersion = mongoDBVersion
	this.Instances = instances
	this.Connections = connections
	this.Location = location
	this.DisplayName = displayName

	return &this
}

// NewCreateClusterPropertiesWithDefaults instantiates a new CreateClusterProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateClusterPropertiesWithDefaults() *CreateClusterProperties {
	this := CreateClusterProperties{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *CreateClusterProperties) SetType(v string) {
	o.Type = &v
}

// GetTemplateID returns the TemplateID field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetTemplateID() string {
	if o == nil || IsNil(o.TemplateID) {
		var ret string
		return ret
	}
	return *o.TemplateID
}

// GetTemplateIDOk returns a tuple with the TemplateID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetTemplateIDOk() (*string, bool) {
	if o == nil || IsNil(o.TemplateID) {
		return nil, false
	}
	return o.TemplateID, true
}

// HasTemplateID returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasTemplateID() bool {
	if o != nil && !IsNil(o.TemplateID) {
		return true
	}

	return false
}

// SetTemplateID gets a reference to the given string and assigns it to the TemplateID field.
func (o *CreateClusterProperties) SetTemplateID(v string) {
	o.TemplateID = &v
}

// GetMongoDBVersion returns the MongoDBVersion field value
func (o *CreateClusterProperties) GetMongoDBVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MongoDBVersion
}

// GetMongoDBVersionOk returns a tuple with the MongoDBVersion field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetMongoDBVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MongoDBVersion, true
}

// SetMongoDBVersion sets field value
func (o *CreateClusterProperties) SetMongoDBVersion(v string) {
	o.MongoDBVersion = v
}

// GetInstances returns the Instances field value
func (o *CreateClusterProperties) GetInstances() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Instances
}

// GetInstancesOk returns a tuple with the Instances field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetInstancesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Instances, true
}

// SetInstances sets field value
func (o *CreateClusterProperties) SetInstances(v int32) {
	o.Instances = v
}

// GetShards returns the Shards field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetShards() int32 {
	if o == nil || IsNil(o.Shards) {
		var ret int32
		return ret
	}
	return *o.Shards
}

// GetShardsOk returns a tuple with the Shards field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetShardsOk() (*int32, bool) {
	if o == nil || IsNil(o.Shards) {
		return nil, false
	}
	return o.Shards, true
}

// HasShards returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasShards() bool {
	if o != nil && !IsNil(o.Shards) {
		return true
	}

	return false
}

// SetShards gets a reference to the given int32 and assigns it to the Shards field.
func (o *CreateClusterProperties) SetShards(v int32) {
	o.Shards = &v
}

// GetConnections returns the Connections field value
func (o *CreateClusterProperties) GetConnections() []Connection {
	if o == nil {
		var ret []Connection
		return ret
	}

	return o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetConnectionsOk() ([]Connection, bool) {
	if o == nil {
		return nil, false
	}
	return o.Connections, true
}

// SetConnections sets field value
func (o *CreateClusterProperties) SetConnections(v []Connection) {
	o.Connections = v
}

// GetLocation returns the Location field value
func (o *CreateClusterProperties) GetLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *CreateClusterProperties) SetLocation(v string) {
	o.Location = v
}

// GetBackup returns the Backup field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetBackup() BackupProperties {
	if o == nil || IsNil(o.Backup) {
		var ret BackupProperties
		return ret
	}
	return *o.Backup
}

// GetBackupOk returns a tuple with the Backup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetBackupOk() (*BackupProperties, bool) {
	if o == nil || IsNil(o.Backup) {
		return nil, false
	}
	return o.Backup, true
}

// HasBackup returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasBackup() bool {
	if o != nil && !IsNil(o.Backup) {
		return true
	}

	return false
}

// SetBackup gets a reference to the given BackupProperties and assigns it to the Backup field.
func (o *CreateClusterProperties) SetBackup(v BackupProperties) {
	o.Backup = &v
}

// GetDisplayName returns the DisplayName field value
func (o *CreateClusterProperties) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *CreateClusterProperties) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetMaintenanceWindow returns the MaintenanceWindow field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetMaintenanceWindow() MaintenanceWindow {
	if o == nil || IsNil(o.MaintenanceWindow) {
		var ret MaintenanceWindow
		return ret
	}
	return *o.MaintenanceWindow
}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil || IsNil(o.MaintenanceWindow) {
		return nil, false
	}
	return o.MaintenanceWindow, true
}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasMaintenanceWindow() bool {
	if o != nil && !IsNil(o.MaintenanceWindow) {
		return true
	}

	return false
}

// SetMaintenanceWindow gets a reference to the given MaintenanceWindow and assigns it to the MaintenanceWindow field.
func (o *CreateClusterProperties) SetMaintenanceWindow(v MaintenanceWindow) {
	o.MaintenanceWindow = &v
}

// GetBiConnector returns the BiConnector field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetBiConnector() BiConnectorProperties {
	if o == nil || IsNil(o.BiConnector) {
		var ret BiConnectorProperties
		return ret
	}
	return *o.BiConnector
}

// GetBiConnectorOk returns a tuple with the BiConnector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetBiConnectorOk() (*BiConnectorProperties, bool) {
	if o == nil || IsNil(o.BiConnector) {
		return nil, false
	}
	return o.BiConnector, true
}

// HasBiConnector returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasBiConnector() bool {
	if o != nil && !IsNil(o.BiConnector) {
		return true
	}

	return false
}

// SetBiConnector gets a reference to the given BiConnectorProperties and assigns it to the BiConnector field.
func (o *CreateClusterProperties) SetBiConnector(v BiConnectorProperties) {
	o.BiConnector = &v
}

// GetFromBackup returns the FromBackup field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetFromBackup() CreateRestoreRequest {
	if o == nil || IsNil(o.FromBackup) {
		var ret CreateRestoreRequest
		return ret
	}
	return *o.FromBackup
}

// GetFromBackupOk returns a tuple with the FromBackup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetFromBackupOk() (*CreateRestoreRequest, bool) {
	if o == nil || IsNil(o.FromBackup) {
		return nil, false
	}
	return o.FromBackup, true
}

// HasFromBackup returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasFromBackup() bool {
	if o != nil && !IsNil(o.FromBackup) {
		return true
	}

	return false
}

// SetFromBackup gets a reference to the given CreateRestoreRequest and assigns it to the FromBackup field.
func (o *CreateClusterProperties) SetFromBackup(v CreateRestoreRequest) {
	o.FromBackup = &v
}

// GetEdition returns the Edition field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetEdition() string {
	if o == nil || IsNil(o.Edition) {
		var ret string
		return ret
	}
	return *o.Edition
}

// GetEditionOk returns a tuple with the Edition field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetEditionOk() (*string, bool) {
	if o == nil || IsNil(o.Edition) {
		return nil, false
	}
	return o.Edition, true
}

// HasEdition returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasEdition() bool {
	if o != nil && !IsNil(o.Edition) {
		return true
	}

	return false
}

// SetEdition gets a reference to the given string and assigns it to the Edition field.
func (o *CreateClusterProperties) SetEdition(v string) {
	o.Edition = &v
}

// GetCores returns the Cores field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetCores() int32 {
	if o == nil || IsNil(o.Cores) {
		var ret int32
		return ret
	}
	return *o.Cores
}

// GetCoresOk returns a tuple with the Cores field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetCoresOk() (*int32, bool) {
	if o == nil || IsNil(o.Cores) {
		return nil, false
	}
	return o.Cores, true
}

// HasCores returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasCores() bool {
	if o != nil && !IsNil(o.Cores) {
		return true
	}

	return false
}

// SetCores gets a reference to the given int32 and assigns it to the Cores field.
func (o *CreateClusterProperties) SetCores(v int32) {
	o.Cores = &v
}

// GetRam returns the Ram field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetRam() int32 {
	if o == nil || IsNil(o.Ram) {
		var ret int32
		return ret
	}
	return *o.Ram
}

// GetRamOk returns a tuple with the Ram field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetRamOk() (*int32, bool) {
	if o == nil || IsNil(o.Ram) {
		return nil, false
	}
	return o.Ram, true
}

// HasRam returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasRam() bool {
	if o != nil && !IsNil(o.Ram) {
		return true
	}

	return false
}

// SetRam gets a reference to the given int32 and assigns it to the Ram field.
func (o *CreateClusterProperties) SetRam(v int32) {
	o.Ram = &v
}

// GetStorageSize returns the StorageSize field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetStorageSize() int32 {
	if o == nil || IsNil(o.StorageSize) {
		var ret int32
		return ret
	}
	return *o.StorageSize
}

// GetStorageSizeOk returns a tuple with the StorageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetStorageSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.StorageSize) {
		return nil, false
	}
	return o.StorageSize, true
}

// HasStorageSize returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasStorageSize() bool {
	if o != nil && !IsNil(o.StorageSize) {
		return true
	}

	return false
}

// SetStorageSize gets a reference to the given int32 and assigns it to the StorageSize field.
func (o *CreateClusterProperties) SetStorageSize(v int32) {
	o.StorageSize = &v
}

// GetStorageType returns the StorageType field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetStorageType() StorageType {
	if o == nil || IsNil(o.StorageType) {
		var ret StorageType
		return ret
	}
	return *o.StorageType
}

// GetStorageTypeOk returns a tuple with the StorageType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil || IsNil(o.StorageType) {
		return nil, false
	}
	return o.StorageType, true
}

// HasStorageType returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasStorageType() bool {
	if o != nil && !IsNil(o.StorageType) {
		return true
	}

	return false
}

// SetStorageType gets a reference to the given StorageType and assigns it to the StorageType field.
func (o *CreateClusterProperties) SetStorageType(v StorageType) {
	o.StorageType = &v
}

func (o CreateClusterProperties) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateClusterProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.TemplateID) {
		toSerialize["templateID"] = o.TemplateID
	}
	toSerialize["mongoDBVersion"] = o.MongoDBVersion
	toSerialize["instances"] = o.Instances
	if !IsNil(o.Shards) {
		toSerialize["shards"] = o.Shards
	}
	toSerialize["connections"] = o.Connections
	toSerialize["location"] = o.Location
	if !IsNil(o.Backup) {
		toSerialize["backup"] = o.Backup
	}
	toSerialize["displayName"] = o.DisplayName
	if !IsNil(o.MaintenanceWindow) {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}
	if !IsNil(o.BiConnector) {
		toSerialize["biConnector"] = o.BiConnector
	}
	if !IsNil(o.FromBackup) {
		toSerialize["fromBackup"] = o.FromBackup
	}
	if !IsNil(o.Edition) {
		toSerialize["edition"] = o.Edition
	}
	if !IsNil(o.Cores) {
		toSerialize["cores"] = o.Cores
	}
	if !IsNil(o.Ram) {
		toSerialize["ram"] = o.Ram
	}
	if !IsNil(o.StorageSize) {
		toSerialize["storageSize"] = o.StorageSize
	}
	if !IsNil(o.StorageType) {
		toSerialize["storageType"] = o.StorageType
	}
	return toSerialize, nil
}

type NullableCreateClusterProperties struct {
	value *CreateClusterProperties
	isSet bool
}

func (v NullableCreateClusterProperties) Get() *CreateClusterProperties {
	return v.value
}

func (v *NullableCreateClusterProperties) Set(val *CreateClusterProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateClusterProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateClusterProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateClusterProperties(val *CreateClusterProperties) *NullableCreateClusterProperties {
	return &NullableCreateClusterProperties{value: val, isSet: true}
}

func (v NullableCreateClusterProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateClusterProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
