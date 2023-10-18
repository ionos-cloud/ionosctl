# GroupPolicy

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Metric** | [**Metric**](Metric.md) |  | |
|**Range** | Pointer to **string** | Specifies the time range for which the samples are to be aggregated. Must be &gt;&#x3D; 2 minutes. | [optional] [default to "120s"]|
|**ScaleInAction** | [**GroupPolicyScaleInAction**](GroupPolicyScaleInAction.md) |  | |
|**ScaleInThreshold** | **float32** | The lower threshold for the value of the &#39;metric&#39;. Used with the &#x60;less than&#x60; (&lt;) operator. When this value is exceeded, a scale-in action is triggered, specified by the &#39;scaleInAction&#39; property. The value must have a higher minimum delta to the &#39;scaleOutThreshold&#39;, depending on the &#39;metric&#39;, to avoid competing for actions at the same time. | |
|**ScaleOutAction** | [**GroupPolicyScaleOutAction**](GroupPolicyScaleOutAction.md) |  | |
|**ScaleOutThreshold** | **float32** | The upper threshold for the value of the &#39;metric&#39;. Used with the &#39;greater than&#39; (&gt;) operator. A scale-out action is triggered when this value is exceeded, specified by the &#39;scaleOutAction&#39; property. The value must have a lower minimum delta to the &#39;scaleInThreshold&#39;, depending on the metric, to avoid competing for actions simultaneously. If &#39;properties.policy.unit&#x3D;TOTAL&#39;, a value &gt;&#x3D; 40 must be chosen. | |
|**Unit** | [**QueryUnit**](QueryUnit.md) |  | [default to QUERYUNIT_TOTAL]|

## Methods

### NewGroupPolicy

`func NewGroupPolicy(metric Metric, scaleInAction GroupPolicyScaleInAction, scaleInThreshold float32, scaleOutAction GroupPolicyScaleOutAction, scaleOutThreshold float32, unit QueryUnit, ) *GroupPolicy`

NewGroupPolicy instantiates a new GroupPolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPolicyWithDefaults

`func NewGroupPolicyWithDefaults() *GroupPolicy`

NewGroupPolicyWithDefaults instantiates a new GroupPolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMetric

`func (o *GroupPolicy) GetMetric() Metric`

GetMetric returns the Metric field if non-nil, zero value otherwise.

### GetMetricOk

`func (o *GroupPolicy) GetMetricOk() (*Metric, bool)`

GetMetricOk returns a tuple with the Metric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetric

`func (o *GroupPolicy) SetMetric(v Metric)`

SetMetric sets Metric field to given value.


### GetRange

`func (o *GroupPolicy) GetRange() string`

GetRange returns the Range field if non-nil, zero value otherwise.

### GetRangeOk

`func (o *GroupPolicy) GetRangeOk() (*string, bool)`

GetRangeOk returns a tuple with the Range field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRange

`func (o *GroupPolicy) SetRange(v string)`

SetRange sets Range field to given value.

### HasRange

`func (o *GroupPolicy) HasRange() bool`

HasRange returns a boolean if a field has been set.

### GetScaleInAction

`func (o *GroupPolicy) GetScaleInAction() GroupPolicyScaleInAction`

GetScaleInAction returns the ScaleInAction field if non-nil, zero value otherwise.

### GetScaleInActionOk

`func (o *GroupPolicy) GetScaleInActionOk() (*GroupPolicyScaleInAction, bool)`

GetScaleInActionOk returns a tuple with the ScaleInAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScaleInAction

`func (o *GroupPolicy) SetScaleInAction(v GroupPolicyScaleInAction)`

SetScaleInAction sets ScaleInAction field to given value.


### GetScaleInThreshold

`func (o *GroupPolicy) GetScaleInThreshold() float32`

GetScaleInThreshold returns the ScaleInThreshold field if non-nil, zero value otherwise.

### GetScaleInThresholdOk

`func (o *GroupPolicy) GetScaleInThresholdOk() (*float32, bool)`

GetScaleInThresholdOk returns a tuple with the ScaleInThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScaleInThreshold

`func (o *GroupPolicy) SetScaleInThreshold(v float32)`

SetScaleInThreshold sets ScaleInThreshold field to given value.


### GetScaleOutAction

`func (o *GroupPolicy) GetScaleOutAction() GroupPolicyScaleOutAction`

GetScaleOutAction returns the ScaleOutAction field if non-nil, zero value otherwise.

### GetScaleOutActionOk

`func (o *GroupPolicy) GetScaleOutActionOk() (*GroupPolicyScaleOutAction, bool)`

GetScaleOutActionOk returns a tuple with the ScaleOutAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScaleOutAction

`func (o *GroupPolicy) SetScaleOutAction(v GroupPolicyScaleOutAction)`

SetScaleOutAction sets ScaleOutAction field to given value.


### GetScaleOutThreshold

`func (o *GroupPolicy) GetScaleOutThreshold() float32`

GetScaleOutThreshold returns the ScaleOutThreshold field if non-nil, zero value otherwise.

### GetScaleOutThresholdOk

`func (o *GroupPolicy) GetScaleOutThresholdOk() (*float32, bool)`

GetScaleOutThresholdOk returns a tuple with the ScaleOutThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScaleOutThreshold

`func (o *GroupPolicy) SetScaleOutThreshold(v float32)`

SetScaleOutThreshold sets ScaleOutThreshold field to given value.


### GetUnit

`func (o *GroupPolicy) GetUnit() QueryUnit`

GetUnit returns the Unit field if non-nil, zero value otherwise.

### GetUnitOk

`func (o *GroupPolicy) GetUnitOk() (*QueryUnit, bool)`

GetUnitOk returns a tuple with the Unit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnit

`func (o *GroupPolicy) SetUnit(v QueryUnit)`

SetUnit sets Unit field to given value.



