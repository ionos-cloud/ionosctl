# Error404Message

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ErrorCode** | Pointer to **string** |  | [optional] |
|**Message** | Pointer to **string** |  | [optional] |

## Methods

### NewError404Message

`func NewError404Message() *Error404Message`

NewError404Message instantiates a new Error404Message object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewError404MessageWithDefaults

`func NewError404MessageWithDefaults() *Error404Message`

NewError404MessageWithDefaults instantiates a new Error404Message object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorCode

`func (o *Error404Message) GetErrorCode() string`

GetErrorCode returns the ErrorCode field if non-nil, zero value otherwise.

### GetErrorCodeOk

`func (o *Error404Message) GetErrorCodeOk() (*string, bool)`

GetErrorCodeOk returns a tuple with the ErrorCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCode

`func (o *Error404Message) SetErrorCode(v string)`

SetErrorCode sets ErrorCode field to given value.

### HasErrorCode

`func (o *Error404Message) HasErrorCode() bool`

HasErrorCode returns a boolean if a field has been set.

### GetMessage

`func (o *Error404Message) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *Error404Message) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *Error404Message) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *Error404Message) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


