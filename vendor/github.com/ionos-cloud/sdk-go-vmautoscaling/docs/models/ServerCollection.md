# ServerCollection

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |
|**Items** | Pointer to [**[]ServerResource**](ServerResource.md) |  | [optional] |

## Methods

### NewServerCollection

`func NewServerCollection(id string, ) *ServerCollection`

NewServerCollection instantiates a new ServerCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerCollectionWithDefaults

`func NewServerCollectionWithDefaults() *ServerCollection`

NewServerCollectionWithDefaults instantiates a new ServerCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ServerCollection) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ServerCollection) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ServerCollection) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *ServerCollection) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ServerCollection) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ServerCollection) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ServerCollection) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *ServerCollection) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *ServerCollection) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *ServerCollection) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *ServerCollection) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetItems

`func (o *ServerCollection) GetItems() []ServerResource`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *ServerCollection) GetItemsOk() (*[]ServerResource, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *ServerCollection) SetItems(v []ServerResource)`

SetItems sets Items field to given value.

### HasItems

`func (o *ServerCollection) HasItems() bool`

HasItems returns a boolean if a field has been set.


