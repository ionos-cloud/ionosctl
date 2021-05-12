package commands

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func location() *core.Command {
	ctx := context.TODO()
	locationCmd := &core.Command{
		NS: "location",
		Command: &cobra.Command{
			Use:              "location",
			Short:            "Location Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl location` + "`" + ` allows you to see information about locations available to create objects.`,
			TraverseChildren: true,
		},
	}
	globalFlags := locationCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultLocationCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(locationCmd.NS, config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	core.NewCommand(ctx, locationCmd, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "location",
		Verb:       "list",
		ShortDesc:  "List Locations",
		LongDesc:   "Use this command to get a list of available locations to create objects on.",
		Example:    listLocationExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunLocationList,
		InitClient: true,
	})

	return locationCmd
}

func RunLocationList(c *core.CommandConfig) error {
	locations, _, err := c.Locations().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: locations,
		KeyValue:   getLocationsKVMaps(getLocations(locations)),
		Columns:    getLocationCols(core.GetGlobalFlagName(c.Namespace, config.ArgCols), c.Printer.GetStderr()),
	})
}

var defaultLocationCols = []string{"LocationId", "Name", "Features"}

type LocationPrint struct {
	LocationId string   `json:"LocationId,omitempty"`
	Name       string   `json:"Name,omitempty"`
	Features   []string `json:"Features,omitempty"`
}

func getLocationCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultLocationCols
	}

	columnsMap := map[string]string{
		"LocationId": "LocationId",
		"Name":       "Name",
		"Features":   "Features",
	}
	var datacenterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			datacenterCols = append(datacenterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
}

func getLocations(datacenters resources.Locations) []resources.Location {
	dc := make([]resources.Location, 0)
	for _, d := range *datacenters.Items {
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
		o := structs.Map(dcPrint)
		out = append(out, o)
	}
	return out
}

func getLocationIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	locationSvc := resources.NewLocationService(clientSvc.Get(), context.TODO())
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
