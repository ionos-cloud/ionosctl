/*
	This is used for getting query parameters from options in the CLI.
	And also for validate the parameters set - especially for filters.
*/
package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = "="

// AvailableFilters for resources are usually split
// between Properties and Metadata Filters.
type AvailableFilters struct {
	PropertiesFilters []string
	MetadataFilters   []string
}

func ValidateFilters(c *core.PreCommandConfig, filters AvailableFilters) error {
	filtersKV, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Validating %v filters...", len(filtersKV))
	invalidFilters := make([]string, 0)
	for filterKey, _ := range filtersKV {
		if !isValidFilter(filterKey, filters.PropertiesFilters, filters.MetadataFilters) {
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
				getUsage(filters),
				c.Command.CommandPath(),
			),
		)
	}
	return nil
}

func GetListQueryParams(c *core.CommandConfig) (resources.ListQueryParams, error) {
	queryParams := resources.ListQueryParams{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		filters, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
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
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMaxResults)) {
		maxResults := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgMaxResults))
		queryParams = queryParams.SetMaxResults(maxResults)
	}

	return queryParams, nil
}

// getFilters should get the input from the user: --filters key=value,key=value
// and return a map with the corresponding key values
func getFilters(args []string, cmd *core.Command) (map[string]string, error) {
	filtersKV := map[string]string{}
	if len(args) == 0 {
		return filtersKV, errors.New("len of args must be different than 0")
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

func getUsage(filters AvailableFilters) string {
	usage := "Available Filters:\n"
	if len(filters.PropertiesFilters) > 0 {
		usage = fmt.Sprintf("%s* filter by property: %s", usage, filters.PropertiesFilters)
	}
	if len(filters.MetadataFilters) > 0 {
		usage = fmt.Sprintf("%s\n* filter by metadata: %s", usage, filters.MetadataFilters)
	}
	return usage
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
