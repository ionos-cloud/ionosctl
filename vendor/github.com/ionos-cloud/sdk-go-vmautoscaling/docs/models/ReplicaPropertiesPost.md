# ReplicaPropertiesPost

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**AvailabilityZone** | Pointer to [**NullableAvailabilityZone**](AvailabilityZone.md) |  | [optional] |
|**Cores** | **int32** | The total number of cores for the VMs. | |
|**CpuFamily** | Pointer to [**CpuFamily**](CpuFamily.md) |  | [optional] |
|**Nics** | Pointer to [**[]ReplicaNic**](ReplicaNic.md) | The list of NICs associated with this replica. | [optional] |
|**Ram** | **int32** | The size of the memory for the VMs in MB. The size must be in multiples of 256 MB, with a minimum of 256 MB; if you set &#39;ramHotPlug&#x3D;TRUE&#39;, you must use at least 1024 MB. If you set the RAM size to more than 240 GB, &#39;ramHotPlug&#x3D;FALSE&#39; is fixed. | |
|**Volumes** | Pointer to [**[]ReplicaVolumePost**](ReplicaVolumePost.md) | List of volumes associated with this Replica. | [optional] |

## Methods

### NewReplicaPropertiesPost

`func NewReplicaPropertiesPost(cores int32, ram int32, ) *ReplicaPropertiesPost`

NewReplicaPropertiesPost instantiates a new ReplicaPropertiesPost object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaPropertiesPostWithDefaults

`func NewReplicaPropertiesPostWithDefaults() *ReplicaPropertiesPost`

NewReplicaPropertiesPostWithDefaults instantiates a new ReplicaPropertiesPost object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailabilityZone

`func (o *ReplicaPropertiesPost) GetAvailabilityZone() AvailabilityZone`

GetAvailabilityZone returns the AvailabilityZone field if non-nil, zero value otherwise.

### GetAvailabilityZoneOk

`func (o *ReplicaPropertiesPost) GetAvailabilityZoneOk() (*AvailabilityZone, bool)`

GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityZone

`func (o *ReplicaPropertiesPost) SetAvailabilityZone(v AvailabilityZone)`

SetAvailabilityZone sets AvailabilityZone field to given value.

### HasAvailabilityZone

`func (o *ReplicaPropertiesPost) HasAvailabilityZone() bool`

HasAvailabilityZone returns a boolean if a field has been set.

### SetAvailabilityZoneNil

`func (o *ReplicaPropertiesPost) SetAvailabilityZoneNil(b bool)`

 SetAvailabilityZoneNil sets the value for AvailabilityZone to be an explicit nil

### UnsetAvailabilityZone
`func (o *ReplicaPropertiesPost) UnsetAvailabilityZone()`

UnsetAvailabilityZone ensures that no value is present for AvailabilityZone, not even an explicit nil
### GetCores

`func (o *ReplicaPropertiesPost) GetCores() int32`

GetCores returns the Cores field if non-nil, zero value otherwise.

### GetCoresOk

`func (o *ReplicaPropertiesPost) GetCoresOk() (*int32, bool)`

GetCoresOk returns a tuple with the Cores field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCores

`func (o *ReplicaPropertiesPost) SetCores(v int32)`

SetCores sets Cores field to given value.


### GetCpuFamily

`func (o *ReplicaPropertiesPost) GetCpuFamily() CpuFamily`

GetCpuFamily returns the CpuFamily field if non-nil, zero value otherwise.

### GetCpuFamilyOk

`func (o *ReplicaPropertiesPost) GetCpuFamilyOk() (*CpuFamily, bool)`

GetCpuFamilyOk returns a tuple with the CpuFamily field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpuFamily

`func (o *ReplicaPropertiesPost) SetCpuFamily(v CpuFamily)`

SetCpuFamily sets CpuFamily field to given value.

### HasCpuFamily

`func (o *ReplicaPropertiesPost) HasCpuFamily() bool`

HasCpuFamily returns a boolean if a field has been set.

### GetNics

`func (o *ReplicaPropertiesPost) GetNics() []ReplicaNic`

GetNics returns the Nics field if non-nil, zero value otherwise.

### GetNicsOk

`func (o *ReplicaPropertiesPost) GetNicsOk() (*[]ReplicaNic, bool)`

GetNicsOk returns a tuple with the Nics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNics

`func (o *ReplicaPropertiesPost) SetNics(v []ReplicaNic)`

SetNics sets Nics field to given value.

### HasNics

`func (o *ReplicaPropertiesPost) HasNics() bool`

HasNics returns a boolean if a field has been set.

### GetRam

`func (o *ReplicaPropertiesPost) GetRam() int32`

GetRam returns the Ram field if non-nil, zero value otherwise.

### GetRamOk

`func (o *ReplicaPropertiesPost) GetRamOk() (*int32, bool)`

GetRamOk returns a tuple with the Ram field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRam

`func (o *ReplicaPropertiesPost) SetRam(v int32)`

SetRam sets Ram field to given value.


### GetVolumes

`func (o *ReplicaPropertiesPost) GetVolumes() []ReplicaVolumePost`

GetVolumes returns the Volumes field if non-nil, zero value otherwise.

### GetVolumesOk

`func (o *ReplicaPropertiesPost) GetVolumesOk() (*[]ReplicaVolumePost, bool)`

GetVolumesOk returns a tuple with the Volumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumes

`func (o *ReplicaPropertiesPost) SetVolumes(v []ReplicaVolumePost)`

SetVolumes sets Volumes field to given value.

### HasVolumes

`func (o *ReplicaPropertiesPost) HasVolumes() bool`

HasVolumes returns a boolean if a field has been set.


