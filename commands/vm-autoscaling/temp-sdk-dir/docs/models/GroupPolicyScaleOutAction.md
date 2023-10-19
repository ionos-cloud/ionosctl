# GroupPolicyScaleOutAction

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Amount** | **float32** | &#39;amountType&#x3D;ABSOLUTE&#39; specifies the absolute number of VMs that are added or removed. The value must be between 1 to 10.   &#39;amountType&#x3D;PERCENTAGE&#39; specifies the percentage value that is applied to the current number of replicas of the VM Auto Scaling Group. The value must be between 1 to 200.   At least one VM is always added or removed.   Note that for &#39;SCALE_IN&#39; operations, volumes are not deleted after the server is deleted. | |
|**AmountType** | [**ActionAmount**](ActionAmount.md) |  | |
|**CooldownPeriod** | Pointer to **NullableString** | The minimum time that elapses after the start of this scaling action until the following scaling action is started. While a scaling action is in progress, no second action is initiated for the same VM Auto Scaling Group. Instead, the metric is re-evaluated after the current scaling action completes (either successfully or with errors). This is currently validated with a minimum value of 2 minutes and a maximum of 24 hours. The default value is 5 minutes if not specified. | [optional] [default to "5m"]|

## Methods

### NewGroupPolicyScaleOutAction

`func NewGroupPolicyScaleOutAction(amount float32, amountType ActionAmount, ) *GroupPolicyScaleOutAction`

NewGroupPolicyScaleOutAction instantiates a new GroupPolicyScaleOutAction object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPolicyScaleOutActionWithDefaults

`func NewGroupPolicyScaleOutActionWithDefaults() *GroupPolicyScaleOutAction`

NewGroupPolicyScaleOutActionWithDefaults instantiates a new GroupPolicyScaleOutAction object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAmount

`func (o *GroupPolicyScaleOutAction) GetAmount() float32`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *GroupPolicyScaleOutAction) GetAmountOk() (*float32, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *GroupPolicyScaleOutAction) SetAmount(v float32)`

SetAmount sets Amount field to given value.


### GetAmountType

`func (o *GroupPolicyScaleOutAction) GetAmountType() ActionAmount`

GetAmountType returns the AmountType field if non-nil, zero value otherwise.

### GetAmountTypeOk

`func (o *GroupPolicyScaleOutAction) GetAmountTypeOk() (*ActionAmount, bool)`

GetAmountTypeOk returns a tuple with the AmountType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmountType

`func (o *GroupPolicyScaleOutAction) SetAmountType(v ActionAmount)`

SetAmountType sets AmountType field to given value.


### GetCooldownPeriod

`func (o *GroupPolicyScaleOutAction) GetCooldownPeriod() string`

GetCooldownPeriod returns the CooldownPeriod field if non-nil, zero value otherwise.

### GetCooldownPeriodOk

`func (o *GroupPolicyScaleOutAction) GetCooldownPeriodOk() (*string, bool)`

GetCooldownPeriodOk returns a tuple with the CooldownPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCooldownPeriod

`func (o *GroupPolicyScaleOutAction) SetCooldownPeriod(v string)`

SetCooldownPeriod sets CooldownPeriod field to given value.

### HasCooldownPeriod

`func (o *GroupPolicyScaleOutAction) HasCooldownPeriod() bool`

HasCooldownPeriod returns a boolean if a field has been set.

### SetCooldownPeriodNil

`func (o *GroupPolicyScaleOutAction) SetCooldownPeriodNil(b bool)`

 SetCooldownPeriodNil sets the value for CooldownPeriod to be an explicit nil

### UnsetCooldownPeriod
`func (o *GroupPolicyScaleOutAction) UnsetCooldownPeriod()`

UnsetCooldownPeriod ensures that no value is present for CooldownPeriod, not even an explicit nil

