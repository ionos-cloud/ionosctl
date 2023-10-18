# ErrorAuthorize

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**HttpStatus** | Pointer to **int32** |  | [optional] |
|**Messages** | Pointer to [**[]ErrorAuthorizeMessage**](ErrorAuthorizeMessage.md) |  | [optional] |

## Methods

### NewErrorAuthorize

`func NewErrorAuthorize() *ErrorAuthorize`

NewErrorAuthorize instantiates a new ErrorAuthorize object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorAuthorizeWithDefaults

`func NewErrorAuthorizeWithDefaults() *ErrorAuthorize`

NewErrorAuthorizeWithDefaults instantiates a new ErrorAuthorize object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHttpStatus

`func (o *ErrorAuthorize) GetHttpStatus() int32`

GetHttpStatus returns the HttpStatus field if non-nil, zero value otherwise.

### GetHttpStatusOk

`func (o *ErrorAuthorize) GetHttpStatusOk() (*int32, bool)`

GetHttpStatusOk returns a tuple with the HttpStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpStatus

`func (o *ErrorAuthorize) SetHttpStatus(v int32)`

SetHttpStatus sets HttpStatus field to given value.

### HasHttpStatus

`func (o *ErrorAuthorize) HasHttpStatus() bool`

HasHttpStatus returns a boolean if a field has been set.

### GetMessages

`func (o *ErrorAuthorize) GetMessages() []ErrorAuthorizeMessage`

GetMessages returns the Messages field if non-nil, zero value otherwise.

### GetMessagesOk

`func (o *ErrorAuthorize) GetMessagesOk() (*[]ErrorAuthorizeMessage, bool)`

GetMessagesOk returns a tuple with the Messages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessages

`func (o *ErrorAuthorize) SetMessages(v []ErrorAuthorizeMessage)`

SetMessages sets Messages field to given value.

### HasMessages

`func (o *ErrorAuthorize) HasMessages() bool`

HasMessages returns a boolean if a field has been set.


