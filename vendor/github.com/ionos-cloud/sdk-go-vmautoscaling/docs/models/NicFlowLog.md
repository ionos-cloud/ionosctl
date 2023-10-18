# NicFlowLog

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Name** | **string** | The flow log name. | |
|**Action** | **string** | Specifies the traffic action pattern. | |
|**Direction** | **string** | Specifies the traffic direction pattern. | |
|**Bucket** | **string** | The S3 bucket name of an existing IONOS Cloud S3 bucket. | |

## Methods

### NewNicFlowLog

`func NewNicFlowLog(name string, action string, direction string, bucket string, ) *NicFlowLog`

NewNicFlowLog instantiates a new NicFlowLog object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNicFlowLogWithDefaults

`func NewNicFlowLogWithDefaults() *NicFlowLog`

NewNicFlowLogWithDefaults instantiates a new NicFlowLog object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NicFlowLog) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NicFlowLog) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NicFlowLog) SetName(v string)`

SetName sets Name field to given value.


### GetAction

`func (o *NicFlowLog) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *NicFlowLog) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *NicFlowLog) SetAction(v string)`

SetAction sets Action field to given value.


### GetDirection

`func (o *NicFlowLog) GetDirection() string`

GetDirection returns the Direction field if non-nil, zero value otherwise.

### GetDirectionOk

`func (o *NicFlowLog) GetDirectionOk() (*string, bool)`

GetDirectionOk returns a tuple with the Direction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirection

`func (o *NicFlowLog) SetDirection(v string)`

SetDirection sets Direction field to given value.


### GetBucket

`func (o *NicFlowLog) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *NicFlowLog) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *NicFlowLog) SetBucket(v string)`

SetBucket sets Bucket field to given value.



