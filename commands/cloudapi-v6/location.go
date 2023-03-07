package commands

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLocationCols, printer.ColsMessage(allLocationCols))
	_ = viper.BindPFlag(core.GetFlagName(locationCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddStringFlag(cloudapiv6.ArgLocationId, cloudapiv6.ArgIdShort, "", cloudapiv6.LocationId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	locationCmd.AddCommand(CpuCmd())

	return locationCmd
}

func PreRunLocationsList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.LocationsFilters(), completer.LocationsFiltersUsage())
	}
	return nil
}

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgLocationId)
}

func RunLocationList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	locations, resp, err := c.CloudApiV6Services.Locations().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: locations,
		KeyValue:   getLocationsKVMaps(getLocations(locations)),
		Columns:    printer.GetHeaders(allLocationCols, defaultLocationCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))),
	})
}

func RunLocationGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	locId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId))
	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}
	c.Printer.Verbose("Location with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId)))
	loc, resp, err := c.CloudApiV6Services.Locations().GetByRegionAndLocationId(ids[0], ids[1], queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: loc,
		KeyValue:   getLocationsKVMaps(getLocation(loc)),
		Columns:    printer.GetHeaders(allLocationCols, defaultLocationCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))),
	})
}

// Output Printing

var (
	defaultLocationCols = []string{"LocationId", "Name", "CpuFamily"}
	allLocationCols     = []string{"LocationId", "Name", "Features", "ImageAliases", "CpuFamily"}
)

type LocationPrint struct {
	LocationId   string   `json:"LocationId,omitempty"`
	Name         string   `json:"Name,omitempty"`
	Features     []string `json:"Features,omitempty"`
	CpuFamily    []string `json:"CpuFamily,omitempty"`
	ImageAliases []string `json:"ImageAliases,omitempty"`
}

func getLocation(u *resources.Location) []resources.Location {
	locs := make([]resources.Location, 0)
	if u != nil {
		locs = append(locs, resources.Location{Location: u.Location})
	}
	return locs
}

func getLocations(locations resources.Locations) []resources.Location {
	locationObjs := make([]resources.Location, 0)
	if items, ok := locations.GetItemsOk(); ok && items != nil {
		for _, location := range *items {
			locationObjs = append(locationObjs, resources.Location{Location: location})
		}
	}
	return locationObjs
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
		if cpus, ok := properties.GetCpuArchitectureOk(); ok && cpus != nil {
			cpuFamilies := make([]string, 0)
			for _, cpu := range *cpus {
				if cpuFamily, ok := cpu.GetCpuFamilyOk(); ok && cpuFamily != nil {
					cpuFamilies = append(cpuFamilies, *cpuFamily)
				}
			}
			dcPrint.CpuFamily = cpuFamilies
		}
		o := structs.Map(dcPrint)
		out = append(out, o)
	}
	return out
}
