# ErrorReplicaValidateMessage

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ErrorCode** | Pointer to **string** |  | [optional] |
|**Message** | Pointer to **string** |  | [optional] |

## Methods

### NewErrorReplicaValidateMessage

`func NewErrorReplicaValidateMessage() *ErrorReplicaValidateMessage`

NewErrorReplicaValidateMessage instantiates a new ErrorReplicaValidateMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorReplicaValidateMessageWithDefaults

`func NewErrorReplicaValidateMessageWithDefaults() *ErrorReplicaValidateMessage`

NewErrorReplicaValidateMessageWithDefaults instantiates a new ErrorReplicaValidateMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorCode

`func (o *ErrorReplicaValidateMessage) GetErrorCode() string`

GetErrorCode returns the ErrorCode field if non-nil, zero value otherwise.

### GetErrorCodeOk

`func (o *ErrorReplicaValidateMessage) GetErrorCodeOk() (*string, bool)`

GetErrorCodeOk returns a tuple with the ErrorCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCode

`func (o *ErrorReplicaValidateMessage) SetErrorCode(v string)`

SetErrorCode sets ErrorCode field to given value.

### HasErrorCode

`func (o *ErrorReplicaValidateMessage) HasErrorCode() bool`

HasErrorCode returns a boolean if a field has been set.

### GetMessage

`func (o *ErrorReplicaValidateMessage) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ErrorReplicaValidateMessage) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ErrorReplicaValidateMessage) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ErrorReplicaValidateMessage) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


