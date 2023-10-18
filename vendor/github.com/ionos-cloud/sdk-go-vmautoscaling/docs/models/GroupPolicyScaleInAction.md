# GroupPolicyScaleInAction

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Amount** | **float32** | &#39;amountType&#x3D;ABSOLUTE&#39; specifies the absolute number of VMs that are added or removed. The value must be between 1 to 10.   &#39;amountType&#x3D;PERCENTAGE&#39; specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200.   At least one VM is always added or removed.   Note that for &#39;SCALE_IN&#39; operations, volumes are not deleted after the server is deleted. | |
|**AmountType** | [**ActionAmount**](ActionAmount.md) |  | |
|**CooldownPeriod** | Pointer to **NullableString** | The minimum time that elapses after the start of this scaling action until the next scaling action is started. With a scaling action in progress, no second scaling action is started for the same VM Auto Scaling Group. Instead, the metric is re-evaluated after the current scaling action completes (either successfully or with errors). This is currently validated with a minimum value of 2 minutes and a maximum value of 24 hours. The default value is 5 minutes if not specified. | [optional] [default to "5m"]|
|**TerminationPolicy** | Pointer to [**NullableTerminationPolicyType**](TerminationPolicyType.md) |  | [optional] |
|**DeleteVolumes** | **bool** | If set to &#x60;true&#x60;, when deleting an replica during scale in, any attached volume will also be deleted. When set to &#x60;false&#x60;, all volumes remain in the datacenter and must be deleted manually.  **Note**, that every scale-out creates new volumes. When they are not deleted, they will eventually use all of your contracts resource limits. At this point, scaling out would not be possible anymore. | |

## Methods

### NewGroupPolicyScaleInAction

`func NewGroupPolicyScaleInAction(amount float32, amountType ActionAmount, deleteVolumes bool, ) *GroupPolicyScaleInAction`

NewGroupPolicyScaleInAction instantiates a new GroupPolicyScaleInAction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPolicyScaleInActionWithDefaults

`func NewGroupPolicyScaleInActionWithDefaults() *GroupPolicyScaleInAction`

NewGroupPolicyScaleInActionWithDefaults instantiates a new GroupPolicyScaleInAction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *GroupPolicyScaleInAction) GetAmount() float32`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *GroupPolicyScaleInAction) GetAmountOk() (*float32, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *GroupPolicyScaleInAction) SetAmount(v float32)`

SetAmount sets Amount field to given value.


### GetAmountType

`func (o *GroupPolicyScaleInAction) GetAmountType() ActionAmount`

GetAmountType returns the AmountType field if non-nil, zero value otherwise.

### GetAmountTypeOk

`func (o *GroupPolicyScaleInAction) GetAmountTypeOk() (*ActionAmount, bool)`

GetAmountTypeOk returns a tuple with the AmountType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmountType

`func (o *GroupPolicyScaleInAction) SetAmountType(v ActionAmount)`

SetAmountType sets AmountType field to given value.


### GetCooldownPeriod

`func (o *GroupPolicyScaleInAction) GetCooldownPeriod() string`

GetCooldownPeriod returns the CooldownPeriod field if non-nil, zero value otherwise.

### GetCooldownPeriodOk

`func (o *GroupPolicyScaleInAction) GetCooldownPeriodOk() (*string, bool)`

GetCooldownPeriodOk returns a tuple with the CooldownPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCooldownPeriod

`func (o *GroupPolicyScaleInAction) SetCooldownPeriod(v string)`

SetCooldownPeriod sets CooldownPeriod field to given value.

### HasCooldownPeriod

`func (o *GroupPolicyScaleInAction) HasCooldownPeriod() bool`

HasCooldownPeriod returns a boolean if a field has been set.

### SetCooldownPeriodNil

`func (o *GroupPolicyScaleInAction) SetCooldownPeriodNil(b bool)`

 SetCooldownPeriodNil sets the value for CooldownPeriod to be an explicit nil

### UnsetCooldownPeriod
`func (o *GroupPolicyScaleInAction) UnsetCooldownPeriod()`

UnsetCooldownPeriod ensures that no value is present for CooldownPeriod, not even an explicit nil
### GetTerminationPolicy

`func (o *GroupPolicyScaleInAction) GetTerminationPolicy() TerminationPolicyType`

GetTerminationPolicy returns the TerminationPolicy field if non-nil, zero value otherwise.

### GetTerminationPolicyOk

`func (o *GroupPolicyScaleInAction) GetTerminationPolicyOk() (*TerminationPolicyType, bool)`

GetTerminationPolicyOk returns a tuple with the TerminationPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminationPolicy

`func (o *GroupPolicyScaleInAction) SetTerminationPolicy(v TerminationPolicyType)`

SetTerminationPolicy sets TerminationPolicy field to given value.

### HasTerminationPolicy

`func (o *GroupPolicyScaleInAction) HasTerminationPolicy() bool`

HasTerminationPolicy returns a boolean if a field has been set.

### SetTerminationPolicyNil

`func (o *GroupPolicyScaleInAction) SetTerminationPolicyNil(b bool)`

 SetTerminationPolicyNil sets the value for TerminationPolicy to be an explicit nil

### UnsetTerminationPolicy
`func (o *GroupPolicyScaleInAction) UnsetTerminationPolicy()`

UnsetTerminationPolicy ensures that no value is present for TerminationPolicy, not even an explicit nil
### GetDeleteVolumes

`func (o *GroupPolicyScaleInAction) GetDeleteVolumes() bool`

GetDeleteVolumes returns the DeleteVolumes field if non-nil, zero value otherwise.

### GetDeleteVolumesOk

`func (o *GroupPolicyScaleInAction) GetDeleteVolumesOk() (*bool, bool)`

GetDeleteVolumesOk returns a tuple with the DeleteVolumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteVolumes

`func (o *GroupPolicyScaleInAction) SetDeleteVolumes(v bool)`

SetDeleteVolumes sets DeleteVolumes field to given value.



