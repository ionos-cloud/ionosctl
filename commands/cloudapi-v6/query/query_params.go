package query

import (
	"errors"
	"strings"

	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = ":"

func GetListQueryParams(c *core.CommandConfig) (resources.ListQueryParams, error) {
	queryParams := resources.ListQueryParams{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		filters, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)))
		if err != nil {
			return queryParams, err
		}
		if len(filters) > 0 {
			queryParams = queryParams.SetFilters(filters)
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgOrderBy)) {
		orderBy := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgOrderBy))
		queryParams = queryParams.SetOrderBy(orderBy)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgOffset)) {
		offset := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgOffset))
		queryParams = queryParams.SetOffset(offset)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLimit)) {
		limit := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLimit))
		queryParams = queryParams.SetLimit(limit)
	}
	defaultQueryParams := resources.QueryParams{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDepth)) {
		depth := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgDepth))
		defaultQueryParams = defaultQueryParams.SetDepth(depth)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPretty)) {
		pretty := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPretty))
		defaultQueryParams = defaultQueryParams.SetPretty(pretty)
	}
	queryParams = queryParams.SetDefaultQueryParams(defaultQueryParams)

	return queryParams, nil
}

// getFilters should get the input from the user: --filters key:value,key:value
// and return a map with the corresponding key values
func getFilters(args []string) (map[string]string, error) {
	filtersKV := map[string]string{}
	if len(args) == 0 {
		return filtersKV, errors.New("len of args must be different than 0")
	}
	for _, arg := range args {
		if strings.Contains(arg, FiltersPartitionChar) {
			kv := strings.Split(arg, FiltersPartitionChar)
			filtersKV[kv[0]] = kv[1]
		} else {
			return filtersKV, errors.New("--filters should have the following format: --filters KEY1:VALUE1, KEY2:VALUE2")
		}
	}
	return filtersKV, nil
}
