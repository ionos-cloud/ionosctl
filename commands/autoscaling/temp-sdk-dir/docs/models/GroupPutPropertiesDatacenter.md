# GroupPutPropertiesDatacenter

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |

## Methods

### NewGroupPutPropertiesDatacenter

`func NewGroupPutPropertiesDatacenter(id string, ) *GroupPutPropertiesDatacenter`

NewGroupPutPropertiesDatacenter instantiates a new GroupPutPropertiesDatacenter object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupPutPropertiesDatacenterWithDefaults

`func NewGroupPutPropertiesDatacenterWithDefaults() *GroupPutPropertiesDatacenter`

NewGroupPutPropertiesDatacenterWithDefaults instantiates a new GroupPutPropertiesDatacenter object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *GroupPutPropertiesDatacenter) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GroupPutPropertiesDatacenter) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GroupPutPropertiesDatacenter) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *GroupPutPropertiesDatacenter) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GroupPutPropertiesDatacenter) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GroupPutPropertiesDatacenter) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *GroupPutPropertiesDatacenter) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *GroupPutPropertiesDatacenter) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *GroupPutPropertiesDatacenter) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *GroupPutPropertiesDatacenter) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *GroupPutPropertiesDatacenter) HasHref() bool`

HasHref returns a boolean if a field has been set.


