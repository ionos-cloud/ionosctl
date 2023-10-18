# Server

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | **string** | The unique resource identifier. | [readonly] |
|**Type** | Pointer to **string** | The resource type. | [optional] [readonly] |
|**Href** | Pointer to **string** | The absolute URL to the resource&#39;s representation. | [optional] [readonly] |
|**Metadata** | Pointer to [**MetadataBasic**](MetadataBasic.md) |  | [optional] |
|**Properties** | Pointer to [**ServerProperties**](ServerProperties.md) |  | [optional] |

## Methods

### NewServer

`func NewServer(id string, ) *Server`

NewServer instantiates a new Server object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerWithDefaults

`func NewServerWithDefaults() *Server`

NewServerWithDefaults instantiates a new Server object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Server) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Server) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Server) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *Server) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Server) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Server) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *Server) HasType() bool`

HasType returns a boolean if a field has been set.

### GetHref

`func (o *Server) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *Server) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *Server) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *Server) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetMetadata

`func (o *Server) GetMetadata() MetadataBasic`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *Server) GetMetadataOk() (*MetadataBasic, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *Server) SetMetadata(v MetadataBasic)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *Server) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### GetProperties

`func (o *Server) GetProperties() ServerProperties`

GetProperties returns the Properties field if non-nil, zero value otherwise.

### GetPropertiesOk

`func (o *Server) GetPropertiesOk() (*ServerProperties, bool)`

GetPropertiesOk returns a tuple with the Properties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProperties

`func (o *Server) SetProperties(v ServerProperties)`

SetProperties sets Properties field to given value.

### HasProperties

`func (o *Server) HasProperties() bool`

HasProperties returns a boolean if a field has been set.


