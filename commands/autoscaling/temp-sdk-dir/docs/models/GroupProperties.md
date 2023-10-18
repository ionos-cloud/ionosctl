# GroupProperties

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Datacenter** | Pointer to [**GroupPropertiesDatacenter**](GroupPropertiesDatacenter.md) |  | [optional] |
|**Location** | **string** | The data center location. | [readonly] |
|**MaxReplicaCount** | Pointer to **int64** | The maximum value for the number of replicas. Must be &gt;&#x3D; 0 and &lt;&#x3D; 100. Will be enforced for both automatic and manual changes. | [optional] |
|**MinReplicaCount** | Pointer to **int64** | The minimum value for the number of replicas. Must be &gt;&#x3D; 0 and &lt;&#x3D; 100. Will be enforced for both automatic and manual changes | [optional] |
|**Name** | Pointer to **string** | The name of the VM Auto Scaling Group. This field must not be null or blank. | [optional] |
|**Policy** | Pointer to [**GroupPolicy**](GroupPolicy.md) |  | [optional] |
|**ReplicaConfiguration** | Pointer to [**ReplicaPropertiesPost**](ReplicaPropertiesPost.md) |  | [optional] |

## Methods

### NewGroupProperties

`func NewGroupProperties(location string, ) *GroupProperties`

NewGroupProperties instantiates a new GroupProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPropertiesWithDefaults

`func NewGroupPropertiesWithDefaults() *GroupProperties`

NewGroupPropertiesWithDefaults instantiates a new GroupProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDatacenter

`func (o *GroupProperties) GetDatacenter() GroupPropertiesDatacenter`

GetDatacenter returns the Datacenter field if non-nil, zero value otherwise.

### GetDatacenterOk

`func (o *GroupProperties) GetDatacenterOk() (*GroupPropertiesDatacenter, bool)`

GetDatacenterOk returns a tuple with the Datacenter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenter

`func (o *GroupProperties) SetDatacenter(v GroupPropertiesDatacenter)`

SetDatacenter sets Datacenter field to given value.

### HasDatacenter

`func (o *GroupProperties) HasDatacenter() bool`

HasDatacenter returns a boolean if a field has been set.

### GetLocation

`func (o *GroupProperties) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *GroupProperties) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *GroupProperties) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMaxReplicaCount

`func (o *GroupProperties) GetMaxReplicaCount() int64`

GetMaxReplicaCount returns the MaxReplicaCount field if non-nil, zero value otherwise.

### GetMaxReplicaCountOk

`func (o *GroupProperties) GetMaxReplicaCountOk() (*int64, bool)`

GetMaxReplicaCountOk returns a tuple with the MaxReplicaCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxReplicaCount

`func (o *GroupProperties) SetMaxReplicaCount(v int64)`

SetMaxReplicaCount sets MaxReplicaCount field to given value.

### HasMaxReplicaCount

`func (o *GroupProperties) HasMaxReplicaCount() bool`

HasMaxReplicaCount returns a boolean if a field has been set.

### GetMinReplicaCount

`func (o *GroupProperties) GetMinReplicaCount() int64`

GetMinReplicaCount returns the MinReplicaCount field if non-nil, zero value otherwise.

### GetMinReplicaCountOk

`func (o *GroupProperties) GetMinReplicaCountOk() (*int64, bool)`

GetMinReplicaCountOk returns a tuple with the MinReplicaCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinReplicaCount

`func (o *GroupProperties) SetMinReplicaCount(v int64)`

SetMinReplicaCount sets MinReplicaCount field to given value.

### HasMinReplicaCount

`func (o *GroupProperties) HasMinReplicaCount() bool`

HasMinReplicaCount returns a boolean if a field has been set.

### GetName

`func (o *GroupProperties) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GroupProperties) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GroupProperties) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *GroupProperties) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPolicy

`func (o *GroupProperties) GetPolicy() GroupPolicy`

GetPolicy returns the Policy field if non-nil, zero value otherwise.

### GetPolicyOk

`func (o *GroupProperties) GetPolicyOk() (*GroupPolicy, bool)`

GetPolicyOk returns a tuple with the Policy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicy

`func (o *GroupProperties) SetPolicy(v GroupPolicy)`

SetPolicy sets Policy field to given value.

### HasPolicy

`func (o *GroupProperties) HasPolicy() bool`

HasPolicy returns a boolean if a field has been set.

### GetReplicaConfiguration

`func (o *GroupProperties) GetReplicaConfiguration() ReplicaPropertiesPost`

GetReplicaConfiguration returns the ReplicaConfiguration field if non-nil, zero value otherwise.

### GetReplicaConfigurationOk

`func (o *GroupProperties) GetReplicaConfigurationOk() (*ReplicaPropertiesPost, bool)`

GetReplicaConfigurationOk returns a tuple with the ReplicaConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicaConfiguration

`func (o *GroupProperties) SetReplicaConfiguration(v ReplicaPropertiesPost)`

SetReplicaConfiguration sets ReplicaConfiguration field to given value.

### HasReplicaConfiguration

`func (o *GroupProperties) HasReplicaConfiguration() bool`

HasReplicaConfiguration returns a boolean if a field has been set.


