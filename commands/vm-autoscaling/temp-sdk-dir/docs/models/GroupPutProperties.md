# GroupPutProperties

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Datacenter** | Pointer to [**GroupPutPropertiesDatacenter**](GroupPutPropertiesDatacenter.md) |  | [optional] |
|**Location** | **string** | The data center location. | [readonly] |
|**MaxReplicaCount** | **int64** | The maximum value for the number of replicas on a VM Auto Scaling Group. Must be &gt;&#x3D; 0 and &lt;&#x3D; 200. Will be enforced for both automatic and manual changes. | |
|**MinReplicaCount** | **int64** | The minimum value for the number of replicas on a VM Auto Scaling Group. Must be &gt;&#x3D; 0 and &lt;&#x3D; 200. Will be enforced for both automatic and manual changes | |
|**Name** | **string** | The name of the VM Auto Scaling Group. This field must not be null or blank. | |
|**Policy** | [**GroupPolicy**](GroupPolicy.md) |  | |
|**ReplicaConfiguration** | [**ReplicaPropertiesPost**](ReplicaPropertiesPost.md) |  | |

## Methods

### NewGroupPutProperties

`func NewGroupPutProperties(location string, maxReplicaCount int64, minReplicaCount int64, name string, policy GroupPolicy, replicaConfiguration ReplicaPropertiesPost, ) *GroupPutProperties`

NewGroupPutProperties instantiates a new GroupPutProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPutPropertiesWithDefaults

`func NewGroupPutPropertiesWithDefaults() *GroupPutProperties`

NewGroupPutPropertiesWithDefaults instantiates a new GroupPutProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDatacenter

`func (o *GroupPutProperties) GetDatacenter() GroupPutPropertiesDatacenter`

GetDatacenter returns the Datacenter field if non-nil, zero value otherwise.

### GetDatacenterOk

`func (o *GroupPutProperties) GetDatacenterOk() (*GroupPutPropertiesDatacenter, bool)`

GetDatacenterOk returns a tuple with the Datacenter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenter

`func (o *GroupPutProperties) SetDatacenter(v GroupPutPropertiesDatacenter)`

SetDatacenter sets Datacenter field to given value.

### HasDatacenter

`func (o *GroupPutProperties) HasDatacenter() bool`

HasDatacenter returns a boolean if a field has been set.

### GetLocation

`func (o *GroupPutProperties) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *GroupPutProperties) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *GroupPutProperties) SetLocation(v string)`

SetLocation sets Location field to given value.


### GetMaxReplicaCount

`func (o *GroupPutProperties) GetMaxReplicaCount() int64`

GetMaxReplicaCount returns the MaxReplicaCount field if non-nil, zero value otherwise.

### GetMaxReplicaCountOk

`func (o *GroupPutProperties) GetMaxReplicaCountOk() (*int64, bool)`

GetMaxReplicaCountOk returns a tuple with the MaxReplicaCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxReplicaCount

`func (o *GroupPutProperties) SetMaxReplicaCount(v int64)`

SetMaxReplicaCount sets MaxReplicaCount field to given value.


### GetMinReplicaCount

`func (o *GroupPutProperties) GetMinReplicaCount() int64`

GetMinReplicaCount returns the MinReplicaCount field if non-nil, zero value otherwise.

### GetMinReplicaCountOk

`func (o *GroupPutProperties) GetMinReplicaCountOk() (*int64, bool)`

GetMinReplicaCountOk returns a tuple with the MinReplicaCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinReplicaCount

`func (o *GroupPutProperties) SetMinReplicaCount(v int64)`

SetMinReplicaCount sets MinReplicaCount field to given value.


### GetName

`func (o *GroupPutProperties) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GroupPutProperties) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GroupPutProperties) SetName(v string)`

SetName sets Name field to given value.


### GetPolicy

`func (o *GroupPutProperties) GetPolicy() GroupPolicy`

GetPolicy returns the Policy field if non-nil, zero value otherwise.

### GetPolicyOk

`func (o *GroupPutProperties) GetPolicyOk() (*GroupPolicy, bool)`

GetPolicyOk returns a tuple with the Policy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicy

`func (o *GroupPutProperties) SetPolicy(v GroupPolicy)`

SetPolicy sets Policy field to given value.


### GetReplicaConfiguration

`func (o *GroupPutProperties) GetReplicaConfiguration() ReplicaPropertiesPost`

GetReplicaConfiguration returns the ReplicaConfiguration field if non-nil, zero value otherwise.

### GetReplicaConfigurationOk

`func (o *GroupPutProperties) GetReplicaConfigurationOk() (*ReplicaPropertiesPost, bool)`

GetReplicaConfigurationOk returns a tuple with the ReplicaConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicaConfiguration

`func (o *GroupPutProperties) SetReplicaConfiguration(v ReplicaPropertiesPost)`

SetReplicaConfiguration sets ReplicaConfiguration field to given value.



