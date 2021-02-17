package commands

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func location() *builder.Command {
	locationCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "location",
			Short:            "Location Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl location` + "`" + ` allows you to see information about locations available to create objects.`,
			TraverseChildren: true,
		},
	}
	globalFlags := locationCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultLocationCols, "Columns to be printed in the standard output")
	viper.BindPFlag(builder.GetGlobalFlagName(locationCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), locationCmd, noPreRun, RunLocationList, "list", "List Locations",
		"Use this command to get a list of available locations to create objects on.",
		listLocationExample, true)

	return locationCmd
}

func RunLocationList(c *builder.CommandConfig) error {
	locations, _, err := c.Locations().List()
	if err != nil {
		return err
	}
	c.Printer.Print(utils.Result{
		OutputJSON: locations,
		KeyValue:   getLocationsKVMaps(getLocations(locations)),
		Columns:    getLocationCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
	return nil
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
			utils.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
}

func getLocations(datacenters resources.Locations) []resources.Location {
	dc := make([]resources.Location, 0)
	for _, d := range *datacenters.Items {
		dc = append(dc, resources.Location{d})
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
	err := config.LoadFile()
	utils.CheckError(err, outErr)

	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.ArgServerUrl),
	)
	utils.CheckError(err, outErr)

	locationSvc := resources.NewLocationService(clientSvc.Get(), context.TODO())
	locations, _, err := locationSvc.List()
	utils.CheckError(err, outErr)

	lcIds := make([]string, 0)
	if locations.Locations.Items != nil {
		for _, d := range *locations.Locations.Items {
			lcIds = append(lcIds, *d.GetId())
		}
	} else {
		return nil
	}
	return lcIds
}
