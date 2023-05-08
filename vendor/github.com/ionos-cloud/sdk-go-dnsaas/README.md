# Go API client for ionoscloud

DNS API Specification


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.1.0
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [https://docs.ionos.com/support/general-information/contact-information](https://docs.ionos.com/support/general-information/contact-information)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import sw "./ionoscloud"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), sw.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identifield by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```
ctx := context.WithValue(context.Background(), sw.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), sw.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://dns.de-fra.ionos.com*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*RecordsApi* | [**RecordsGet**](docs/api/RecordsApi.md#recordsget) | **Get** /records | Retrieve all records
*RecordsApi* | [**ZonesRecordsDelete**](docs/api/RecordsApi.md#zonesrecordsdelete) | **Delete** /zones/{zoneId}/records/{recordId} | Delete a record
*RecordsApi* | [**ZonesRecordsFindById**](docs/api/RecordsApi.md#zonesrecordsfindbyid) | **Get** /zones/{zoneId}/records/{recordId} | Retrieve a record
*RecordsApi* | [**ZonesRecordsGet**](docs/api/RecordsApi.md#zonesrecordsget) | **Get** /zones/{zoneId}/records | Retrieve records
*RecordsApi* | [**ZonesRecordsPost**](docs/api/RecordsApi.md#zonesrecordspost) | **Post** /zones/{zoneId}/records | Create a record
*RecordsApi* | [**ZonesRecordsPut**](docs/api/RecordsApi.md#zonesrecordsput) | **Put** /zones/{zoneId}/records/{recordId} | Ensure a record
*ZonesApi* | [**ZonesDelete**](docs/api/ZonesApi.md#zonesdelete) | **Delete** /zones/{zoneId} | Delete a zone
*ZonesApi* | [**ZonesFindById**](docs/api/ZonesApi.md#zonesfindbyid) | **Get** /zones/{zoneId} | Retrieve a zone
*ZonesApi* | [**ZonesGet**](docs/api/ZonesApi.md#zonesget) | **Get** /zones | Retrieve zones
*ZonesApi* | [**ZonesPost**](docs/api/ZonesApi.md#zonespost) | **Post** /zones | Create a zone
*ZonesApi* | [**ZonesPut**](docs/api/ZonesApi.md#zonesput) | **Put** /zones/{zoneId} | Ensure a zone


## Documentation For Models

 - [ErrorMessage](docs/models/ErrorMessage.md)
 - [ErrorResponse](docs/models/ErrorResponse.md)
 - [ProvisioningState](docs/models/ProvisioningState.md)
 - [RecordCreateRequest](docs/models/RecordCreateRequest.md)
 - [RecordMetadata](docs/models/RecordMetadata.md)
 - [RecordProperties](docs/models/RecordProperties.md)
 - [RecordResponse](docs/models/RecordResponse.md)
 - [RecordType](docs/models/RecordType.md)
 - [RecordUpdateRequest](docs/models/RecordUpdateRequest.md)
 - [RecordsResponse](docs/models/RecordsResponse.md)
 - [ZoneCreateRequest](docs/models/ZoneCreateRequest.md)
 - [ZoneCreateRequestProperties](docs/models/ZoneCreateRequestProperties.md)
 - [ZoneResponse](docs/models/ZoneResponse.md)
 - [ZoneResponseMetadata](docs/models/ZoneResponseMetadata.md)
 - [ZoneResponseProperties](docs/models/ZoneResponseProperties.md)
 - [ZoneUpdateRequest](docs/models/ZoneUpdateRequest.md)
 - [ZoneUpdateRequestProperties](docs/models/ZoneUpdateRequestProperties.md)
 - [ZonesResponse](docs/models/ZonesResponse.md)


## Documentation For Authorization



### tokenAuth

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: Authorization and passed in as the auth context for each request.


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

support@cloud.ionos.com
