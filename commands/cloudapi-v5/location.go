package cloudapi_v5

import (
	"context"
	"errors"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LocationCmd() *core.Command {
	ctx := context.TODO()
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc"},
			Short:            "Location Operations",
			Long:             "The sub-command of `ionosctl location` allows you to see information about locations available to create objects.",
			TraverseChildren: true,
		},
	}
	globalFlags := locationCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultLocationCols, printer.ColsMessage(allLocationCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(locationCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLocationCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, locationCmd, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "location",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Locations",
		LongDesc:   "Use this command to get a list of available locations to create objects on.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LocationsFiltersUsage(),
		Example:    listLocationExample,
		PreCmdRun:  PreRunLocationsList,
		CmdRun:     RunLocationList,
		InitClient: true,
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationsFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, locationCmd, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "location",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Location",
		LongDesc:   "Use this command to get information about a specific Location from a Region.\n\nRequired values to run command:\n\n* Location Id",
		Example:    getLocationExample,
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv5.ArgLocationId, cloudapiv5.ArgIdShort, "", cloudapiv5.LocationId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return locationCmd
}

func PreRunLocationsList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.LocationsFilters(), completer.LocationsFiltersUsage())
	}
	return nil
}

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgLocationId)
}

func RunLocationList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	locations, resp, err := c.CloudApiV5Services.Locations().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: locations,
		KeyValue:   getLocationsKVMaps(getLocations(locations)),
		Columns:    getLocationCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunLocationGet(c *core.CommandConfig) error {
	locId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLocationId))
	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}
	c.Printer.Verbose("Location with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLocationId)))
	loc, resp, err := c.CloudApiV5Services.Locations().GetByRegionAndLocationId(ids[0], ids[1])
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: loc,
		KeyValue:   getLocationsKVMaps(getLocation(loc)),
		Columns:    getLocationCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
	})
}

// Output Printing

var (
	defaultLocationCols = []string{"LocationId", "Name", "Features"}
	allLocationCols     = []string{"LocationId", "Name", "Features", "ImageAliases"}
)

type LocationPrint struct {
	LocationId   string   `json:"LocationId,omitempty"`
	Name         string   `json:"Name,omitempty"`
	Features     []string `json:"Features,omitempty"`
	ImageAliases []string `json:"ImageAliases,omitempty"`
}

func getLocationCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultLocationCols
	}

	columnsMap := map[string]string{
		"LocationId":   "LocationId",
		"Name":         "Name",
		"Features":     "Features",
		"ImageAliases": "ImageAliases",
	}
	var locationsCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			locationsCols = append(locationsCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return locationsCols
}

func getLocation(u *resources.Location) []resources.Location {
	locs := make([]resources.Location, 0)
	if u != nil {
		locs = append(locs, resources.Location{Location: u.Location})
	}
	return locs
}

func getLocations(locations resources.Locations) []resources.Location {
	dc := make([]resources.Location, 0)
	for _, d := range *locations.Items {
		dc = append(dc, resources.Location{Location: d})
	}
	return dc
}

func getLocationsKVMaps(dcs []resources.Location) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		properties := dc.GetProperties()
		var dcPrint LocationPrint
		if dcid, ok := dc.GetIdOk(); ok && dcid != nil {
			dcPrint.LocationId = *dcid
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			dcPrint.Name = *name
		}
		if features, ok := properties.GetFeaturesOk(); ok && features != nil {
			dcPrint.Features = *features
		}
		if aliases, ok := properties.GetImageAliasesOk(); ok && aliases != nil {
			dcPrint.ImageAliases = *aliases
		}
		o := structs.Map(dcPrint)
		out = append(out, o)
	}
	return out
}
