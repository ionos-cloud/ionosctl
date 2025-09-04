package client

import (
	"fmt"
	"net/url"

	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

	"github.com/spf13/viper"
)

// AuthSource represents a human-readable description of where the client's authentication credentials were sourced from.
type AuthSource string

const (
	AuthSourceEnvBearer AuthSource = "environment variable: IONOS_TOKEN"
	AuthSourceEnvBasic  AuthSource = "environment variables: IONOS_USERNAME, IONOS_PASSWORD"
	AuthSourceCfgBearer AuthSource = "credentials from config file: token"
	AuthSourceCfgBasic  AuthSource = "credentials from config file: username, password"
	AuthSourceNone      AuthSource = "no authentication provided"
)

// all possible sources in priority order
var AuthOrder = []AuthSource{
	AuthSourceEnvBearer,
	AuthSourceEnvBasic,
	AuthSourceCfgBearer,
	AuthSourceCfgBasic,
}

type Client struct {
	Config      *fileconfiguration.FileConfig
	ConfigPath  string // Path to the config file used to create this client, if any.
	AuthSource  AuthSource
	URLOverride string // If the client was created with a specific URL override, this will hold that value. If we notice a change in the URL, we need to re-create the client.

	Apigateway           *apigateway.APIClient
	CloudClient          *cloudv6.APIClient
	AuthClient           *auth.APIClient
	CertManagerClient    *cert.APIClient
	RegistryClient       *containerregistry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logging.APIClient
	VMAscClient          *vmasc.AutoScalingGroupsApiService
	VPNClient            *vpn.APIClient
	CDNClient            *cdn.APIClient
	Kafka                *kafka.APIClient
	Monitoring           *monitoring.APIClient

	PostgresClient   *psql.APIClient
	MongoClient      *mongo.APIClient
	MariaClient      *mariadb.APIClient
	InMemoryDBClient *inmemorydb.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

// hostWithoutPath strips any path from hostUrl; so that SDK clients append their own product paths,
// avoiding double basepaths ('/databases/postgresql/cloudapi/v6')
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
	copyCfg, err := auth.DeepCopy(cfg)
	if err != nil {
		// should never happen
		panic("failed to copy config: " + err.Error())
	}

	url := hostWithoutPath(copyCfg.Servers[0].URL)
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
