/*
This is used for getting query parameters from options in the CLI.
And also for validate the parameters set - especially for filters.
*/
package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/utils"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = "="

func ValidateFilters(c *core.PreCommandConfig, availableFilters []string, usageFilters string) error {
	filtersKV, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Validating %v filters...", len(filtersKV))
	invalidFilters := make([]string, 0)
	for filterKey, _ := range filtersKV {
		if !isValidFilter(filterKey, availableFilters) {
			c.Printer.Verbose("Invalid Filter: %s", filterKey)
			invalidFilters = append(invalidFilters, filterKey)
		} else {
			c.Printer.Verbose("Valid Filter: %s", filterKey)
		}
	}
	if len(invalidFilters) > 0 {
		return errors.New(
			fmt.Sprintf("%q has at least %d invalid %s.\n\n%s\n\nFor more details, see '%s --help'.",
				c.Command.CommandPath(),
				len(invalidFilters),
				pluralize("filter", len(invalidFilters)),
				usageFilters,
				c.Command.CommandPath(),
			),
		)
	}
	return nil
}

func GetListQueryParams(c *core.CommandConfig) (resources.ListQueryParams, error) {
	listQueryParams := resources.ListQueryParams{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		filters, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
		if err != nil {
			return listQueryParams, err
		}
		if len(filters) > 0 {
			listQueryParams = listQueryParams.SetFilters(filters)
		}
	}

	if c.Command.Command.Flags().Changed(cloudapiv6.ArgMaxResults) {
		orderBy, _ := c.Command.Command.Flags().GetString(cloudapiv6.ArgOrderBy)
		listQueryParams = listQueryParams.SetOrderBy(orderBy)
	}

	if c.Command.Command.Flags().Changed(cloudapiv6.ArgMaxResults) {
		maxResults, _ := c.Command.Command.Flags().GetInt32(cloudapiv6.ArgMaxResults)
		listQueryParams = listQueryParams.SetMaxResults(maxResults)
	}

	// No guard against "changed", as we want the pflag imposed defaults
	depth, _ := c.Command.Command.Flags().GetInt32(cloudapiv6.ArgDepth)
	listQueryParams = listQueryParams.SetDepth(depth)

	if !structs.IsZero(listQueryParams) || !structs.IsZero(listQueryParams.QueryParams) {
		c.Printer.Verbose("Query Parameters set: %v, %v", utils.GetPropertiesKVSet(listQueryParams), utils.GetPropertiesKVSet(listQueryParams.QueryParams))
	}

	return listQueryParams, nil
}

// getFilters should get the input from the user: --filters key=value,key=value
// and return a map with the corresponding key values
func getFilters(args []string, cmd *core.Command) (map[string]string, error) {
	filtersKV := map[string]string{}
	if len(args) == 0 {
		return filtersKV, errors.New("must provide at least one filter")
	}
	for _, arg := range args {
		if strings.Contains(arg, FiltersPartitionChar) {
			kv := strings.Split(arg, FiltersPartitionChar)
			filtersKV[kv[0]] = kv[1]
		} else {
			return filtersKV, errors.New(
				fmt.Sprintf("\"%s --filters\" option set incorrectly.\n\nUsage: %s --filters KEY1%sVALUE1,KEY2%sVALUE2\n\nFor more details, see '%s --help'.",
					cmd.CommandPath(),
					cmd.CommandPath(),
					FiltersPartitionChar,
					FiltersPartitionChar,
					cmd.CommandPath(),
				),
			)
		}
	}
	return filtersKV, nil
}

// isValidFilter will return true if the filter is part
// of the available filters array and false if is not.
func isValidFilter(filter string, availableFiltersObjs ...[]string) bool {
	for _, availableFilters := range availableFiltersObjs {
		for _, availableFilter := range availableFilters {
			if availableFilter == filter {
				return true
			}
		}
	}
	return false
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
