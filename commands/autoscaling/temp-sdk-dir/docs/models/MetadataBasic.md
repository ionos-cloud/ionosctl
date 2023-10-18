# MetadataBasic

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**CreatedDate** | [**time.Time**](time.Time.md) | The date the resource was created. | |
|**Etag** | **string** | The resource etag. | |
|**LastModifiedDate** | [**time.Time**](time.Time.md) | The date the resource was last modified. | |
|**State** | **string** | The resource state. | |

## Methods

### NewMetadataBasic

`func NewMetadataBasic(createdDate time.Time, etag string, lastModifiedDate time.Time, state string, ) *MetadataBasic`

NewMetadataBasic instantiates a new MetadataBasic object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMetadataBasicWithDefaults

`func NewMetadataBasicWithDefaults() *MetadataBasic`

NewMetadataBasicWithDefaults instantiates a new MetadataBasic object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedDate

`func (o *MetadataBasic) GetCreatedDate() time.Time`

GetCreatedDate returns the CreatedDate field if non-nil, zero value otherwise.

### GetCreatedDateOk

`func (o *MetadataBasic) GetCreatedDateOk() (*time.Time, bool)`

GetCreatedDateOk returns a tuple with the CreatedDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedDate

`func (o *MetadataBasic) SetCreatedDate(v time.Time)`

SetCreatedDate sets CreatedDate field to given value.


### GetEtag

`func (o *MetadataBasic) GetEtag() string`

GetEtag returns the Etag field if non-nil, zero value otherwise.

### GetEtagOk

`func (o *MetadataBasic) GetEtagOk() (*string, bool)`

GetEtagOk returns a tuple with the Etag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEtag

`func (o *MetadataBasic) SetEtag(v string)`

SetEtag sets Etag field to given value.


### GetLastModifiedDate

`func (o *MetadataBasic) GetLastModifiedDate() time.Time`

GetLastModifiedDate returns the LastModifiedDate field if non-nil, zero value otherwise.

### GetLastModifiedDateOk

`func (o *MetadataBasic) GetLastModifiedDateOk() (*time.Time, bool)`

GetLastModifiedDateOk returns a tuple with the LastModifiedDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModifiedDate

`func (o *MetadataBasic) SetLastModifiedDate(v time.Time)`

SetLastModifiedDate sets LastModifiedDate field to given value.


### GetState

`func (o *MetadataBasic) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *MetadataBasic) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *MetadataBasic) SetState(v string)`

SetState sets State field to given value.



