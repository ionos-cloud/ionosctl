# GroupPostEntities

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Actions** | Pointer to [**ActionsLinkResource**](ActionsLinkResource.md) |  | [optional] |
|**Servers** | Pointer to [**ServersLinkResource**](ServersLinkResource.md) |  | [optional] |

## Methods

### NewGroupPostEntities

`func NewGroupPostEntities() *GroupPostEntities`

NewGroupPostEntities instantiates a new GroupPostEntities object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPostEntitiesWithDefaults

`func NewGroupPostEntitiesWithDefaults() *GroupPostEntities`

NewGroupPostEntitiesWithDefaults instantiates a new GroupPostEntities object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActions

`func (o *GroupPostEntities) GetActions() ActionsLinkResource`

GetActions returns the Actions field if non-nil, zero value otherwise.

### GetActionsOk

`func (o *GroupPostEntities) GetActionsOk() (*ActionsLinkResource, bool)`

GetActionsOk returns a tuple with the Actions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActions

`func (o *GroupPostEntities) SetActions(v ActionsLinkResource)`

SetActions sets Actions field to given value.

### HasActions

`func (o *GroupPostEntities) HasActions() bool`

HasActions returns a boolean if a field has been set.

### GetServers

`func (o *GroupPostEntities) GetServers() ServersLinkResource`

GetServers returns the Servers field if non-nil, zero value otherwise.

### GetServersOk

`func (o *GroupPostEntities) GetServersOk() (*ServersLinkResource, bool)`

GetServersOk returns a tuple with the Servers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServers

`func (o *GroupPostEntities) SetServers(v ServersLinkResource)`

SetServers sets Servers field to given value.

### HasServers

`func (o *GroupPostEntities) HasServers() bool`

HasServers returns a boolean if a field has been set.


