# GroupCollection

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |
|**Items** | Pointer to [**[]GroupResource**](GroupResource.md) |  | [optional] |

## Methods

### NewGroupCollection

`func NewGroupCollection(id string, ) *GroupCollection`

NewGroupCollection instantiates a new GroupCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupCollectionWithDefaults

`func NewGroupCollectionWithDefaults() *GroupCollection`

NewGroupCollectionWithDefaults instantiates a new GroupCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *GroupCollection) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GroupCollection) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GroupCollection) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *GroupCollection) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GroupCollection) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GroupCollection) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *GroupCollection) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *GroupCollection) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *GroupCollection) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *GroupCollection) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *GroupCollection) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetItems

`func (o *GroupCollection) GetItems() []GroupResource`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *GroupCollection) GetItemsOk() (*[]GroupResource, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *GroupCollection) SetItems(v []GroupResource)`

SetItems sets Items field to given value.

### HasItems

`func (o *GroupCollection) HasItems() bool`

HasItems returns a boolean if a field has been set.


