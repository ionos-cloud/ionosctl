# GroupPost

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |
|**Metadata** | Pointer to [**Metadata**](Metadata.md) |  | [optional] |
|**Properties** | [**GroupProperties**](GroupProperties.md) |  | |
|**Entities** | Pointer to [**GroupPostEntities**](GroupPostEntities.md) |  | [optional] |

## Methods

### NewGroupPost

`func NewGroupPost(id string, properties GroupProperties, ) *GroupPost`

NewGroupPost instantiates a new GroupPost object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPostWithDefaults

`func NewGroupPostWithDefaults() *GroupPost`

NewGroupPostWithDefaults instantiates a new GroupPost object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *GroupPost) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GroupPost) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GroupPost) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *GroupPost) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GroupPost) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GroupPost) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *GroupPost) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *GroupPost) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *GroupPost) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *GroupPost) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *GroupPost) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetMetadata

`func (o *GroupPost) GetMetadata() Metadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *GroupPost) GetMetadataOk() (*Metadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *GroupPost) SetMetadata(v Metadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *GroupPost) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetProperties

`func (o *GroupPost) GetProperties() GroupProperties`

GetProperties returns the Properties field if non-nil, zero value otherwise.

### GetPropertiesOk

`func (o *GroupPost) GetPropertiesOk() (*GroupProperties, bool)`

GetPropertiesOk returns a tuple with the Properties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProperties

`func (o *GroupPost) SetProperties(v GroupProperties)`

SetProperties sets Properties field to given value.


### GetEntities

`func (o *GroupPost) GetEntities() GroupPostEntities`

GetEntities returns the Entities field if non-nil, zero value otherwise.

### GetEntitiesOk

`func (o *GroupPost) GetEntitiesOk() (*GroupPostEntities, bool)`

GetEntitiesOk returns a tuple with the Entities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntities

`func (o *GroupPost) SetEntities(v GroupPostEntities)`

SetEntities sets Entities field to given value.

### HasEntities

`func (o *GroupPost) HasEntities() bool`

HasEntities returns a boolean if a field has been set.


