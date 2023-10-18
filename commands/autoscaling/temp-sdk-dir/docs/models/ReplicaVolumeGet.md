# ReplicaVolumeGet

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Image** | Pointer to **NullableString** | The image installed on the disk. Currently, only the UUID of the image is supported.  &gt;Note that either &#39;image&#39; or &#39;imageAlias&#39; must be specified, but not both. | [optional] |
|**ImageAlias** | Pointer to **NullableString** | The image installed on the volume. Must be an &#39;imageAlias&#39; as specified via the images API. Note that one of &#39;image&#39; or &#39;imageAlias&#39; must be set, but not both. | [optional] |
|**Name** | **string** | The replica volume name. | |
|**Size** | **int32** | The size of this replica volume in GB. | |
|**SshKeys** | Pointer to **[]string** | The SSH keys of this volume. | [optional] |
|**Type** | [**VolumeHwType**](VolumeHwType.md) |  | |
|**UserData** | Pointer to **string** | The user data (Cloud Init) for this replica volume. | [optional] |
|**Bus** | Pointer to [**BusType**](BusType.md) |  | [optional] [default to BUSTYPE_VIRTIO]|
|**BackupunitId** | Pointer to **string** | The ID of the backup unit that the user has access to. The property is immutable and is only allowed to be set on creation of a new a volume. It is mandatory to provide either &#39;public image&#39; or &#39;imageAlias&#39; in conjunction with this property. | [optional] |
|**BootOrder** | **string** | Determines whether the volume will be used as a boot volume. Set to NONE, the volume will not be used as boot volume. Set to PRIMARY, the volume will be used as boot volume and set to AUTO will delegate the decision to the provisioning engine to decide whether to use the voluem as boot volume. Notice that exactly one volume can be set to PRIMARY or all of them set to AUTO. | |
|**VolumeId** | Pointer to **NullableString** | Identifies volumes to be updated. Must only be used in &#x60;GET&#x60; responses and &#x60;PUT&#x60; requests, so that updated volumes in the request will be matched with existing volumes, and passwords don&#39;t need to be re-send. | [optional] |

## Methods

### NewReplicaVolumeGet

`func NewReplicaVolumeGet(name string, size int32, type_ VolumeHwType, bootOrder string, ) *ReplicaVolumeGet`

NewReplicaVolumeGet instantiates a new ReplicaVolumeGet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaVolumeGetWithDefaults

`func NewReplicaVolumeGetWithDefaults() *ReplicaVolumeGet`

NewReplicaVolumeGetWithDefaults instantiates a new ReplicaVolumeGet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetImage

`func (o *ReplicaVolumeGet) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ReplicaVolumeGet) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ReplicaVolumeGet) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ReplicaVolumeGet) HasImage() bool`

HasImage returns a boolean if a field has been set.

### SetImageNil

`func (o *ReplicaVolumeGet) SetImageNil(b bool)`

 SetImageNil sets the value for Image to be an explicit nil

### UnsetImage
`func (o *ReplicaVolumeGet) UnsetImage()`

UnsetImage ensures that no value is present for Image, not even an explicit nil
### GetImageAlias

`func (o *ReplicaVolumeGet) GetImageAlias() string`

GetImageAlias returns the ImageAlias field if non-nil, zero value otherwise.

### GetImageAliasOk

`func (o *ReplicaVolumeGet) GetImageAliasOk() (*string, bool)`

GetImageAliasOk returns a tuple with the ImageAlias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAlias

`func (o *ReplicaVolumeGet) SetImageAlias(v string)`

SetImageAlias sets ImageAlias field to given value.

### HasImageAlias

`func (o *ReplicaVolumeGet) HasImageAlias() bool`

HasImageAlias returns a boolean if a field has been set.

### SetImageAliasNil

`func (o *ReplicaVolumeGet) SetImageAliasNil(b bool)`

 SetImageAliasNil sets the value for ImageAlias to be an explicit nil

### UnsetImageAlias
`func (o *ReplicaVolumeGet) UnsetImageAlias()`

UnsetImageAlias ensures that no value is present for ImageAlias, not even an explicit nil
### GetName

`func (o *ReplicaVolumeGet) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ReplicaVolumeGet) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ReplicaVolumeGet) SetName(v string)`

SetName sets Name field to given value.


### GetSize

`func (o *ReplicaVolumeGet) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *ReplicaVolumeGet) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *ReplicaVolumeGet) SetSize(v int32)`

SetSize sets Size field to given value.


### GetSshKeys

`func (o *ReplicaVolumeGet) GetSshKeys() []string`

GetSshKeys returns the SshKeys field if non-nil, zero value otherwise.

### GetSshKeysOk

`func (o *ReplicaVolumeGet) GetSshKeysOk() (*[]string, bool)`

GetSshKeysOk returns a tuple with the SshKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeys

`func (o *ReplicaVolumeGet) SetSshKeys(v []string)`

SetSshKeys sets SshKeys field to given value.

### HasSshKeys

`func (o *ReplicaVolumeGet) HasSshKeys() bool`

HasSshKeys returns a boolean if a field has been set.

### GetType

`func (o *ReplicaVolumeGet) GetType() VolumeHwType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ReplicaVolumeGet) GetTypeOk() (*VolumeHwType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ReplicaVolumeGet) SetType(v VolumeHwType)`

SetType sets Type field to given value.


### GetUserData

`func (o *ReplicaVolumeGet) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *ReplicaVolumeGet) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *ReplicaVolumeGet) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *ReplicaVolumeGet) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### GetBus

`func (o *ReplicaVolumeGet) GetBus() BusType`

GetBus returns the Bus field if non-nil, zero value otherwise.

### GetBusOk

`func (o *ReplicaVolumeGet) GetBusOk() (*BusType, bool)`

GetBusOk returns a tuple with the Bus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBus

`func (o *ReplicaVolumeGet) SetBus(v BusType)`

SetBus sets Bus field to given value.

### HasBus

`func (o *ReplicaVolumeGet) HasBus() bool`

HasBus returns a boolean if a field has been set.

### GetBackupunitId

`func (o *ReplicaVolumeGet) GetBackupunitId() string`

GetBackupunitId returns the BackupunitId field if non-nil, zero value otherwise.

### GetBackupunitIdOk

`func (o *ReplicaVolumeGet) GetBackupunitIdOk() (*string, bool)`

GetBackupunitIdOk returns a tuple with the BackupunitId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackupunitId

`func (o *ReplicaVolumeGet) SetBackupunitId(v string)`

SetBackupunitId sets BackupunitId field to given value.

### HasBackupunitId

`func (o *ReplicaVolumeGet) HasBackupunitId() bool`

HasBackupunitId returns a boolean if a field has been set.

### GetBootOrder

`func (o *ReplicaVolumeGet) GetBootOrder() string`

GetBootOrder returns the BootOrder field if non-nil, zero value otherwise.

### GetBootOrderOk

`func (o *ReplicaVolumeGet) GetBootOrderOk() (*string, bool)`

GetBootOrderOk returns a tuple with the BootOrder field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBootOrder

`func (o *ReplicaVolumeGet) SetBootOrder(v string)`

SetBootOrder sets BootOrder field to given value.


### GetVolumeId

`func (o *ReplicaVolumeGet) GetVolumeId() string`

GetVolumeId returns the VolumeId field if non-nil, zero value otherwise.

### GetVolumeIdOk

`func (o *ReplicaVolumeGet) GetVolumeIdOk() (*string, bool)`

GetVolumeIdOk returns a tuple with the VolumeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeId

`func (o *ReplicaVolumeGet) SetVolumeId(v string)`

SetVolumeId sets VolumeId field to given value.

### HasVolumeId

`func (o *ReplicaVolumeGet) HasVolumeId() bool`

HasVolumeId returns a boolean if a field has been set.

### SetVolumeIdNil

`func (o *ReplicaVolumeGet) SetVolumeIdNil(b bool)`

 SetVolumeIdNil sets the value for VolumeId to be an explicit nil

### UnsetVolumeId
`func (o *ReplicaVolumeGet) UnsetVolumeId()`

UnsetVolumeId ensures that no value is present for VolumeId, not even an explicit nil

