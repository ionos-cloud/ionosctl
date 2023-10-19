# ServerResource

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |

## Methods

### NewServerResource

`func NewServerResource(id string, ) *ServerResource`

NewServerResource instantiates a new ServerResource object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerResourceWithDefaults

`func NewServerResourceWithDefaults() *ServerResource`

NewServerResourceWithDefaults instantiates a new ServerResource object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ServerResource) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ServerResource) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ServerResource) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *ServerResource) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ServerResource) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ServerResource) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ServerResource) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *ServerResource) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *ServerResource) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *ServerResource) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *ServerResource) HasHref() bool`

HasHref returns a boolean if a field has been set.


