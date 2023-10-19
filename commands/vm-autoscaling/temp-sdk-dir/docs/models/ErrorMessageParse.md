# ErrorMessageParse

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ErrorCode** | Pointer to **string** |  | [optional] |
|**Message** | Pointer to **string** |  | [optional] |

## Methods

### NewErrorMessageParse

`func NewErrorMessageParse() *ErrorMessageParse`

NewErrorMessageParse instantiates a new ErrorMessageParse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorMessageParseWithDefaults

`func NewErrorMessageParseWithDefaults() *ErrorMessageParse`

NewErrorMessageParseWithDefaults instantiates a new ErrorMessageParse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorCode

`func (o *ErrorMessageParse) GetErrorCode() string`

GetErrorCode returns the ErrorCode field if non-nil, zero value otherwise.

### GetErrorCodeOk

`func (o *ErrorMessageParse) GetErrorCodeOk() (*string, bool)`

GetErrorCodeOk returns a tuple with the ErrorCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCode

`func (o *ErrorMessageParse) SetErrorCode(v string)`

SetErrorCode sets ErrorCode field to given value.

### HasErrorCode

`func (o *ErrorMessageParse) HasErrorCode() bool`

HasErrorCode returns a boolean if a field has been set.

### GetMessage

`func (o *ErrorMessageParse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ErrorMessageParse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ErrorMessageParse) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ErrorMessageParse) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


