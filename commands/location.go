package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func location() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultLocationCols, utils.ColsMessage(allLocationCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(locationCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLocationCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, locationCmd, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "location",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Locations",
		LongDesc:   "Use this command to get a list of available locations to create objects on.",
		Example:    listLocationExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunLocationList,
		InitClient: true,
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
	get.AddStringFlag(config.ArgLocationId, config.ArgIdShort, "", config.LocationId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return locationCmd
}

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgLocationId)
}

func RunLocationList(c *core.CommandConfig) error {
	locations, _, err := c.Locations().List()
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
	locId := viper.GetString(core.GetFlagName(c.NS, config.ArgLocationId))
	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}
	c.Printer.Verbose("Location with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgLocationId)))
	loc, _, err := c.Locations().GetByRegionAndLocationId(ids[0], ids[1])
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

func getLocation(u *v5.Location) []v5.Location {
	locs := make([]v5.Location, 0)
	if u != nil {
		locs = append(locs, v5.Location{Location: u.Location})
	}
	return locs
}

func getLocations(locations v5.Locations) []v5.Location {
	dc := make([]v5.Location, 0)
	for _, d := range *locations.Items {
		dc = append(dc, v5.Location{Location: d})
	}
	return dc
}

func getLocationsKVMaps(dcs []v5.Location) []map[string]interface{} {
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

func getLocationIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	locationSvc := v5.NewLocationService(clientSvc.Get(), context.TODO())
	locations, _, err := locationSvc.List()
	clierror.CheckError(err, outErr)
	lcIds := make([]string, 0)
	if items, ok := locations.Locations.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				lcIds = append(lcIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return lcIds
}
