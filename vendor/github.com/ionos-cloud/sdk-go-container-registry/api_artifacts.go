/*
 * Container Registry service
 *
 * ## Overview Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls. ## Changelog ### 1.1.0  - Added new endpoints for Repositories  - Added new endpoints for Artifacts  - Added new endpoints for Vulnerabilities  - Added registry vulnerabilityScanning feature
 *
 * API version: 1.1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	_context "context"
	"fmt"
	"io"
	_nethttp "net/http"
	_neturl "net/url"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

// ArtifactsApiService ArtifactsApi service
type ArtifactsApiService service

type ApiRegistriesArtifactsGetRequest struct {
	ctx                   _context.Context
	ApiService            *ArtifactsApiService
	registryId            string
	offset                *int32
	limit                 *int32
	filterVulnerabilityId *string
	orderBy               *string
}

func (r ApiRegistriesArtifactsGetRequest) Offset(offset int32) ApiRegistriesArtifactsGetRequest {
	r.offset = &offset
	return r
}
func (r ApiRegistriesArtifactsGetRequest) Limit(limit int32) ApiRegistriesArtifactsGetRequest {
	r.limit = &limit
	return r
}
func (r ApiRegistriesArtifactsGetRequest) FilterVulnerabilityId(filterVulnerabilityId string) ApiRegistriesArtifactsGetRequest {
	r.filterVulnerabilityId = &filterVulnerabilityId
	return r
}
func (r ApiRegistriesArtifactsGetRequest) OrderBy(orderBy string) ApiRegistriesArtifactsGetRequest {
	r.orderBy = &orderBy
	return r
}

func (r ApiRegistriesArtifactsGetRequest) Execute() (RegistryArtifactsReadList, *APIResponse, error) {
	return r.ApiService.RegistriesArtifactsGetExecute(r)
}

/*
  - RegistriesArtifactsGet Retrieve all Artifacts by Registry
  - This endpoint enables retrieving all Artifacts using

pagination and optional filters.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param registryId The ID (UUID) of the Registry.
  - @return ApiRegistriesArtifactsGetRequest
*/
func (a *ArtifactsApiService) RegistriesArtifactsGet(ctx _context.Context, registryId string) ApiRegistriesArtifactsGetRequest {
	return ApiRegistriesArtifactsGetRequest{
		ApiService: a,
		ctx:        ctx,
		registryId: registryId,
	}
}

/*
 * Execute executes the request
 * @return RegistryArtifactsReadList
 */
func (a *ArtifactsApiService) RegistriesArtifactsGetExecute(r ApiRegistriesArtifactsGetRequest) (RegistryArtifactsReadList, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  RegistryArtifactsReadList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ArtifactsApiService.RegistriesArtifactsGet")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/registries/{registryId}/artifacts"
	localVarPath = strings.Replace(localVarPath, "{"+"registryId"+"}", _neturl.PathEscape(parameterToString(r.registryId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.offset != nil {
		localVarQueryParams.Add("offset", parameterToString(*r.offset, ""))
	}
	if r.limit != nil {
		localVarQueryParams.Add("limit", parameterToString(*r.limit, ""))
	}
	if r.filterVulnerabilityId != nil {
		localVarQueryParams.Add("filter.vulnerabilityId", parameterToString(*r.filterVulnerabilityId, ""))
	}
	if r.orderBy != nil {
		localVarQueryParams.Add("orderBy", parameterToString(*r.orderBy, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["tokenAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "RegistriesArtifactsGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiRegistriesRepositoriesArtifactsFindByDigestRequest struct {
	ctx            _context.Context
	ApiService     *ArtifactsApiService
	registryId     string
	repositoryName string
	digest         string
}

func (r ApiRegistriesRepositoriesArtifactsFindByDigestRequest) Execute() (ArtifactRead, *APIResponse, error) {
	return r.ApiService.RegistriesRepositoriesArtifactsFindByDigestExecute(r)
}

/*
 * RegistriesRepositoriesArtifactsFindByDigest Retrieve Artifact
 * Returns the Artifact by Digest.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param registryId The ID (UUID) of the Registry.
 * @param repositoryName The Name of the Repository.
 * @param digest The Digest of the Artifact that should be retrieved.
 * @return ApiRegistriesRepositoriesArtifactsFindByDigestRequest
 */
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsFindByDigest(ctx _context.Context, registryId string, repositoryName string, digest string) ApiRegistriesRepositoriesArtifactsFindByDigestRequest {
	return ApiRegistriesRepositoriesArtifactsFindByDigestRequest{
		ApiService:     a,
		ctx:            ctx,
		registryId:     registryId,
		repositoryName: repositoryName,
		digest:         digest,
	}
}

/*
 * Execute executes the request
 * @return ArtifactRead
 */
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsFindByDigestExecute(r ApiRegistriesRepositoriesArtifactsFindByDigestRequest) (ArtifactRead, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  ArtifactRead
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ArtifactsApiService.RegistriesRepositoriesArtifactsFindByDigest")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/registries/{registryId}/repositories/{repositoryName}/artifacts/{digest}"
	localVarPath = strings.Replace(localVarPath, "{"+"registryId"+"}", _neturl.PathEscape(parameterToString(r.registryId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repositoryName"+"}", _neturl.PathEscape(parameterToString(r.repositoryName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"digest"+"}", _neturl.PathEscape(parameterToString(r.digest, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if strlen(r.repositoryName) > 256 {
		return localVarReturnValue, nil, reportError("repositoryName must have less than 256 elements")
	}
	if strlen(r.digest) > 128 {
		return localVarReturnValue, nil, reportError("digest must have less than 128 elements")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["tokenAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "RegistriesRepositoriesArtifactsFindByDigest",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiRegistriesRepositoriesArtifactsGetRequest struct {
	ctx            _context.Context
	ApiService     *ArtifactsApiService
	registryId     string
	repositoryName string
	offset         *int32
	limit          *int32
	orderBy        *string
}

func (r ApiRegistriesRepositoriesArtifactsGetRequest) Offset(offset int32) ApiRegistriesRepositoriesArtifactsGetRequest {
	r.offset = &offset
	return r
}
func (r ApiRegistriesRepositoriesArtifactsGetRequest) Limit(limit int32) ApiRegistriesRepositoriesArtifactsGetRequest {
	r.limit = &limit
	return r
}
func (r ApiRegistriesRepositoriesArtifactsGetRequest) OrderBy(orderBy string) ApiRegistriesRepositoriesArtifactsGetRequest {
	r.orderBy = &orderBy
	return r
}

func (r ApiRegistriesRepositoriesArtifactsGetRequest) Execute() (ArtifactReadList, *APIResponse, error) {
	return r.ApiService.RegistriesRepositoriesArtifactsGetExecute(r)
}

/*
  - RegistriesRepositoriesArtifactsGet Retrieve all Artifacts by Repository
  - This endpoint enables retrieving all Artifacts using

pagination and optional filters.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param registryId The ID (UUID) of the Registry.
  - @param repositoryName The Name of the Repository.
  - @return ApiRegistriesRepositoriesArtifactsGetRequest
*/
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsGet(ctx _context.Context, registryId string, repositoryName string) ApiRegistriesRepositoriesArtifactsGetRequest {
	return ApiRegistriesRepositoriesArtifactsGetRequest{
		ApiService:     a,
		ctx:            ctx,
		registryId:     registryId,
		repositoryName: repositoryName,
	}
}

/*
 * Execute executes the request
 * @return ArtifactReadList
 */
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsGetExecute(r ApiRegistriesRepositoriesArtifactsGetRequest) (ArtifactReadList, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  ArtifactReadList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ArtifactsApiService.RegistriesRepositoriesArtifactsGet")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/registries/{registryId}/repositories/{repositoryName}/artifacts"
	localVarPath = strings.Replace(localVarPath, "{"+"registryId"+"}", _neturl.PathEscape(parameterToString(r.registryId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repositoryName"+"}", _neturl.PathEscape(parameterToString(r.repositoryName, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if strlen(r.repositoryName) > 256 {
		return localVarReturnValue, nil, reportError("repositoryName must have less than 256 elements")
	}

	if r.offset != nil {
		localVarQueryParams.Add("offset", parameterToString(*r.offset, ""))
	}
	if r.limit != nil {
		localVarQueryParams.Add("limit", parameterToString(*r.limit, ""))
	}
	if r.orderBy != nil {
		localVarQueryParams.Add("orderBy", parameterToString(*r.orderBy, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["tokenAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "RegistriesRepositoriesArtifactsGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest struct {
	ctx            _context.Context
	ApiService     *ArtifactsApiService
	registryId     string
	repositoryName string
	digest         string
	offset         *int32
	limit          *int32
	filterSeverity *string
	filterFixable  *bool
	orderBy        *string
}

func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) Offset(offset int32) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	r.offset = &offset
	return r
}
func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) Limit(limit int32) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	r.limit = &limit
	return r
}
func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) FilterSeverity(filterSeverity string) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	r.filterSeverity = &filterSeverity
	return r
}
func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) FilterFixable(filterFixable bool) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	r.filterFixable = &filterFixable
	return r
}
func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) OrderBy(orderBy string) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	r.orderBy = &orderBy
	return r
}

func (r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) Execute() (ArtifactVulnerabilityReadList, *APIResponse, error) {
	return r.ApiService.RegistriesRepositoriesArtifactsVulnerabilitiesGetExecute(r)
}

/*
  - RegistriesRepositoriesArtifactsVulnerabilitiesGet Retrieve all Vulnerabilities
  - This endpoint enables retrieving all Vulnerabilities using

pagination and optional filters.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param registryId The ID (UUID) of the Registry.
  - @param repositoryName The Name of the Repository.
  - @param digest The Digest of the Artifact.
  - @return ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest
*/
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsVulnerabilitiesGet(ctx _context.Context, registryId string, repositoryName string, digest string) ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest {
	return ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest{
		ApiService:     a,
		ctx:            ctx,
		registryId:     registryId,
		repositoryName: repositoryName,
		digest:         digest,
	}
}

/*
 * Execute executes the request
 * @return ArtifactVulnerabilityReadList
 */
func (a *ArtifactsApiService) RegistriesRepositoriesArtifactsVulnerabilitiesGetExecute(r ApiRegistriesRepositoriesArtifactsVulnerabilitiesGetRequest) (ArtifactVulnerabilityReadList, *APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  ArtifactVulnerabilityReadList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ArtifactsApiService.RegistriesRepositoriesArtifactsVulnerabilitiesGet")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/registries/{registryId}/repositories/{repositoryName}/artifacts/{digest}/vulnerabilities"
	localVarPath = strings.Replace(localVarPath, "{"+"registryId"+"}", _neturl.PathEscape(parameterToString(r.registryId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repositoryName"+"}", _neturl.PathEscape(parameterToString(r.repositoryName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"digest"+"}", _neturl.PathEscape(parameterToString(r.digest, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if strlen(r.repositoryName) > 256 {
		return localVarReturnValue, nil, reportError("repositoryName must have less than 256 elements")
	}
	if strlen(r.digest) > 128 {
		return localVarReturnValue, nil, reportError("digest must have less than 128 elements")
	}

	if r.offset != nil {
		localVarQueryParams.Add("offset", parameterToString(*r.offset, ""))
	}
	if r.limit != nil {
		localVarQueryParams.Add("limit", parameterToString(*r.limit, ""))
	}
	if r.filterSeverity != nil {
		localVarQueryParams.Add("filter.severity", parameterToString(*r.filterSeverity, ""))
	}
	if r.filterFixable != nil {
		localVarQueryParams.Add("filter.fixable", parameterToString(*r.filterFixable, ""))
	}
	if r.orderBy != nil {
		localVarQueryParams.Add("orderBy", parameterToString(*r.orderBy, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["tokenAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "RegistriesRepositoriesArtifactsVulnerabilitiesGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)),
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			statusCode: localVarHTTPResponse.StatusCode,
			body:       localVarBody,
			error:      err.Error(),
		}
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}
