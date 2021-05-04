package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func datacenter() *builder.Command {
	ctx := context.TODO()
	datacenterCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"dc"},
			Short:            "Data Center Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl datacenter` + "`" + ` allow you to create, list, get, update and delete Data Centers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(datacenterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, datacenterCmd, noPreRun, RunDataCenterList, "list", "List Data Centers",
		"Use this command to retrieve a complete list of Virtual Data Centers provisioned under your account.", listDatacenterExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, datacenterCmd, PreRunDataCenterId, RunDataCenterGet, "get", "Get a Data Center",
		"Use this command to retrieve details about a Virtual Data Center by using its ID.\n\nRequired values to run command:\n\n* Data Center Id", getDatacenterExample, true)
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, datacenterCmd, noPreRun, RunDataCenterCreate, "create", "Create a Data Center",
		`Use this command to create a Virtual Data Center. You can specify the name, description or location for the object.

Virtual Data Centers (VDCs) are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.`, createDatacenterExample, true)
	create.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	create.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")
	create.AddStringFlag(config.ArgDataCenterRegion, "", "de/txl", "Location for the Data Center")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterRegion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, datacenterCmd, PreRunDataCenterId, RunDataCenterUpdate, "update", "Update a Data Center",
		`Use this command to change a Virtual Data Center's name, description.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, updateDatacenterExample, true)
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	update.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, datacenterCmd, PreRunDataCenterId, RunDataCenterDelete, "delete", "Delete a Data Center",
		`Use this command to delete a specified Virtual Data Center (VDC) from your account. This will remove all objects within the VDC and remove the VDC object itself. 

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id`, deleteDatacenterExample, true)
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be deleted [seconds]")

	datacenterCmd.AddCommand(datacenterLabel())

	return datacenterCmd
}

func PreRunDataCenterId(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId)
}

func RunDataCenterList(c *builder.CommandConfig) error {
	datacenters, _, err := c.DataCenters().List()
	if err != nil {
		return err
	}
	dcs := getDataCenters(datacenters)
	return c.Printer.Print(printer.Result{
		OutputJSON: datacenters,
		KeyValue:   getDataCentersKVMaps(dcs),
		Columns:    getDataCenterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunDataCenterGet(c *builder.CommandConfig) error {
	datacenter, _, err := c.DataCenters().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		KeyValue:   getDataCentersKVMaps([]resources.Datacenter{*datacenter}),
		Columns:    getDataCenterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		OutputJSON: datacenter,
	})
}

func RunDataCenterCreate(c *builder.CommandConfig) error {
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterName))
	description := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterDescription))
	region := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterRegion))
	dc, resp, err := c.DataCenters().Create(name, description, region)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		KeyValue:    getDataCentersKVMaps([]resources.Datacenter{*dc}),
		Columns:     getDataCenterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		OutputJSON:  dc,
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunDataCenterUpdate(c *builder.CommandConfig) error {
	input := resources.DatacenterProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterDescription)) {
		input.SetDescription(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterDescription)))
	}
	dc, resp, err := c.DataCenters().Update(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		KeyValue:    getDataCentersKVMaps([]resources.Datacenter{*dc}),
		Columns:     getDataCenterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		OutputJSON:  dc,
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "update",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunDataCenterDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete data center"); err != nil {
		return err
	}
	resp, err := c.DataCenters().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)))
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "datacenter",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultDatacenterCols = []string{"DatacenterId", "Name", "Location"}

type DatacenterPrint struct {
	DatacenterId string `json:"DatacenterId,omitempty"`
	Name         string `json:"Name,omitempty"`
	Location     string `json:"Location,omitempty"`
	Description  string `json:"Description,omitempty"`
	Version      int32  `json:"Version,omitempty"`
}

func getDataCenterCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultDatacenterCols
	}

	columnsMap := map[string]string{
		"DatacenterId": "DatacenterId",
		"Name":         "Name",
		"Location":     "Location",
		"Version":      "Version",
		"Description":  "Description",
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

func getDataCenters(datacenters resources.Datacenters) []resources.Datacenter {
	dc := make([]resources.Datacenter, 0)
	for _, d := range *datacenters.Items {
		dc = append(dc, resources.Datacenter{Datacenter: d})
	}
	return dc
}

func getDataCentersKVMaps(dcs []resources.Datacenter) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		properties := dc.GetProperties()
		var dcPrint DatacenterPrint
		if dcid, ok := dc.GetIdOk(); ok && dcid != nil {
			dcPrint.DatacenterId = *dcid
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			dcPrint.Name = *name
		}
		if location, ok := properties.GetLocationOk(); ok && location != nil {
			dcPrint.Location = *location
		}
		if description, ok := properties.GetDescriptionOk(); ok && description != nil {
			dcPrint.Description = *description
		}
		if version, ok := properties.GetVersionOk(); ok && version != nil {
			dcPrint.Version = *version
		}
		o := structs.Map(dcPrint)
		out = append(out, o)
	}
	return out
}

func getDataCentersIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	datacenterSvc := resources.NewDataCenterService(clientSvc.Get(), context.TODO())
	datacenters, _, err := datacenterSvc.List()
	clierror.CheckError(err, outErr)
	dcIds := make([]string, 0)
	if items, ok := datacenters.Datacenters.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				dcIds = append(dcIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return dcIds
}
