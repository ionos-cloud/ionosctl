[![CI](https://github.com/ionos-cloud/ionosctl/workflows/CI/badge.svg)](https://github.com/ionos-cloud/sdk-go-auth/actions)
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=alert_status)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=bugs)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=security_rating)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=sdk-go-auth&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=sdk-go-auth)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/sdk-go-auth.svg)](https://github.com/ionos-cloud/sdk-go-auth/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/ionos-cloud/sdk-go-auth.svg)](https://github.com/ionos-cloud/sdk-go-auth/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/sdk-go-auth.svg)](https://github.com/ionos-cloud/sdk-go-auth)

![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")

# Go Auth API client for ionoscloud

The IONOS Cloud SDK Auth for GO provides you with access to the IONOS Cloud Auth API. Use the Auth API to manage tokens for secure access to IONOS Cloud APIs (Auth API, Cloud API, Reseller API, Activity Log API, and others).

## Overview

This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project. By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 1.0
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Getting Started

An IONOS account is required for access to the Cloud API; credentials from your registration are used to authenticate against the IONOS Cloud API.

## Installation

Install the Go language from the official [Go installation](https://golang.org/doc/install) guide.

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
go get github.com/antihax/optional
```

Put the package under your project folder by adding it to the `go.mod`:

```golang
require "github.com/ionos-cloud/sdk-go-auth"
```

You can get the package locally using:

```golang
go get "github.com/ionos-cloud/sdk-go-auth"
```

## Authentication

The username and password or the authentication token can be manually specified when initializing the sdk client:

```golang

client := ionossdk.NewAPIClient(ionossdk.NewConfiguration(username, password, token, hostUrl))
```

Environment variables can also be used. The sdk uses the following variables:

- IONOS_USERNAME - to specify the username used to log in
- IONOS_PASSWORD - to specify the password
- IONOS_TOKEN - if an authentication token is being used
- IONOS_API_URL - to specify the API endpoint in order to overwrite it

In this case, the client configuration needs to be initialized using `NewConfigurationFromEnv()`

```golang

client := ionossdk.NewAPIClient(ionossdk.NewConfigurationFromEnv())

```

## Environment Variables

Environment Variable | Description
--- | --- 
`IONOS_USERNAME` | Specify the username used to login, to authenticate against the IONOS Cloud API | 
`IONOS_PASSWORD` | Specify the password used to login, to authenticate against the IONOS Cloud API | 
`IONOS_TOKEN` | Specify the token used to login, if a token is being used instead of username and password |
`IONOS_API_URL` | Specify the API URL. It will overwrite the API endpoint default value `api.ionos.com`. Note: the host URL does not contain the `/cloudapi/v5` path, so it should _not_ be included in the `IONOS_API_URL` environment variable | 

## Documentation for API Endpoints

All paths are relative to *https://api.ionos.com/auth/v1*

### TokenApi

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateToken**](docs/api/TokenApi.md#createtoken) | **Get** /tokens/generate | Create new tokens
[**DeleteTokenByCriteria**](docs/api/TokenApi.md#deletetokenbycriteria) | **Delete** /tokens | Delete tokens by criteria
[**DeleteTokenById**](docs/api/TokenApi.md#deletetokenbyid) | **Delete** /tokens/{tokenId} | Delete tokens
[**GetAllTokens**](docs/api/TokenApi.md#getalltokens) | **Get** /tokens | List all tokens
[**GetTokenById**](docs/api/TokenApi.md#gettokenbyid) | **Get** /tokens/{tokenId} | Get tokens by Key ID

## Documentation For Models

- [DeleteResponse](docs/models/DeleteResponse.md)
- [Error](docs/models/Error.md)
- [ErrorMessages](docs/models/ErrorMessages.md)
- [Jwt](docs/models/Jwt.md)
- [Token](docs/models/Token.md)
- [Tokens](docs/models/Tokens.md)

## FAQ

* How can I open a bug report/feature request?

Bug reports and feature requests can be opened in the Issues repository: [https://github.com/ionos-cloud/sdk-go-auth/issues/new/choose](https://github.com/ionos-cloud/sdk-go-auth/issues/new/choose)

* Can I contribute to the GO SDK for Auth API?

Pure SDKs are automatically generated using OpenAPI Generator and don’t support manual changes. If you require changes, please open an issue, and we will try to address it.
