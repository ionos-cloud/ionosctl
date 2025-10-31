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
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"golang.org/x/exp/slices"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = "="

// ValidateFilters checks for invalid user-provided filters in PreRun. It compares user input against a list of valid filters (`availableFilters`),
// which are derived from struct properties (like 'Name', 'Description') with the first letter in lowercase, following the SDK struct naming conventions.
//
// Usage:
//   - `availableFilters`: List of valid filter keys (struct properties with lowercase first letter).
//   - `usageFilters`: Usage information string displayed if an invalid filter is provided.
//
// IMPORTANT:
//   - This function relies on the naming convention of SDK struct fields.
//   - If these names change, a refactor of filters.go in commands/cloudapi/completer may be necessary.
//
// WORKAROUND:
//
//   - This function 'corrects' the capitalization of valid keys. For instance, it would adjust
//
//     '--filters NAME=myserver,NAME=otherserver,imagetype=LINUX' to
//
//     '--filters name=myserver,name=otherserver,imageType=LINUX'
//
//   - This adjustment is done by manually setting both the flag and global config (viper) values.
//
// The function first retrieves and processes filter key-value pairs. If any invalid filters are found, it reports them and returns an error.
// Otherwise, it applies the capitalization correction and updates the flag and viper settings accordingly.
func ValidateFilters(c *core.PreCommandConfig, availableFilters []string, usageFilters string) error {
	filters, err := c.Command.Command.Flags().GetStringSlice(cloudapiv6.ArgFilters)
	if err != nil {
		return fmt.Errorf("failed getting stringSlice: %w", err)
	}
	filtersKV, err := getFilters(filters, c.Command)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Validating %v filters...", len(filtersKV)))

	correctedFilters := make(map[string][]string, 0)
	for key, vals := range filtersKV {
		keyWithProperCaps, errValidFilter := getFirstValidFilter(key, availableFilters...)
		if errValidFilter != nil {
			err = errors.Join(err, errValidFilter)
			continue
		}

		correctedFilters[keyWithProperCaps] = append(correctedFilters[keyWithProperCaps], vals...)
	}

	if err != nil {
		return fmt.Errorf("encountered invalid filters:\n%w\n\n%s", err, usageFilters)
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

	return nil
}

func GetListQueryParams(c *core.CommandConfig) (resources.ListQueryParams, error) {
	listQueryParams := resources.ListQueryParams{}

	if c.Command.Command.Flags().Changed(cloudapiv6.ArgFilters) {
		filters, err := c.Command.Command.Flags().GetStringSlice(cloudapiv6.ArgFilters)
		if err != nil {
			return listQueryParams, err
		}
		mapFilters, err := getFilters(filters, c.Command)
		if err != nil {
			return listQueryParams, err
		}

		if len(filters) > 0 {
			listQueryParams = listQueryParams.SetFilters(mapFilters)
		}
	}

	if c.Command.Command.Flags().Changed(cloudapiv6.ArgOrderBy) {
		orderBy, _ := c.Command.Command.Flags().GetString(cloudapiv6.ArgOrderBy)
		listQueryParams = listQueryParams.SetOrderBy(orderBy)
	}

	// No guard against "changed", as we want the pflag imposed defaults
	depth, _ := c.Command.Command.Flags().GetInt32(cloudapiv6.ArgDepth)
	listQueryParams = listQueryParams.SetDepth(depth)

	if !structs.IsZero(listQueryParams) || !structs.IsZero(listQueryParams.QueryParams) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Query Parameters set: %v, %v",
			utils.GetPropertiesKVSet(listQueryParams), utils.GetPropertiesKVSet(listQueryParams.QueryParams)))
	}

	return listQueryParams, nil
}

// getFilters should get the input from the user: --filters key=value,key=value
// and return a map with the corresponding key values
func getFilters(args []string, cmd *core.Command) (map[string][]string, error) {
	filtersKV := map[string][]string{}

	for _, arg := range args {
		if arg == "" {
			// Workaround for interactive shell:
			// further usages after resetting to an empty slice results in an empty string as the first value
			continue
		}

		if strings.Contains(arg, FiltersPartitionChar) {
			kv := strings.Split(arg, FiltersPartitionChar)
			filtersKV[kv[0]] = append(filtersKV[kv[0]], kv[1])
		} else {
			return filtersKV, errors.New(
				fmt.Sprintf("\"%s --filters\" option set incorrectly. \n\nUsage: %s --filters KEY1%sVALUE1,KEY2%sVALUE2\n\nFor more details, see '%s --help'.",
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

// getFirstValidFilter will return the filter with the expected case
// if it is part of availableFilters - the comparison is case insensitive.
// An error is returned if the filter is not found in the slice.
// Examples:
//
//	getFirstValidFilter("imagetype", "imageType", "OtherFilterHere", "foo") -> "imageType"
//	getFirstValidFilter("imagetype", "OtherFilterHere", "foo", "bar") ->
//	  -> error("imagetype is not case insensitively equal to any of OtherFilterHere, foo, bar")
func getFirstValidFilter(filter string, availableFilters ...string) (string, error) {
	idx := slices.IndexFunc(availableFilters, func(s string) bool {
		return strings.EqualFold(filter, s)
	})

	if idx == -1 {
		return "", fmt.Errorf("%s is not a valid filter", filter)
	}

	return availableFilters[idx], nil
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
