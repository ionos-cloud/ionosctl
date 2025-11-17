package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

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

func configGuaranteeBasepath(cfg *shared.Configuration, defaultBasepath string) *shared.Configuration {
	var url string
	if len(cfg.Servers) > 0 {
		url = hostWithoutPath(cfg.Servers[0].URL)
	} else {
		// fallback
		url = constants.DefaultApiURL
	}
	return shared.NewConfiguration(cfg.Username, cfg.Password, cfg.Token, url+defaultBasepath)
}

func newClient(name, pwd, token, hostUrl string) *Client {
	sharedConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	sharedConfig.UserAgent = appendUserAgent(sharedConfig.UserAgent)

	cloudUrl := hostWithoutPath(hostUrl) + "/cloudapi/v6"
	clientConfig := cloudv6.NewConfiguration(name, pwd, token, cloudUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)
	// Set Depth Query Parameter globally
	clientConfig.SetDepth(1)

	vmascConfig := vmasc.NewConfiguration(name, pwd, token, hostUrl)
	vmascConfig.UserAgent = appendUserAgent(vmascConfig.UserAgent)

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

	queryParams := map[string]string{
		"limit":    viper.GetString(constants.FlagLimit),
		"offset":   viper.GetString(constants.FlagOffset),
		"depth":    viper.GetString(constants.FlagDepth),
		"order-by": viper.GetString(constants.FlagOrderBy),
		"filter":   strings.Join(viper.GetStringSlice(constants.FlagFilters), ","),
	}

	s := jsontabwriter.GenerateVerboseOutput("queryParams: %v", queryParams)
	fmt.Println(s)

	fmt.Println("int depth:", viper.GetInt(constants.FlagDepth), "string depth:", viper.GetString(constants.FlagDepth))

	setQueryParams(sharedConfig, queryParams)
	setQueryParams(clientConfig, queryParams)
	setQueryParams(vmascConfig, queryParams)

	return &Client{
		URLOverride: hostUrl,

		// api.ionos.com
		AuthClient:     auth.NewAPIClient(configGuaranteeBasepath(sharedConfig, "/auth/v1")),
		CloudClient:    cloudv6.NewAPIClient(clientConfig),
		RegistryClient: containerregistry.NewAPIClient(configGuaranteeBasepath(sharedConfig, "/containerregistries")),

		PostgresClient: psql.NewAPIClient(configGuaranteeBasepath(sharedConfig, "/databases/postgresql")),
		MongoClient:    mongo.NewAPIClient(configGuaranteeBasepath(sharedConfig, "/databases/mongodb")),

		// regional APIs
		Apigateway:           apigateway.NewAPIClient(sharedConfig),
		CDNClient:            cdn.NewAPIClient(sharedConfig),
		CertManagerClient:    cert.NewAPIClient(sharedConfig),
		DnsClient:            dns.NewAPIClient(sharedConfig),
		Kafka:                kafka.NewAPIClient(sharedConfig),
		LoggingServiceClient: logging.NewAPIClient(sharedConfig),
		Monitoring:           monitoring.NewAPIClient(sharedConfig),
		VMAscClient:          vmasc.NewAPIClient(vmascConfig).AutoScalingGroupsApi,
		VPNClient:            vpn.NewAPIClient(sharedConfig),

		MariaClient:      mariadb.NewAPIClient(sharedConfig),
		InMemoryDBClient: inmemorydb.NewAPIClient(sharedConfig),
	}
}

type hasQueryParam interface {
	AddDefaultQueryParam(key, val string)
}

func setQueryParams(cfg hasQueryParam, params map[string]string) {
	for k, v := range params {
		cfg.AddDefaultQueryParam(k, v)
	}
}
