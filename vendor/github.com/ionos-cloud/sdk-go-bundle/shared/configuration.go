/*
 * IONOS Shared Libraries
 */

package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var DefaultIonosBasePath = ""

const (
	IonosUsernameEnvVar       = "IONOS_USERNAME"
	IonosPasswordEnvVar       = "IONOS_PASSWORD"
	IonosTokenEnvVar          = "IONOS_TOKEN"
	IonosApiUrlEnvVar         = "IONOS_API_URL"
	IonosPinnedCertEnvVar     = "IONOS_PINNED_CERT"
	IonosLogLevelEnvVar       = "IONOS_LOG_LEVEL"
	IonosFilePathEnvVar       = "IONOS_CONFIG_FILE"
	IonosCurrentProfileEnvVar = "IONOS_CURRENT_PROFILE"
	DefaultIonosServerUrl     = "https://api.ionos.com/"

	defaultMaxRetries   = 3
	defaultWaitTime     = time.Duration(100) * time.Millisecond
	defaultMaxWaitTime  = time.Duration(2000) * time.Millisecond
	defaultPollInterval = 1 * time.Second
)

// contextKeys are used to identify the type of value in the context.
// Since these are string, it is possible to get a short description of the
// context key for logging and debugging using key.String().

type contextKey string

func (c contextKey) String() string {
	return "auth " + string(c)
}

var (
	// ContextOAuth2 takes an oauth2.TokenSource as authentication for the request.
	ContextOAuth2 = contextKey("token")

	// ContextBasicAuth takes BasicAuth as authentication for the request.
	ContextBasicAuth = contextKey("basic")

	// ContextAccessToken takes a string oauth2 access token as authentication for the request.
	ContextAccessToken = contextKey("accesstoken")

	// ContextAPIKeys takes a string apikey as authentication for the request
	ContextAPIKeys = contextKey("apiKeys")

	// ContextHttpSignatureAuth takes HttpSignatureAuth as authentication for the request.
	ContextHttpSignatureAuth = contextKey("httpsignature")

	// ContextServerIndex uses a server configuration from the index.
	ContextServerIndex = contextKey("serverIndex")

	// ContextOperationServerIndices uses a server configuration from the index mapping.
	ContextOperationServerIndices = contextKey("serverOperationIndices")

	// ContextServerVariables overrides a server configuration variables.
	ContextServerVariables = contextKey("serverVariables")

	// ContextOperationServerVariables overrides a server configuration variables using operation specific values.
	ContextOperationServerVariables = contextKey("serverOperationVariables")
)

// BasicAuth provides basic http authentication to a request passed via context using ContextBasicAuth
type BasicAuth struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

// APIKey provides API key based authentication to a request passed via context using ContextAPIKey
type APIKey struct {
	Key    string
	Prefix string
}

// ServerVariable stores the information about a server variable
type ServerVariable struct {
	Description  string
	DefaultValue string
	EnumValues   []string
}

// ServerConfiguration stores the information about a server
type ServerConfiguration struct {
	URL         string
	Description string
	Variables   map[string]ServerVariable
}

// ServerConfigurations stores multiple ServerConfiguration items
type ServerConfigurations []ServerConfiguration

// Configuration stores the configuration of the API client
type Configuration struct {
	Host               string            `json:"host,omitempty"`
	Scheme             string            `json:"scheme,omitempty"`
	DefaultHeader      map[string]string `json:"defaultHeader,omitempty"`
	DefaultQueryParams url.Values        `json:"defaultQueryParams,omitempty"`
	UserAgent          string            `json:"userAgent,omitempty"`
	Servers            ServerConfigurations
	OperationServers   map[string]ServerConfigurations
	HTTPClient         *http.Client
	Username           string        `json:"username,omitempty"`
	Password           string        `json:"password,omitempty"`
	Token              string        `json:"token,omitempty"`
	MaxRetries         int           `json:"maxRetries,omitempty"`
	WaitTime           time.Duration `json:"waitTime,omitempty"`
	MaxWaitTime        time.Duration `json:"maxWaitTime,omitempty"`
	PollInterval       time.Duration `json:"pollInterval,omitempty"`
}

// NewConfiguration returns a new shared.Configuration object
func NewConfiguration(username, password, token, hostUrl string) *Configuration {
	cfg := &Configuration{
		DefaultHeader:      make(map[string]string),
		DefaultQueryParams: url.Values{},
		UserAgent:          "ionos-cloud-sdk-go/v1.0.4",
		Username:           username,
		Password:           password,
		Token:              token,
		MaxRetries:         defaultMaxRetries,
		MaxWaitTime:        defaultMaxWaitTime,
		PollInterval:       defaultPollInterval,
		WaitTime:           defaultWaitTime,
		Servers:            ServerConfigurations{},
		OperationServers:   map[string]ServerConfigurations{},
	}
	if hostUrl != "" {
		cfg.Servers = ServerConfigurations{
			{
				URL:         getServerUrl(hostUrl),
				Description: "Production",
			},
		}
	}
	return cfg
}

// ClientOptions is a struct that represents the client options
type ClientOptions struct {
	// Endpoint is the endpoint that will be overridden
	Endpoint string
	// SkipTLSVerify skips tls verification. Not recommended for production!
	SkipTLSVerify bool
	// Certificate is the certificate that will be used for tls verification
	Certificate string
	// Credentials are the credentials that will be used for authentication
	Credentials Credentials
}

// Credentials are the credentials that will be used for authentication
type Credentials struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token"`
}

// NewConfigurationFromOptions returns a new shared.Configuration object created from the client options
func NewConfigurationFromOptions(clientOptions ClientOptions) *Configuration {
	cfg := &Configuration{
		DefaultHeader:      make(map[string]string),
		DefaultQueryParams: url.Values{},
		UserAgent:          "shared-sdk-go",
		Username:           clientOptions.Credentials.Username,
		Password:           clientOptions.Credentials.Password,
		Token:              clientOptions.Credentials.Token,
		MaxRetries:         defaultMaxRetries,
		MaxWaitTime:        defaultMaxWaitTime,
		WaitTime:           defaultWaitTime,
		Servers:            ServerConfigurations{},
		OperationServers:   map[string]ServerConfigurations{},
		HTTPClient:         http.DefaultClient,
	}
	if clientOptions.Endpoint != "" {
		cfg.Servers = ServerConfigurations{
			{
				URL:         getServerUrl(clientOptions.Endpoint),
				Description: "Production",
			},
		}
	}
	cfg.HTTPClient.Transport = CreateTransport(clientOptions.SkipTLSVerify, clientOptions.Certificate)
	return cfg
}

func CreateTransport(insecure bool, certificate string) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   3,
		MaxConnsPerHost:       3,
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
	if certificate != "" {
		transport.TLSClientConfig.RootCAs = AddCertsToClient(certificate)
	}
	return transport
}

func NewConfigurationFromEnv() *Configuration {
	return NewConfiguration(os.Getenv(IonosUsernameEnvVar), os.Getenv(IonosPasswordEnvVar), os.Getenv(IonosTokenEnvVar), os.Getenv(IonosApiUrlEnvVar))
}

// AddDefaultHeader adds a new HTTP header to the default header in the request
func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}

func (c *Configuration) AddDefaultQueryParam(key string, value string) {
	c.DefaultQueryParams[key] = []string{value}
}

// URL formats template on a index using given variables
func (sc ServerConfigurations) URL(index int, variables map[string]string) (string, error) {
	if index < 0 || len(sc) <= index {
		return "", fmt.Errorf("index %v out of range %v", index, len(sc)-1)
	}
	server := sc[index]
	url := server.URL
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		return "", fmt.Errorf(
			"the URL provided appears to be missing the protocol scheme prefix (\"https://\" or \"http://\"), please verify and try again: %s",
			url,
		)
	}

	// go through variables and replace placeholders
	for name, variable := range server.Variables {
		if value, ok := variables[name]; ok {
			found := bool(len(variable.EnumValues) == 0)
			for _, enumValue := range variable.EnumValues {
				if value == enumValue {
					found = true
				}
			}
			if !found {
				return "", fmt.Errorf("the variable %s in the server URL has invalid value %v. Must be %v", name, value, variable.EnumValues)
			}
			url = strings.Replace(url, "{"+name+"}", value, -1)
		} else {
			url = strings.Replace(url, "{"+name+"}", variable.DefaultValue, -1)
		}
	}
	return url, nil
}

// ServerURL returns URL based on server settings
func (c *Configuration) ServerURL(index int, variables map[string]string) (string, error) {
	return c.Servers.URL(index, variables)
}

func getServerIndex(ctx context.Context) (int, error) {
	si := ctx.Value(ContextServerIndex)
	if si != nil {
		if index, ok := si.(int); ok {
			return index, nil
		}
		return 0, reportError("invalid type %T should be int", si)
	}
	return 0, nil
}

// Prevent trying to import "fmt"
func reportError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func getServerOperationIndex(ctx context.Context, endpoint string) (int, error) {
	osi := ctx.Value(ContextOperationServerIndices)
	if osi != nil {
		if operationIndices, ok := osi.(map[string]int); !ok {
			return 0, reportError("Invalid type %T should be map[string]int", osi)
		} else {
			index, ok := operationIndices[endpoint]
			if ok {
				return index, nil
			}
		}
	}
	return getServerIndex(ctx)
}

func getServerVariables(ctx context.Context) (map[string]string, error) {
	sv := ctx.Value(ContextServerVariables)
	if sv != nil {
		if variables, ok := sv.(map[string]string); ok {
			return variables, nil
		}
		return nil, reportError("ctx value of ContextServerVariables has invalid type %T should be map[string]string", sv)
	}
	return nil, nil
}

func getServerOperationVariables(ctx context.Context, endpoint string) (map[string]string, error) {
	osv := ctx.Value(ContextOperationServerVariables)
	if osv != nil {
		if operationVariables, ok := osv.(map[string]map[string]string); !ok {
			return nil, reportError("ctx value of ContextOperationServerVariables has invalid type %T should be map[string]map[string]string", osv)
		} else {
			variables, ok := operationVariables[endpoint]
			if ok {
				return variables, nil
			}
		}
	}
	return getServerVariables(ctx)
}
func SetBasePath(basePath string) {
	DefaultIonosBasePath = basePath
}
func getServerUrl(serverUrl string) string {
	if serverUrl == "" {
		return DefaultIonosServerUrl + DefaultIonosBasePath
	}

	if !strings.HasSuffix(serverUrl, DefaultIonosBasePath) {
		serverUrl = fmt.Sprintf("%s%s", serverUrl, DefaultIonosBasePath)
	}

	return EnsureURLFormat(serverUrl)
}

// ServerURLWithContext returns a new server URL given an endpoint
func (c *Configuration) ServerURLWithContext(ctx context.Context, endpoint string) (string, error) {
	sc, ok := c.OperationServers[endpoint]
	if !ok {
		sc = c.Servers
	}

	if ctx == nil {
		return sc.URL(0, nil)
	}

	index, err := getServerOperationIndex(ctx, endpoint)
	if err != nil {
		return "", err
	}

	variables, err := getServerOperationVariables(ctx, endpoint)
	if err != nil {
		return "", err
	}

	return sc.URL(index, variables)
}

// ConfigProvider is an interface that allows to get the configuration of shared clients
type ConfigProvider interface {
	GetConfig() *Configuration
}

// EndpointOverridden is a constant that is used to mark the endpoint as overridden and can be used to search for the location
// in the server configuration.
const EndpointOverridden = "endpoint from config file"

// OverrideLocationFor aims to override the server URL for a given client configuration, based on location and endpoint inputs.
// Mutates the client configuration. It searches for the location in the server configuration and overrides the endpoint.
// If the endpoint is empty, it early exits without making changes.
func OverrideLocationFor(configProvider ConfigProvider, location, endpoint string, replaceServers bool) {
	if endpoint == "" {
		return
	}
	// If the replaceServers flag is set, we replace the servers with the new endpoint
	if replaceServers {
		SdkLogger.Printf("[DEBUG] Replacing all server configurations for location %s", location)
		configProvider.GetConfig().Servers = []ServerConfiguration{
			{
				URL:         endpoint,
				Description: EndpointOverridden + location,
			},
		}
		return
	}
	location = strings.TrimSpace(location)
	endpoint = strings.TrimSpace(endpoint)
	servers := configProvider.GetConfig().Servers
	for idx := range servers {
		if strings.Contains(servers[idx].URL, location) {
			SdkLogger.Printf("[DEBUG] Overriding server configuration for location %s", location)
			servers[idx].URL = endpoint
			servers[idx].Description = EndpointOverridden + location
			return
		}
	}
	SdkLogger.Printf("[DEBUG] Adding new server configuration for location %s", location)
	configProvider.GetConfig().Servers = append(configProvider.GetConfig().Servers, ServerConfiguration{
		URL:         endpoint,
		Description: EndpointOverridden + location,
	})
}

func SetSkipTLSVerify(configProvider ConfigProvider, skipTLSVerify bool) {
	configProvider.GetConfig().HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
	}
}

// AddCertsToClient adds certificates to the http client
func AddCertsToClient(authorityData string) *x509.CertPool {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM([]byte(authorityData)); !ok && SdkLogLevel.Satisfies(Debug) {
		SdkLogger.Printf("No certs appended, using system certs only")
	}
	return rootCAs
}
