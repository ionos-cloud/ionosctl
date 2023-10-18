# GroupPostResponse

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |
|**Metadata** | Pointer to [**Metadata**](Metadata.md) |  | [optional] |
|**Properties** | [**GroupProperties**](GroupProperties.md) |  | |
|**Entities** | Pointer to [**GroupPostEntities**](GroupPostEntities.md) |  | [optional] |
|**StartedActions** | Pointer to [**[]ActionResource**](ActionResource.md) | Any background activity caused by this request. You can use this to track the progress of such activities. | [optional] [readonly] |

## Methods

### NewGroupPostResponse

`func NewGroupPostResponse(id string, properties GroupProperties, ) *GroupPostResponse`

NewGroupPostResponse instantiates a new GroupPostResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPostResponseWithDefaults

`func NewGroupPostResponseWithDefaults() *GroupPostResponse`

NewGroupPostResponseWithDefaults instantiates a new GroupPostResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *GroupPostResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GroupPostResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GroupPostResponse) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *GroupPostResponse) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GroupPostResponse) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GroupPostResponse) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *GroupPostResponse) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *GroupPostResponse) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *GroupPostResponse) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *GroupPostResponse) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *GroupPostResponse) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetMetadata

`func (o *GroupPostResponse) GetMetadata() Metadata`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *GroupPostResponse) GetMetadataOk() (*Metadata, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *GroupPostResponse) SetMetadata(v Metadata)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *GroupPostResponse) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetProperties

`func (o *GroupPostResponse) GetProperties() GroupProperties`

GetProperties returns the Properties field if non-nil, zero value otherwise.

### GetPropertiesOk

`func (o *GroupPostResponse) GetPropertiesOk() (*GroupProperties, bool)`

GetPropertiesOk returns a tuple with the Properties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProperties

`func (o *GroupPostResponse) SetProperties(v GroupProperties)`

SetProperties sets Properties field to given value.


### GetEntities

`func (o *GroupPostResponse) GetEntities() GroupPostEntities`

GetEntities returns the Entities field if non-nil, zero value otherwise.

### GetEntitiesOk

`func (o *GroupPostResponse) GetEntitiesOk() (*GroupPostEntities, bool)`

GetEntitiesOk returns a tuple with the Entities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntities

`func (o *GroupPostResponse) SetEntities(v GroupPostEntities)`

SetEntities sets Entities field to given value.

### HasEntities

`func (o *GroupPostResponse) HasEntities() bool`

HasEntities returns a boolean if a field has been set.

### GetStartedActions

`func (o *GroupPostResponse) GetStartedActions() []ActionResource`

GetStartedActions returns the StartedActions field if non-nil, zero value otherwise.

### GetStartedActionsOk

`func (o *GroupPostResponse) GetStartedActionsOk() (*[]ActionResource, bool)`

GetStartedActionsOk returns a tuple with the StartedActions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedActions

`func (o *GroupPostResponse) SetStartedActions(v []ActionResource)`

SetStartedActions sets StartedActions field to given value.

### HasStartedActions

`func (o *GroupPostResponse) HasStartedActions() bool`

HasStartedActions returns a boolean if a field has been set.


