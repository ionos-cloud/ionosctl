/*
 * IONOS Logging REST API
 *
 * Logging Service is a service that provides a centralized logging system where users are able to push and aggregate their system or application logs. This service also provides a visualization platform where users are able to observe, search and filter the logs and also create dashboards and alerts for their data points. This service can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an API. The API allows you to create logging pipelines or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"log"
	"net/http"
	"time"
)

// APIResponse stores the API response returned by the server.
type APIResponse struct {
	*http.Response `json:"-"`
	Message        string `json:"message,omitempty"`
	// Operation is the name of the OpenAPI operation.
	Operation string `json:"operation,omitempty"`
	// RequestURL is the request URL. This value is always available, even if the
	// embedded *http.Response is nil.
	RequestURL string `json:"url,omitempty"`
	// RequestTime is the time duration from the moment the APIClient sends
	// the HTTP request to the moment it receives an HTTP response.
	RequestTime time.Duration `json:"duration,omitempty"`
	// Method is the HTTP method used for the request.  This value is always
	// available, even if the embedded *http.Response is nil.
	Method string `json:"method,omitempty"`
	// Payload holds the contents of the response body (which may be nil or empty).
	// This is provided here as the raw response.Body() reader will have already
	// been drained.
	Payload []byte `json:"-"`
}

// NewAPIResponse returns a new APIResponse object.
func NewAPIResponse(r *http.Response) *APIResponse {

	response := &APIResponse{Response: r}
	return response
}

// NewAPIResponseWithError returns a new APIResponse object with the provided error message.
func NewAPIResponseWithError(errorMessage string) *APIResponse {

	response := &APIResponse{Message: errorMessage}
	return response
}

// HttpNotFound - returns true if a 404 status code was returned
// returns false for nil APIResponse values
func (resp *APIResponse) HttpNotFound() bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

// LogInfo - logs APIResponse values like RequestTime, Operation and StatusCode
// does not print anything for nil APIResponse values
func (resp *APIResponse) LogInfo() {
	if resp != nil {
		log.Printf("[DEBUG] Request time : %s for operation : %s",
			resp.RequestTime, resp.Operation)
		if resp.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", resp.StatusCode)
		}
	}
}
