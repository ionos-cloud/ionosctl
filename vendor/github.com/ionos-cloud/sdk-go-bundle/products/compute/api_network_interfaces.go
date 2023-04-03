/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package compute

import (
	_context "context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	_neturl "net/url"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

// NetworkInterfacesApiService NetworkInterfacesApi service
type NetworkInterfacesApiService service

type ApiDatacentersServersNicsDeleteRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	datacenterId    string
	serverId        string
	nicId           string
	pretty          *bool
	depth           *int32
	xContractNumber *int32
}

func (r ApiDatacentersServersNicsDeleteRequest) Pretty(pretty bool) ApiDatacentersServersNicsDeleteRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsDeleteRequest) Depth(depth int32) ApiDatacentersServersNicsDeleteRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsDeleteRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsDeleteRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiDatacentersServersNicsDeleteRequest) Execute() (*shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsDeleteExecute(r)
}

/*
 * DatacentersServersNicsDelete Delete NICs
 * Remove the specified NIC.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @param nicId The unique ID of the NIC.
 * @return ApiDatacentersServersNicsDeleteRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsDelete(ctx _context.Context, datacenterId string, serverId string, nicId string) ApiDatacentersServersNicsDeleteRequest {
	return ApiDatacentersServersNicsDeleteRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
		nicId:        nicId,
	}
}

/*
 * Execute executes the request
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsDeleteExecute(r ApiDatacentersServersNicsDeleteRequest) (*shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodDelete
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsDelete")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics/{nicId}"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nicId"+"}", _neturl.PathEscape(parameterToString(r.nicId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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
		return nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsDelete",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarAPIResponse, newErr
	}

	return localVarAPIResponse, nil
}

type ApiDatacentersServersNicsFindByIdRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	datacenterId    string
	serverId        string
	nicId           string
	pretty          *bool
	depth           *int32
	xContractNumber *int32
}

func (r ApiDatacentersServersNicsFindByIdRequest) Pretty(pretty bool) ApiDatacentersServersNicsFindByIdRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsFindByIdRequest) Depth(depth int32) ApiDatacentersServersNicsFindByIdRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsFindByIdRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsFindByIdRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiDatacentersServersNicsFindByIdRequest) Execute() (Nic, *shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsFindByIdExecute(r)
}

/*
 * DatacentersServersNicsFindById Retrieve NICs
 * Retrieve the properties of the specified NIC.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @param nicId The unique ID of the NIC.
 * @return ApiDatacentersServersNicsFindByIdRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsFindById(ctx _context.Context, datacenterId string, serverId string, nicId string) ApiDatacentersServersNicsFindByIdRequest {
	return ApiDatacentersServersNicsFindByIdRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
		nicId:        nicId,
	}
}

/*
 * Execute executes the request
 * @return Nic
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsFindByIdExecute(r ApiDatacentersServersNicsFindByIdRequest) (Nic, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  Nic
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsFindById")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics/{nicId}"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nicId"+"}", _neturl.PathEscape(parameterToString(r.nicId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsFindById",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiDatacentersServersNicsGetRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	filters         _neturl.Values
	orderBy         *string
	maxResults      *int32
	datacenterId    string
	serverId        string
	pretty          *bool
	depth           *int32
	xContractNumber *int32
	offset          *int32
	limit           *int32
}

func (r ApiDatacentersServersNicsGetRequest) Pretty(pretty bool) ApiDatacentersServersNicsGetRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsGetRequest) Depth(depth int32) ApiDatacentersServersNicsGetRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsGetRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsGetRequest {
	r.xContractNumber = &xContractNumber
	return r
}
func (r ApiDatacentersServersNicsGetRequest) Offset(offset int32) ApiDatacentersServersNicsGetRequest {
	r.offset = &offset
	return r
}
func (r ApiDatacentersServersNicsGetRequest) Limit(limit int32) ApiDatacentersServersNicsGetRequest {
	r.limit = &limit
	return r
}

// Filters query parameters limit results to those containing a matching value for a specific property.
func (r ApiDatacentersServersNicsGetRequest) Filter(key string, value string) ApiDatacentersServersNicsGetRequest {
	filterKey := fmt.Sprintf("filter.%s", key)
	r.filters[filterKey] = append(r.filters[filterKey], value)
	return r
}

// OrderBy query param sorts the results alphanumerically in ascending order based on the specified property.
func (r ApiDatacentersServersNicsGetRequest) OrderBy(orderBy string) ApiDatacentersServersNicsGetRequest {
	r.orderBy = &orderBy
	return r
}

// MaxResults query param limits the number of results returned.
func (r ApiDatacentersServersNicsGetRequest) MaxResults(maxResults int32) ApiDatacentersServersNicsGetRequest {
	r.maxResults = &maxResults
	return r
}

func (r ApiDatacentersServersNicsGetRequest) Execute() (Nics, *shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsGetExecute(r)
}

/*
 * DatacentersServersNicsGet List NICs
 * List all NICs, attached to the specified server.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @return ApiDatacentersServersNicsGetRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsGet(ctx _context.Context, datacenterId string, serverId string) ApiDatacentersServersNicsGetRequest {
	return ApiDatacentersServersNicsGetRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
		filters:      _neturl.Values{},
	}
}

/*
 * Execute executes the request
 * @return Nics
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsGetExecute(r ApiDatacentersServersNicsGetRequest) (Nics, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  Nics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsGet")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
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
	if r.maxResults != nil {
		localVarQueryParams.Add("maxResults", parameterToString(*r.maxResults, ""))
	}
	if len(r.filters) > 0 {
		for k, v := range r.filters {
			for _, iv := range v {
				localVarQueryParams.Add(k, iv)
			}
		}
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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiDatacentersServersNicsPatchRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	datacenterId    string
	serverId        string
	nicId           string
	nic             *NicProperties
	pretty          *bool
	depth           *int32
	xContractNumber *int32
}

func (r ApiDatacentersServersNicsPatchRequest) Nic(nic NicProperties) ApiDatacentersServersNicsPatchRequest {
	r.nic = &nic
	return r
}
func (r ApiDatacentersServersNicsPatchRequest) Pretty(pretty bool) ApiDatacentersServersNicsPatchRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsPatchRequest) Depth(depth int32) ApiDatacentersServersNicsPatchRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsPatchRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsPatchRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiDatacentersServersNicsPatchRequest) Execute() (Nic, *shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsPatchExecute(r)
}

/*
 * DatacentersServersNicsPatch Partially modify NICs
 * Update the properties of the specified NIC.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @param nicId The unique ID of the NIC.
 * @return ApiDatacentersServersNicsPatchRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPatch(ctx _context.Context, datacenterId string, serverId string, nicId string) ApiDatacentersServersNicsPatchRequest {
	return ApiDatacentersServersNicsPatchRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
		nicId:        nicId,
	}
}

/*
 * Execute executes the request
 * @return Nic
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPatchExecute(r ApiDatacentersServersNicsPatchRequest) (Nic, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPatch
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  Nic
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsPatch")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics/{nicId}"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nicId"+"}", _neturl.PathEscape(parameterToString(r.nicId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.nic == nil {
		return localVarReturnValue, nil, reportError("nic is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.nic
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsPatch",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiDatacentersServersNicsPostRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	datacenterId    string
	serverId        string
	nic             *Nic
	pretty          *bool
	depth           *int32
	xContractNumber *int32
}

func (r ApiDatacentersServersNicsPostRequest) Nic(nic Nic) ApiDatacentersServersNicsPostRequest {
	r.nic = &nic
	return r
}
func (r ApiDatacentersServersNicsPostRequest) Pretty(pretty bool) ApiDatacentersServersNicsPostRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsPostRequest) Depth(depth int32) ApiDatacentersServersNicsPostRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsPostRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsPostRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiDatacentersServersNicsPostRequest) Execute() (Nic, *shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsPostExecute(r)
}

/*
 * DatacentersServersNicsPost Create NICs
 * Add a NIC to the specified server. The combined total of NICs and attached volumes cannot exceed 24 per server.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @return ApiDatacentersServersNicsPostRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPost(ctx _context.Context, datacenterId string, serverId string) ApiDatacentersServersNicsPostRequest {
	return ApiDatacentersServersNicsPostRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
	}
}

/*
 * Execute executes the request
 * @return Nic
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPostExecute(r ApiDatacentersServersNicsPostRequest) (Nic, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPost
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  Nic
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsPost")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.nic == nil {
		return localVarReturnValue, nil, reportError("nic is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.nic
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsPost",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiDatacentersServersNicsPutRequest struct {
	ctx             _context.Context
	ApiService      *NetworkInterfacesApiService
	datacenterId    string
	serverId        string
	nicId           string
	nic             *NicPut
	pretty          *bool
	depth           *int32
	xContractNumber *int32
}

func (r ApiDatacentersServersNicsPutRequest) Nic(nic NicPut) ApiDatacentersServersNicsPutRequest {
	r.nic = &nic
	return r
}
func (r ApiDatacentersServersNicsPutRequest) Pretty(pretty bool) ApiDatacentersServersNicsPutRequest {
	r.pretty = &pretty
	return r
}
func (r ApiDatacentersServersNicsPutRequest) Depth(depth int32) ApiDatacentersServersNicsPutRequest {
	r.depth = &depth
	return r
}
func (r ApiDatacentersServersNicsPutRequest) XContractNumber(xContractNumber int32) ApiDatacentersServersNicsPutRequest {
	r.xContractNumber = &xContractNumber
	return r
}

func (r ApiDatacentersServersNicsPutRequest) Execute() (Nic, *shared.APIResponse, error) {
	return r.ApiService.DatacentersServersNicsPutExecute(r)
}

/*
 * DatacentersServersNicsPut Modify NICs
 * Modify the properties of the specified NIC.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param datacenterId The unique ID of the data center.
 * @param serverId The unique ID of the server.
 * @param nicId The unique ID of the NIC.
 * @return ApiDatacentersServersNicsPutRequest
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPut(ctx _context.Context, datacenterId string, serverId string, nicId string) ApiDatacentersServersNicsPutRequest {
	return ApiDatacentersServersNicsPutRequest{
		ApiService:   a,
		ctx:          ctx,
		datacenterId: datacenterId,
		serverId:     serverId,
		nicId:        nicId,
	}
}

/*
 * Execute executes the request
 * @return Nic
 */
func (a *NetworkInterfacesApiService) DatacentersServersNicsPutExecute(r ApiDatacentersServersNicsPutRequest) (Nic, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPut
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  Nic
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkInterfacesApiService.DatacentersServersNicsPut")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/datacenters/{datacenterId}/servers/{serverId}/nics/{nicId}"
	localVarPath = strings.Replace(localVarPath, "{"+"datacenterId"+"}", _neturl.PathEscape(parameterToString(r.datacenterId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"serverId"+"}", _neturl.PathEscape(parameterToString(r.serverId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nicId"+"}", _neturl.PathEscape(parameterToString(r.nicId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.nic == nil {
		return localVarReturnValue, nil, reportError("nic is required and must be specified")
	}

	if r.pretty != nil {
		localVarQueryParams.Add("pretty", parameterToString(*r.pretty, ""))
	}
	if r.depth != nil {
		localVarQueryParams.Add("depth", parameterToString(*r.depth, ""))
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	if r.xContractNumber != nil {
		localVarHeaderParams["X-Contract-Number"] = parameterToString(*r.xContractNumber, "")
	}
	// body params
	localVarPostBody = r.nic
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(shared.ContextAPIKeys).(map[string]shared.APIKey); ok {
			if apiKey, ok := auth["Token Authentication"]; ok {
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

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DatacentersServersNicsPut",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}
