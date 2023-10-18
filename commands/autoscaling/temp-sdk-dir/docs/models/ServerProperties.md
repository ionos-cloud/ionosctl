# ServerProperties

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**DatacenterServer** | [**DatacenterServer**](DatacenterServer.md) |  | |
|**Name** | Pointer to **string** |  | [optional] |

## Methods

### NewServerProperties

`func NewServerProperties(datacenterServer DatacenterServer, ) *ServerProperties`

NewServerProperties instantiates a new ServerProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerPropertiesWithDefaults

`func NewServerPropertiesWithDefaults() *ServerProperties`

NewServerPropertiesWithDefaults instantiates a new ServerProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDatacenterServer

`func (o *ServerProperties) GetDatacenterServer() DatacenterServer`

GetDatacenterServer returns the DatacenterServer field if non-nil, zero value otherwise.

### GetDatacenterServerOk

`func (o *ServerProperties) GetDatacenterServerOk() (*DatacenterServer, bool)`

GetDatacenterServerOk returns a tuple with the DatacenterServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDatacenterServer

`func (o *ServerProperties) SetDatacenterServer(v DatacenterServer)`

SetDatacenterServer sets DatacenterServer field to given value.


### GetName

`func (o *ServerProperties) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServerProperties) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServerProperties) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServerProperties) HasName() bool`

HasName returns a boolean if a field has been set.


