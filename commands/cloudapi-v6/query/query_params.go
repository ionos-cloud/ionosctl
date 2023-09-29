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
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"golang.org/x/exp/slices"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

const FiltersPartitionChar = "="

// ValidateFilters is currently being used in PreRun to check if the user provided any invalid filters
// availableFilters is set by the caller to be the properties of the struct (e.g. Name, Description, etc) with the first letter non-caps, by using the SDK generated structs.
// usageFilters is the usage string printed to output if the usage is wrong
//
// WARNING: It just so happens that we can find the `availableFilters` by using the SDK structs namings with non-caps first letter,
// but if the SDK struct fields naming were to change, then we would be forced to refactor the whole file commands/cloudapi/completer/filters.go
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
	filtersKV, err := getFilters(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)), c.Command)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Validating %v filters...", len(filtersKV)))

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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Query Parameters set: %v, %v",
			utils.GetPropertiesKVSet(listQueryParams), utils.GetPropertiesKVSet(listQueryParams.QueryParams)))
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
