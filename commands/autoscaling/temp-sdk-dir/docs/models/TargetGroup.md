# TargetGroup

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**TargetGroupId** | **string** | id | |
|**Port** | **int32** | port | |
|**Weight** | **int32** | weight | |

## Methods

### NewTargetGroup

`func NewTargetGroup(targetGroupId string, port int32, weight int32, ) *TargetGroup`

NewTargetGroup instantiates a new TargetGroup object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTargetGroupWithDefaults

`func NewTargetGroupWithDefaults() *TargetGroup`

NewTargetGroupWithDefaults instantiates a new TargetGroup object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTargetGroupId

`func (o *TargetGroup) GetTargetGroupId() string`

GetTargetGroupId returns the TargetGroupId field if non-nil, zero value otherwise.

### GetTargetGroupIdOk

`func (o *TargetGroup) GetTargetGroupIdOk() (*string, bool)`

GetTargetGroupIdOk returns a tuple with the TargetGroupId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetGroupId

`func (o *TargetGroup) SetTargetGroupId(v string)`

SetTargetGroupId sets TargetGroupId field to given value.


### GetPort

`func (o *TargetGroup) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *TargetGroup) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *TargetGroup) SetPort(v int32)`

SetPort sets Port field to given value.


### GetWeight

`func (o *TargetGroup) GetWeight() int32`

GetWeight returns the Weight field if non-nil, zero value otherwise.

### GetWeightOk

`func (o *TargetGroup) GetWeightOk() (*int32, bool)`

GetWeightOk returns a tuple with the Weight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWeight

`func (o *TargetGroup) SetWeight(v int32)`

SetWeight sets Weight field to given value.



