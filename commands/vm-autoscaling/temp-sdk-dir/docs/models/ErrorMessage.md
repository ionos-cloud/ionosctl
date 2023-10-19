# ErrorMessage

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ErrorCode** | Pointer to **string** |  | [optional] |
|**Message** | Pointer to **string** |  | [optional] |

## Methods

### NewErrorMessage

`func NewErrorMessage() *ErrorMessage`

NewErrorMessage instantiates a new ErrorMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorMessageWithDefaults

`func NewErrorMessageWithDefaults() *ErrorMessage`

NewErrorMessageWithDefaults instantiates a new ErrorMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorCode

`func (o *ErrorMessage) GetErrorCode() string`

GetErrorCode returns the ErrorCode field if non-nil, zero value otherwise.

### GetErrorCodeOk

`func (o *ErrorMessage) GetErrorCodeOk() (*string, bool)`

GetErrorCodeOk returns a tuple with the ErrorCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorCode

`func (o *ErrorMessage) SetErrorCode(v string)`

SetErrorCode sets ErrorCode field to given value.

### HasErrorCode

`func (o *ErrorMessage) HasErrorCode() bool`

HasErrorCode returns a boolean if a field has been set.

### GetMessage

`func (o *ErrorMessage) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ErrorMessage) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ErrorMessage) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ErrorMessage) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


