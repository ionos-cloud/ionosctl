/*
This is used for getting query parameters from options in the CLI.
And also for validate the parameters set - especially for filters.
*/
package query

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = "="

// ValidateFilters is currently being used in PreRun to check if the user provided any invalid filters
// availableFilters is set by the caller to be the properties of the struct (e.g. Name, Description, etc) with the first letter non-caps
// usageFilters is the usage string printed to output if the usage is wrong
//
// HACKY WORKAROUND WARNING: To work around the fact that we don't have access to availableFilters after this func's
// execution ends, and to keep support for 'any caps notation for properties is valid' rule e.g naME=myserver
// is a valid input, this func 'corrects' the caps of all keys that are valid. For example,
//
//	'--filters NAME=myserver,NAME=otherserver,imagetype=LINUX'
//
// would be corrected to
//
//	'--filters name=myserver,name=otherserver,imageType=LINUX'
//
// by MANUALLY setting the flag values, as well as the viper (global config) values.
func ValidateFilters(c *core.PreCommandConfig, availableFilters []string, usageFilters string) error {
	log.Printf("before: %s\n", viper.Get(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)))
	filtersKV, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Validating %v filters...", len(filtersKV))
	invalidFilters := make([]string, 0)
	correctedFilters := make(map[string][]string, 0)
	for key, vals := range filtersKV {
		validFilter, err := isValidFilter(key, availableFilters) // filterKey with proper caps
		correctedFilters[validFilter] = append(correctedFilters[validFilter], vals...)
		if err != nil {
			c.Printer.Verbose("Invalid Filter: %s", key)
			invalidFilters = append(invalidFilters, key)
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

	// Hacky workaround
	setString := ""
	for k, v := range correctedFilters {
		for _, vals := range v {
			setString += fmt.Sprintf("%s%s%s ", k, FiltersPartitionChar, vals)
		}
	}
	setString = fmt.Sprintf("%s", strings.TrimSuffix(setString, " "))
	viper.Set(core.GetFlagName(c.NS, cloudapiv6.ArgFilters), setString)
	_ = c.Command.Command.Flags().Set(cloudapiv6.ArgFilters, setString)
	// End hacky workaround

	log.Printf("after: %s\n", viper.Get(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)))

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

	if c.Command.Command.Flags().Changed(cloudapiv6.ArgOrderBy) {
		orderBy, _ := c.Command.Command.Flags().GetString(cloudapiv6.ArgOrderBy)
		listQueryParams = listQueryParams.SetOrderBy(orderBy)
	}

	if c.Command.Command.Flags().Changed(constants.FlagMaxResults) {
		maxResults, _ := c.Command.Command.Flags().GetInt32(constants.FlagMaxResults)
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
func getFilters(args []string, cmd *core.Command) (map[string][]string, error) {
	filtersKV := map[string][]string{}
	if len(args) == 0 {
		return filtersKV, errors.New("must provide at least one filter")
	}
	for _, arg := range args {
		if strings.Contains(arg, FiltersPartitionChar) {
			kv := strings.Split(arg, FiltersPartitionChar)
			filtersKV[kv[0]] = append(filtersKV[kv[0]], kv[1])
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
func isValidFilter(filter string, availableFiltersObjs ...[]string) (string, error) {
	for _, availableFilters := range availableFiltersObjs {
		for _, availableFilter := range availableFilters {
			if strings.ToLower(availableFilter) == strings.ToLower(filter) {
				return availableFilter, nil
			}
		}
	}
	return "", fmt.Errorf("TODO")
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
