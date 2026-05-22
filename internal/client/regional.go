package client

import (
	"net/http"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
)

// NewRegionalConfig creates a shared.Configuration for a specific regional URL,
// reusing credentials from the existing client singleton. It applies the same
// user agent, query params, and filters as the main client builder.
// This is used by ListAllLocations to create standalone SDK clients per location,
// avoiding mutation of the global singleton.
func NewRegionalConfig(url string) *shared.Configuration {
	cl := Must()
	cloudCfg := cl.CloudClient.GetConfig()

	cfg := shared.NewConfiguration(cloudCfg.Username, cloudCfg.Password, cloudCfg.Token, url)
	cfg.UserAgent = appendUserAgent(cfg.UserAgent)
	cfg.HTTPClient = &http.Client{}

	// Apply log level (shared.SdkLogLevel is package-level, already set by main client init)

	// Apply query params (mirrors builder.go logic, including deprecated --max-results)
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
	for k, v := range queryParams {
		if v == "" {
			delete(queryParams, k)
		}
	}
	setQueryParams(cfg, queryParams)

	// Apply filters
	setFilters(cfg, viper.GetStringSlice(constants.FlagFilters))

	return cfg
}
