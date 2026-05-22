package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/globalwait"

	"github.com/spf13/viper"
)

// hostWithoutPath strips any path from hostUrl; so that SDK clients append their own product paths,
// thus avoiding double basepaths ('/databases/postgresql/cloudapi/v6')
// If for some reason this needs to be removed in the future, then please remove
// the default basepaths in all 'WithConfigOverride' calls too.
func hostWithoutPath(h string) string {
	if h == "" {
		return h
	}
	u, err := url.Parse(h)
	if err != nil || u.Scheme == "" || u.Host == "" {
		// fallback if not a full URL
		return h
	}
	return u.Scheme + "://" + u.Host
}


// resolveProductURL returns the endpoint URL for a given product, respecting
// config file overrides and the global IONOS_API_URL environment variable.
//
// Resolution order:
//  1. IONOS_API_URL env var — overrides ALL products (staging/testing use case)
//  2. Config file per-product override (GetOverride with optional location)
//  3. defaultURL fallback (hardcoded SDK default)
//
// This ensures each SDK client uses its own correct endpoint, eliminating
// cross-SDK URL contamination (e.g., DNS command doesn't affect compute client).
func resolveProductURL(config *fileconfiguration.FileConfig, product, defaultURL string) string {
	// Global override: IONOS_API_URL env var. Overrides ALL products.
	// This is the staging/testing use case where all services run on one host.
	if envURL := os.Getenv(constants.EnvServerUrl); envURL != "" {
		return envURL
	}

	// Config file per-product override (global endpoint, no location).
	// Location-specific overrides are handled at command level via
	// findOverridenURL / NewRegionalConfig, not here.
	if config != nil {
		if ov := config.GetProductGlobalOverrides(product, 0); ov != nil {
			return ov.Name
		}
	}

	return defaultURL
}

// HostWithoutPath strips any path from a URL, returning only scheme://host.
// Exported for use by commands that need to construct product-specific URLs.
func HostWithoutPath(h string) string {
	return hostWithoutPath(h)
}

// NewSharedConfig creates a shared.Configuration for a specific URL,
// applying user agent and fresh HTTPClient. Exported for commands that
// need standalone SDK clients (e.g., login token generation).
func NewSharedConfig(name, pwd, token, url string) *shared.Configuration {
	return newSharedConfig(name, pwd, token, url)
}

// newSharedConfig creates a shared.Configuration for a single product URL,
// applying user agent and fresh HTTPClient.
func newSharedConfig(name, pwd, token, productURL string) *shared.Configuration {
	cfg := shared.NewConfiguration(name, pwd, token, productURL)
	cfg.UserAgent = appendUserAgent(cfg.UserAgent)
	cfg.HTTPClient = &http.Client{} // Prevent mutation of http.DefaultClient
	return cfg
}

func newClient(name, pwd, token string, config *fileconfiguration.FileConfig) *Client {
	// Resolve per-product URLs. Each SDK gets its own correct endpoint.
	// Non-regional products always use api.ionos.com as base.
	// Regional products use their SDK-internal defaults (first location).
	// Config file and IONOS_API_URL can override per-product.
	cloudBase := resolveProductURL(config, fileconfiguration.Cloud, constants.DefaultApiURL)
	authURL := hostWithoutPath(resolveProductURL(config, fileconfiguration.Cloud, constants.DefaultApiURL)) + "/auth/v1"
	registryURL := hostWithoutPath(resolveProductURL(config, fileconfiguration.ContainerRegistry, constants.DefaultApiURL)) + "/containerregistries"
	postgresURL := hostWithoutPath(resolveProductURL(config, fileconfiguration.PSQL, constants.DefaultApiURL)) + "/databases/postgresql"
	mongoURL := hostWithoutPath(resolveProductURL(config, fileconfiguration.Mongo, constants.DefaultApiURL)) + "/databases/mongodb"

	// Regional SDKs: pass empty URL ("") so the SDK uses its built-in default servers.
	// When commands actually execute, they create standalone clients via NewRegionalConfig
	// with the correct location-specific URL. These singleton instances serve as fallbacks
	// and for credential/config access.
	dnsURL := resolveProductURL(config, fileconfiguration.DNS, "")
	cdnURL := resolveProductURL(config, fileconfiguration.CDN, "")
	certURL := resolveProductURL(config, fileconfiguration.Cert, "")
	kafkaURL := resolveProductURL(config, fileconfiguration.Kafka, "")
	loggingURL := resolveProductURL(config, fileconfiguration.Logging, "")
	monitoringURL := resolveProductURL(config, fileconfiguration.Monitoring, "")
	vpnURL := resolveProductURL(config, fileconfiguration.VPN, "")
	mariaURL := resolveProductURL(config, fileconfiguration.Mariadb, "")
	inmemorydbURL := resolveProductURL(config, fileconfiguration.InMemoryDB, "")
	psqlv2URL := resolveProductURL(config, fileconfiguration.PSQLV2, "")
	vmascURL := resolveProductURL(config, fileconfiguration.Autoscaling, constants.DefaultApiURL)

	// Create per-product configs
	cloudUrl := hostWithoutPath(cloudBase) + "/cloudapi/v6"
	clientConfig := cloudv6.NewConfiguration(name, pwd, token, cloudUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)
	clientConfig.HTTPClient = &http.Client{}
	clientConfig.SetDepth(1)

	vmascConfig := vmasc.NewConfiguration(name, pwd, token, vmascURL)
	vmascConfig.UserAgent = appendUserAgent(vmascConfig.UserAgent)
	vmascConfig.HTTPClient = &http.Client{}

	switch v := viper.GetInt(constants.ArgVerbose); {
	case v >= 3:
		shared.SdkLogLevel = shared.Trace
		clientConfig.LogLevel = cloudv6.Trace
		vmascConfig.LogLevel = vmasc.Trace
	case v == 2:
		shared.SdkLogLevel = shared.Debug
		clientConfig.LogLevel = cloudv6.Debug
		vmascConfig.LogLevel = vmasc.Debug
	default:
		// don't explicitly set to Off, as this breaks SDK handling of the IONOS_LOG_LEVEL variable
	}

	// Deprecated --max-results flag: if explicitly passed, use it as --limit
	limit := viper.GetString(constants.FlagLimit)
	if argsContainAny([]string{"--" + constants.DeprecatedFlagMaxResults, "-M"}) {
		limit = viper.GetString(constants.DeprecatedFlagMaxResults)
	}

	queryParams := map[string]string{
		"limit":    limit,
		"offset":   viper.GetString(constants.FlagOffset),
		"depth":    viper.GetString(constants.FlagDepth),
		"order-by": viper.GetString(constants.FlagOrderBy),
	}

	// remove empty values from single-value params
	for k, v := range queryParams {
		if v == "" {
			delete(queryParams, k)
		}
	}

	// Create shared configs for each regional SDK with its own URL
	authConfig := newSharedConfig(name, pwd, token, authURL)
	registryConfig := newSharedConfig(name, pwd, token, registryURL)
	postgresConfig := newSharedConfig(name, pwd, token, postgresURL)
	mongoConfig := newSharedConfig(name, pwd, token, mongoURL)
	dnsConfig := newSharedConfig(name, pwd, token, dnsURL)
	cdnConfig := newSharedConfig(name, pwd, token, cdnURL)
	certConfig := newSharedConfig(name, pwd, token, certURL)
	kafkaConfig := newSharedConfig(name, pwd, token, kafkaURL)
	loggingConfig := newSharedConfig(name, pwd, token, loggingURL)
	monitoringConfig := newSharedConfig(name, pwd, token, monitoringURL)
	vpnConfig := newSharedConfig(name, pwd, token, vpnURL)
	mariaConfig := newSharedConfig(name, pwd, token, mariaURL)
	inmemorydbConfig := newSharedConfig(name, pwd, token, inmemorydbURL)
	psqlv2Config := newSharedConfig(name, pwd, token, psqlv2URL)

	// Collect all configs for batch operations (query params, filters)
	allSharedConfigs := []sdkConfiguration{
		authConfig, registryConfig, postgresConfig, mongoConfig,
		dnsConfig, cdnConfig, certConfig, kafkaConfig, loggingConfig,
		monitoringConfig, vpnConfig, mariaConfig, inmemorydbConfig, psqlv2Config,
		clientConfig, vmascConfig,
	}

	for _, cfg := range allSharedConfigs {
		setQueryParams(cfg, queryParams)
	}

	filterList := viper.GetStringSlice(constants.FlagFilters)
	for _, cfg := range allSharedConfigs {
		setFilters(cfg, filterList)
	}

	c := &Client{
		// api.ionos.com based (non-regional)
		AuthClient:     auth.NewAPIClient(authConfig),
		CloudClient:    cloudv6.NewAPIClient(clientConfig),
		RegistryClient: containerregistry.NewAPIClient(registryConfig),

		PostgresClient:   psql.NewAPIClient(postgresConfig),
		PostgresClientV2: psqlv2.NewAPIClient(psqlv2Config),
		MongoClient:      mongo.NewAPIClient(mongoConfig),

		// Regional APIs — each with its own URL, isolated from others.
		// Commands create standalone clients via NewRegionalConfig for actual calls;
		// these singleton instances provide credential access and serve as fallbacks.
		CDNClient:            cdn.NewAPIClient(cdnConfig),
		CertManagerClient:    cert.NewAPIClient(certConfig),
		DnsClient:            dns.NewAPIClient(dnsConfig),
		Kafka:                kafka.NewAPIClient(kafkaConfig),
		LoggingServiceClient: logging.NewAPIClient(loggingConfig),
		Monitoring:           monitoring.NewAPIClient(monitoringConfig),
		VMAscClient:          vmasc.NewAPIClient(vmascConfig).AutoScalingGroupsApi,
		VPNClient:            vpn.NewAPIClient(vpnConfig),

		MariaClient:      mariadb.NewAPIClient(mariaConfig),
		InMemoryDBClient: inmemorydb.NewAPIClient(inmemorydbConfig),
	}

	// Wrap all SDK HTTP transports so --wait can capture request URLs
	// from any mutating API call (POST/PUT/PATCH/DELETE) automatically.
	// Each SDK client deep-copies its config, so we must wrap each client's
	// HTTPClient individually after construction.
	globalwait.WrapTransport(c.CloudClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.AuthClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.RegistryClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.PostgresClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.PostgresClientV2.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.MongoClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.CDNClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.CertManagerClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.DnsClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.Kafka.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.LoggingServiceClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.Monitoring.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.VPNClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.MariaClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(c.InMemoryDBClient.GetConfig().HTTPClient)
	globalwait.WrapTransport(vmascConfig.HTTPClient)

	return c
}

type sdkConfiguration interface {
	AddDefaultQueryParam(key, val string)
}

// argsContainAny checks if any of the CLI arguments match any of the given names.
// This handles both direct invocations (e.g. "ionosctl image list") and namespaced
// invocations (e.g. "ionosctl compute image list").
func argsContainAny(names []string) bool {
	for _, arg := range os.Args[1:] {
		if slices.Contains(names, arg) {
			return true
		}
	}
	return false
}

func setQueryParams(cfg sdkConfiguration, params map[string]string) {
	for k, v := range params {
		// WARNING: 'images' API expects max-results instead of limit
		// TODO: Instead of 'os.Args': 'commands.GetRootCmd().Command.CommandPath()'. But, causes import cycles. After refactor, change this.
		if k == "limit" && argsContainAny([]string{"image", "img"}) {
			if !viper.IsSet(constants.FlagLimit) && !argsContainAny([]string{"--" + constants.DeprecatedFlagMaxResults, "-M"}) {
				// do NOT apply the default value of 'limit' in this case
				// because 'maxResults' is applied before filtering
				// while 'limit' is applied after filtering
				// which leads to some incredible confusion
				// as for why everything handles differently on this command only
				continue
			}
			cfg.AddDefaultQueryParam("maxResults", v)
			continue
		}

		if k == "depth" && argsContainAny([]string{"logging-service", "log-svc", "monitoring"}) {
			// Logging and Monitoring APIs do not support 'depth'
			continue
		}

		cfg.AddDefaultQueryParam(k, v)
	}
}

// setFilters applies multiple filter query params. Each entry in filters must be "key=value".
// If the same key appears multiple times, values are joined with commas.
func setFilters(cfg sdkConfiguration, filters []string) {
	if len(filters) == 0 {
		return
	}
	grouped := make(map[string][]string)
	for _, f := range filters {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}
		key := normalizeFilterKey(parts[0])
		grouped[key] = append(grouped[key], parts[1])
	}
	for k, vals := range grouped {
		key := fmt.Sprintf("filter.%s", k)
		cfg.AddDefaultQueryParam(key, strings.Join(vals, ","))
	}
}
