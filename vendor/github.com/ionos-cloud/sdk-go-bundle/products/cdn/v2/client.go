/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"golang.org/x/oauth2"
)

var (
	jsonCheck       = regexp.MustCompile(`(?i:(?:application|text)\/(?:vnd\.[^;]+|problem\+)?json)`)
	xmlCheck        = regexp.MustCompile(`(?i:(?:application|text)/xml)`)
	queryParamSplit = regexp.MustCompile(`(^|&)([^&]+)`)
	queryDescape    = strings.NewReplacer("%5B", "[", "%5D", "]")
)

const (
	RequestStatusQueued  = "QUEUED"
	RequestStatusRunning = "RUNNING"
	RequestStatusFailed  = "FAILED"
	RequestStatusDone    = "DONE"

	Version               = "products/cdn/v2.1.2"
	DefaultIonosServerUrl = "https://cdn.de-fra.ionos.com"
	DefaultIonosBasePath  = ""
)

// APIClient manages communication with the IONOS Cloud - CDN Distribution API API v1.2.0
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg    *shared.Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services

	DistributionsApi *DistributionsApiService
}

type service struct {
	client *APIClient
}

func DeepCopy(cfg *shared.Configuration) (*shared.Configuration, error) {
	if cfg == nil {
		return nil, nil
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize configuration: %w", err)
	}

	clone := &shared.Configuration{}
	err = json.Unmarshal(data, clone)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize configuration: %w", err)
	}

	return clone, nil
}

// NewAPIClient creates a new API client. Requires a userAgent string describing your application.
// optionally a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *shared.Configuration) *APIClient {
	// Attempt to deep copy the input configuration. If the configuration contains an httpclient,
	// deepcopy(serialization) will fail. In this case, we fallback to a shallow copy.
	cfgCopy, err := DeepCopy(cfg)
	if err != nil {
		log.Printf("Error creating deep copy of configuration: %v", err)

		// shallow copy instead as a fallback
		cfgCopy = &shared.Configuration{}
		*cfgCopy = *cfg
	}

	cfgCopy.UserAgent = "sdk-go-bundle/products/cdn/v2.1.2"

	// Initialize default values in the copied configuration
	if cfgCopy.HTTPClient == nil {
		cfgCopy.HTTPClient = http.DefaultClient
	}

	if len(cfgCopy.Servers) == 0 {
		cfgCopy.Servers = shared.ServerConfigurations{
			{
				URL:         "https://cdn.de-fra.ionos.com",
				Description: "Frankfurt",
			},
		}
	} else {
		// If the user has provided a custom server configuration, we need to ensure that the basepath is set
		for i := range cfgCopy.Servers {
			if cfgCopy.Servers[i].URL != "" && !strings.HasSuffix(cfgCopy.Servers[i].URL, DefaultIonosBasePath) {
				cfgCopy.Servers[i].URL = fmt.Sprintf("%s%s", cfgCopy.Servers[i].URL, DefaultIonosBasePath)
			}
		}
	}

	// Enable certificate pinning if the environment variable is set
	pkFingerprint := os.Getenv(shared.IonosPinnedCertEnvVar)
	if pkFingerprint != "" {
		httpTransport := &http.Transport{}
		AddPinnedCert(httpTransport, pkFingerprint)
		cfgCopy.HTTPClient.Transport = httpTransport
	}

	// Create and initialize the API client
	c := &APIClient{}
	c.cfg = cfgCopy
	c.common.client = c

	// API Services
	c.DistributionsApi = (*DistributionsApiService)(&c.common)

	return c
}

// AddPinnedCert - enables pinning of the sha256 public fingerprint to the http client's transport
func AddPinnedCert(transport *http.Transport, pkFingerprint string) {
	if pkFingerprint != "" {
		transport.DialTLSContext = addPinnedCertVerification([]byte(pkFingerprint), new(tls.Config))
	}
}

// TLSDial can be assigned to a http.Transport's DialTLS field.
type TLSDial func(ctx context.Context, network, addr string) (net.Conn, error)

// addPinnedCertVerification returns a TLSDial function which checks that
// the remote server provides a certificate whose SHA256 fingerprint matches
// the provided value.
//
// The returned dialer function can be plugged into a http.Transport's DialTLS
// field to allow for certificate pinning.
func addPinnedCertVerification(fingerprint []byte, tlsConfig *tls.Config) TLSDial {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		//fingerprints can be added with ':', we need to trim
		fingerprint = bytes.ReplaceAll(fingerprint, []byte(":"), []byte(""))
		fingerprint = bytes.ReplaceAll(fingerprint, []byte(" "), []byte(""))
		//we are manually checking a certificate, so we need to enable insecure
		tlsConfig.InsecureSkipVerify = true

		// Dial the connection to get certificates to check
		conn, err := tls.Dial(network, addr, tlsConfig)
		if err != nil {
			return nil, err
		}

		if err := verifyPinnedCert(fingerprint, conn.ConnectionState().PeerCertificates); err != nil {
			_ = conn.Close()
			return nil, err
		}

		return conn, nil
	}
}

// verifyPinnedCert iterates the list of peer certificates and attempts to
// locate a certificate that is not a CA and whose public key fingerprint matches pkFingerprint.
func verifyPinnedCert(pkFingerprint []byte, peerCerts []*x509.Certificate) error {
	for _, cert := range peerCerts {
		fingerprint := sha256.Sum256(cert.Raw)

		var bytesFingerPrint = make([]byte, hex.EncodedLen(len(fingerprint[:])))
		hex.Encode(bytesFingerPrint, fingerprint[:])

		// we have a match, and it's not an authority certificate
		if cert.IsCA == false && bytes.EqualFold(bytesFingerPrint, pkFingerprint) {
			return nil
		}
	}

	return fmt.Errorf("remote server presented a certificate which does not match the provided fingerprint")
}

func atoi(in string) (int, error) {
	return strconv.Atoi(in)
}

// selectHeaderContentType select a content type from the available list.
func selectHeaderContentType(contentTypes []string) string {
	if len(contentTypes) == 0 {
		return ""
	}
	if contains(contentTypes, "application/json") {
		return "application/json"
	}
	return contentTypes[0] // use the first content type specified in 'consumes'
}

// selectHeaderAccept join all accept types and return
func selectHeaderAccept(accepts []string) string {
	if len(accepts) == 0 {
		return ""
	}

	if contains(accepts, "application/json") {
		return "application/json"
	}

	return strings.Join(accepts, ",")
}

// contains is a case insenstive match, finding needle in a haystack
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.ToLower(a) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

// Verify optional parameters are of the correct type.
func typeCheckParameter(obj interface{}, expected string, name string) error {
	// Make sure there is an object.
	if obj == nil {
		return nil
	}

	// Check the type is as expected.
	if reflect.TypeOf(obj).String() != expected {
		return fmt.Errorf("Expected %s to be of type %s but received %s.", name, expected, reflect.TypeOf(obj).String())
	}
	return nil
}

// parameterToString convert interface{} parameters to string, using a delimiter if format is provided.
func parameterToString(obj interface{}, collectionFormat string) string {
	var delimiter string

	switch collectionFormat {
	case "pipes":
		delimiter = "|"
	case "ssv":
		delimiter = " "
	case "tsv":
		delimiter = "\t"
	case "csv":
		delimiter = ","
	}

	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		return strings.Trim(strings.Replace(fmt.Sprint(obj), " ", delimiter, -1), "[]")
	} else if t, ok := obj.(time.Time); ok {
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", obj)
}

// helper for converting interface{} parameters to json strings
func parameterToJson(obj interface{}) (string, error) {
	jsonBuf, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonBuf), err
}

func parameterValueToString(obj interface{}, key string) string {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return fmt.Sprintf("%v", obj)
	}
	var param, ok = obj.(MappedNullable)
	if !ok {
		return ""
	}
	dataMap, err := param.ToMap()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", dataMap[key])
}

// parameterAddToHeaderOrQuery adds the provided object to the request header or url query
// supporting deep object syntax
func parameterAddToHeaderOrQuery(headerOrQueryParams interface{}, keyPrefix string, obj interface{}, collectionType string) {
	var v = reflect.ValueOf(obj)
	var value = ""
	if v == reflect.ValueOf(nil) {
		value = "null"
	} else {
		switch v.Kind() {
		case reflect.Invalid:
			value = "invalid"
		case reflect.Struct:
			if t, ok := obj.(MappedNullable); ok {
				dataMap, err := t.ToMap()
				if err != nil {
					return
				}
				parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, dataMap, collectionType)
				return
			}
			if t, ok := obj.(time.Time); ok {
				parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, t.Format(time.RFC3339), collectionType)
				return
			}
			value = v.Type().String() + " value"
		case reflect.Slice:
			var indValue = reflect.ValueOf(obj)
			if indValue == reflect.ValueOf(nil) {
				return
			}
			var lenIndValue = indValue.Len()
			for i := 0; i < lenIndValue; i++ {
				var arrayValue = indValue.Index(i)
				parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, arrayValue.Interface(), collectionType)
			}
			return
		case reflect.Map:
			var indValue = reflect.ValueOf(obj)
			if indValue == reflect.ValueOf(nil) {
				return
			}
			iter := indValue.MapRange()
			for iter.Next() {
				k, v := iter.Key(), iter.Value()
				parameterAddToHeaderOrQuery(headerOrQueryParams, fmt.Sprintf("%s[%s]", keyPrefix, k.String()), v.Interface(), collectionType)
			}
			return
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, v.Elem().Interface(), collectionType)
			return
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(v.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = strconv.FormatUint(v.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(v.Float(), 'g', -1, 32)
		case reflect.Bool:
			value = strconv.FormatBool(v.Bool())
		case reflect.String:
			value = v.String()
		default:
			value = v.Type().String() + " value"
		}
	}
	switch valuesMap := headerOrQueryParams.(type) {
	case url.Values:
		if collectionType == "csv" && valuesMap.Get(keyPrefix) != "" {
			valuesMap.Set(keyPrefix, valuesMap.Get(keyPrefix)+","+value)
		} else {
			valuesMap.Add(keyPrefix, value)
		}
		break
	case map[string]string:
		valuesMap[keyPrefix] = value
		break
	}
}

// callAPI do the request.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, time.Duration, error) {
	retryCount := 0

	var resp *http.Response
	var httpRequestTime time.Duration
	var err error

	for {

		retryCount++

		/* we need to clone the request with every retry time because Body closes after the request */
		var clonedRequest *http.Request = request.Clone(request.Context())
		if request.Body != nil {
			clonedRequest.Body, err = request.GetBody()
			if err != nil {
				return nil, httpRequestTime, err
			}
		}

		if shared.SdkLogLevel.Satisfies(shared.Debug) {
			logRequest := request.Clone(request.Context())

			// Remove the Authorization header if Debug is enabled (but not in Trace mode)
			if !shared.SdkLogLevel.Satisfies(shared.Trace) {
				logRequest.Header.Del("Authorization")
			}

			dump, err := httputil.DumpRequestOut(logRequest, true)
			if err == nil {
				shared.SdkLogger.Printf(" DumpRequestOut : %s\n", string(dump))
			} else {
				shared.SdkLogger.Printf(" DumpRequestOut err: %+v", err)
			}
			shared.SdkLogger.Printf("\n try no: %d\n", retryCount)
		}

		httpRequestStartTime := time.Now()
		clonedRequest.Close = true
		resp, err = c.cfg.HTTPClient.Do(clonedRequest)
		httpRequestTime = time.Since(httpRequestStartTime)
		if err != nil {
			return resp, httpRequestTime, err
		}

		if shared.SdkLogLevel.Satisfies(shared.Trace) {
			dump, err := httputil.DumpResponse(resp, true)
			if err == nil {
				shared.SdkLogger.Printf("\n DumpResponse : %s\n", string(dump))
			} else {
				shared.SdkLogger.Printf(" DumpResponse err %+v", err)
			}
		}

		var backoffTime time.Duration

		switch resp.StatusCode {
		case http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusBadGateway:
			if request.Method == http.MethodPost {
				return resp, httpRequestTime, err
			}
			backoffTime = c.GetConfig().WaitTime

		case http.StatusTooManyRequests:
			if retryAfterSeconds := resp.Header.Get("Retry-After"); retryAfterSeconds != "" {
				waitTime, err := time.ParseDuration(retryAfterSeconds + "s")
				if err != nil {
					return resp, httpRequestTime, err
				}
				backoffTime = waitTime
			} else {
				backoffTime = c.GetConfig().WaitTime
			}
		default:
			return resp, httpRequestTime, err

		}

		if retryCount >= c.GetConfig().MaxRetries {
			if shared.SdkLogLevel.Satisfies(shared.Debug) {
				shared.SdkLogger.Printf(" Number of maximum retries exceeded (%d retries)\n", c.cfg.MaxRetries)
			}
			break
		} else {
			c.backOff(request.Context(), backoffTime)
		}
	}

	return resp, httpRequestTime, err
}

func (c *APIClient) backOff(ctx context.Context, t time.Duration) {
	if t > c.GetConfig().MaxWaitTime {
		t = c.GetConfig().MaxWaitTime
	}
	if shared.SdkLogLevel.Satisfies(shared.Debug) {
		shared.SdkLogger.Printf(" Sleeping %s before retrying request\n", t.String())
	}
	if t <= 0 {
		return
	}
	timer := time.NewTimer(t)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

// Allow modification of underlying config for alternate implementations and testing
// Caution: modifying the configuration while live can cause data races and potentially unwanted behavior
func (c *APIClient) GetConfig() *shared.Configuration {
	return c.cfg
}

// prepareRequest build the request
func (c *APIClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	formFileName string,
	fileName string,
	fileBytes []byte) (localVarRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// add form parameters and file if available.
	if strings.HasPrefix(headerParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 || (len(fileBytes) > 0 && fileName != "") {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and multipart form at the same time.")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					w.WriteField(k, iv)
				}
			}
		}
		if len(fileBytes) > 0 && fileName != "" {
			w.Boundary()
			//_, fileNm := filepath.Split(fileName)
			part, err := w.CreateFormFile(formFileName, filepath.Base(fileName))
			if err != nil {
				return nil, err
			}
			_, err = part.Write(fileBytes)
			if err != nil {
				return nil, err
			}
		}

		// Set the Boundary in the Content-Type
		headerParams["Content-Type"] = w.FormDataContentType()

		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and x-www-form-urlencoded form at the same time.")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Override request host, if applicable
	if c.cfg.Host != "" {
		url.Host = c.cfg.Host
	}

	// Override request scheme, if applicable
	if c.cfg.Scheme != "" {
		url.Scheme = c.cfg.Scheme
	}

	// Adding Query Param
	query := url.Query()
	/* adding default query params */
	for k, v := range c.cfg.DefaultQueryParams {
		if _, ok := queryParams[k]; !ok {
			queryParams[k] = v
		}
	}
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	url.RawQuery = queryParamSplit.ReplaceAllStringFunc(query.Encode(), func(s string) string {
		pieces := strings.Split(s, "=")
		pieces[0] = queryDescape.Replace(pieces[0])
		return strings.Join(pieces, "=")
	})

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, url.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if c.cfg.Token != "" {
		localVarRequest.Header.Add("Authorization", "Bearer "+c.cfg.Token)
	} else {
		if c.cfg.Username != "" {
			localVarRequest.SetBasicAuth(c.cfg.Username, c.cfg.Password)
		}
	}

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)

		// Walk through any authentication.

		// OAuth2 authentication
		if tok, ok := ctx.Value(shared.ContextOAuth2).(oauth2.TokenSource); ok {
			// We were able to grab an oauth2 token from the context
			var latestToken *oauth2.Token
			if latestToken, err = tok.Token(); err != nil {
				return nil, err
			}

			latestToken.SetAuthHeader(localVarRequest)
		}

		// Basic HTTP Authentication
		if auth, ok := ctx.Value(shared.ContextBasicAuth).(shared.BasicAuth); ok {
			localVarRequest.SetBasicAuth(auth.UserName, auth.Password)
		}

		// AccessToken Authentication
		if auth, ok := ctx.Value(shared.ContextAccessToken).(string); ok {
			localVarRequest.Header.Add("Authorization", "Bearer "+auth)
		}

	}

	for header, value := range c.cfg.DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}
	return localVarRequest, nil
}

func (c *APIClient) decode(v interface{}, b []byte, contentType string) (err error) {
	if len(b) == 0 {
		return nil
	}
	if s, ok := v.(*string); ok {
		*s = string(b)
		return nil
	}
	if xmlCheck.MatchString(contentType) {
		if err = xml.Unmarshal(b, v); err != nil {
			return err
		}
		return nil
	}
	if jsonCheck.MatchString(contentType) {
		if actualObj, ok := v.(interface{ GetActualInstance() interface{} }); ok { // oneOf, anyOf schemas
			if unmarshalObj, ok := actualObj.(interface{ UnmarshalJSON([]byte) error }); ok { // make sure it has UnmarshalJSON defined
				if err = unmarshalObj.UnmarshalJSON(b); err != nil {
					return err
				}
			} else {
				return errors.New("unknown type with GetActualInstance but no unmarshalObj.UnmarshalJSON defined")
			}
		} else if err = json.Unmarshal(b, v); err != nil { // simple model
			return err
		}
		return nil
	}
	return fmt.Errorf("undefined response type for content %s", contentType)
}

// Add a file to the multipart request
func addFile(w *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := w.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	return err
}

// Prevent trying to import "fmt"
func reportError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		err = xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("Invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// Ripped from https://github.com/gregjones/httpcache/blob/master/httpcache.go
type cacheControl map[string]string

func parseCacheControl(headers http.Header) cacheControl {
	cc := cacheControl{}
	ccHeader := headers.Get("Cache-Control")
	for _, part := range strings.Split(ccHeader, ",") {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}
		if strings.ContainsRune(part, '=') {
			keyval := strings.Split(part, "=")
			cc[strings.Trim(keyval[0], " ")] = strings.Trim(keyval[1], ",")
		} else {
			cc[part] = ""
		}
	}
	return cc
}

// CacheExpires helper function to determine remaining time before repeating a request.
func CacheExpires(r *http.Response) time.Time {
	// Figure out when the cache expires.
	var expires time.Time
	now, err := time.Parse(time.RFC1123, r.Header.Get("date"))
	if err != nil {
		return time.Now()
	}
	respCacheControl := parseCacheControl(r.Header)

	if maxAge, ok := respCacheControl["max-age"]; ok {
		lifetime, err := time.ParseDuration(maxAge + "s")
		if err != nil {
			expires = now
		} else {
			expires = now.Add(lifetime)
		}
	} else {
		expiresHeader := r.Header.Get("Expires")
		if expiresHeader != "" {
			expires, err = time.Parse(time.RFC1123, expiresHeader)
			if err != nil {
				expires = now
			}
		}
	}
	return expires
}
