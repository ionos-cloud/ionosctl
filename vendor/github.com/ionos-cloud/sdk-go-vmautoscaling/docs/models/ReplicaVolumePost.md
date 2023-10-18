# ReplicaVolumePost

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
|**ImagePassword** | Pointer to **string** | The image password for this replica volume. | [optional] |

## Methods

### NewReplicaVolumePost

`func NewReplicaVolumePost(name string, size int32, type_ VolumeHwType, bootOrder string, ) *ReplicaVolumePost`

NewReplicaVolumePost instantiates a new ReplicaVolumePost object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaVolumePostWithDefaults

`func NewReplicaVolumePostWithDefaults() *ReplicaVolumePost`

NewReplicaVolumePostWithDefaults instantiates a new ReplicaVolumePost object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetImage

`func (o *ReplicaVolumePost) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ReplicaVolumePost) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ReplicaVolumePost) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ReplicaVolumePost) HasImage() bool`

HasImage returns a boolean if a field has been set.

### SetImageNil

`func (o *ReplicaVolumePost) SetImageNil(b bool)`

 SetImageNil sets the value for Image to be an explicit nil

### UnsetImage
`func (o *ReplicaVolumePost) UnsetImage()`

UnsetImage ensures that no value is present for Image, not even an explicit nil
### GetImageAlias

`func (o *ReplicaVolumePost) GetImageAlias() string`

GetImageAlias returns the ImageAlias field if non-nil, zero value otherwise.

### GetImageAliasOk

`func (o *ReplicaVolumePost) GetImageAliasOk() (*string, bool)`

GetImageAliasOk returns a tuple with the ImageAlias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageAlias

`func (o *ReplicaVolumePost) SetImageAlias(v string)`

SetImageAlias sets ImageAlias field to given value.

### HasImageAlias

`func (o *ReplicaVolumePost) HasImageAlias() bool`

HasImageAlias returns a boolean if a field has been set.

### SetImageAliasNil

`func (o *ReplicaVolumePost) SetImageAliasNil(b bool)`

 SetImageAliasNil sets the value for ImageAlias to be an explicit nil

### UnsetImageAlias
`func (o *ReplicaVolumePost) UnsetImageAlias()`

UnsetImageAlias ensures that no value is present for ImageAlias, not even an explicit nil
### GetName

`func (o *ReplicaVolumePost) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ReplicaVolumePost) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ReplicaVolumePost) SetName(v string)`

SetName sets Name field to given value.


### GetSize

`func (o *ReplicaVolumePost) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *ReplicaVolumePost) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *ReplicaVolumePost) SetSize(v int32)`

SetSize sets Size field to given value.


### GetSshKeys

`func (o *ReplicaVolumePost) GetSshKeys() []string`

GetSshKeys returns the SshKeys field if non-nil, zero value otherwise.

### GetSshKeysOk

`func (o *ReplicaVolumePost) GetSshKeysOk() (*[]string, bool)`

GetSshKeysOk returns a tuple with the SshKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeys

`func (o *ReplicaVolumePost) SetSshKeys(v []string)`

SetSshKeys sets SshKeys field to given value.

### HasSshKeys

`func (o *ReplicaVolumePost) HasSshKeys() bool`

HasSshKeys returns a boolean if a field has been set.

### GetType

`func (o *ReplicaVolumePost) GetType() VolumeHwType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ReplicaVolumePost) GetTypeOk() (*VolumeHwType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ReplicaVolumePost) SetType(v VolumeHwType)`

SetType sets Type field to given value.


### GetUserData

`func (o *ReplicaVolumePost) GetUserData() string`

GetUserData returns the UserData field if non-nil, zero value otherwise.

### GetUserDataOk

`func (o *ReplicaVolumePost) GetUserDataOk() (*string, bool)`

GetUserDataOk returns a tuple with the UserData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserData

`func (o *ReplicaVolumePost) SetUserData(v string)`

SetUserData sets UserData field to given value.

### HasUserData

`func (o *ReplicaVolumePost) HasUserData() bool`

HasUserData returns a boolean if a field has been set.

### GetBus

`func (o *ReplicaVolumePost) GetBus() BusType`

GetBus returns the Bus field if non-nil, zero value otherwise.

### GetBusOk

`func (o *ReplicaVolumePost) GetBusOk() (*BusType, bool)`

GetBusOk returns a tuple with the Bus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBus

`func (o *ReplicaVolumePost) SetBus(v BusType)`

SetBus sets Bus field to given value.

### HasBus

`func (o *ReplicaVolumePost) HasBus() bool`

HasBus returns a boolean if a field has been set.

### GetBackupunitId

`func (o *ReplicaVolumePost) GetBackupunitId() string`

GetBackupunitId returns the BackupunitId field if non-nil, zero value otherwise.

### GetBackupunitIdOk

`func (o *ReplicaVolumePost) GetBackupunitIdOk() (*string, bool)`

GetBackupunitIdOk returns a tuple with the BackupunitId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackupunitId

`func (o *ReplicaVolumePost) SetBackupunitId(v string)`

SetBackupunitId sets BackupunitId field to given value.

### HasBackupunitId

`func (o *ReplicaVolumePost) HasBackupunitId() bool`

HasBackupunitId returns a boolean if a field has been set.

### GetBootOrder

`func (o *ReplicaVolumePost) GetBootOrder() string`

GetBootOrder returns the BootOrder field if non-nil, zero value otherwise.

### GetBootOrderOk

`func (o *ReplicaVolumePost) GetBootOrderOk() (*string, bool)`

GetBootOrderOk returns a tuple with the BootOrder field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBootOrder

`func (o *ReplicaVolumePost) SetBootOrder(v string)`

SetBootOrder sets BootOrder field to given value.


### GetImagePassword

`func (o *ReplicaVolumePost) GetImagePassword() string`

GetImagePassword returns the ImagePassword field if non-nil, zero value otherwise.

### GetImagePasswordOk

`func (o *ReplicaVolumePost) GetImagePasswordOk() (*string, bool)`

GetImagePasswordOk returns a tuple with the ImagePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImagePassword

`func (o *ReplicaVolumePost) SetImagePassword(v string)`

SetImagePassword sets ImagePassword field to given value.

### HasImagePassword

`func (o *ReplicaVolumePost) HasImagePassword() bool`

HasImagePassword returns a boolean if a field has been set.


