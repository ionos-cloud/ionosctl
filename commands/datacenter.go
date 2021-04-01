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
	datacenterCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"dc"},
			Short:            "Data Center Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl datacenter` + "`" + ` allow you to create, list, get, update and delete Data Centers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	viper.BindPFlag(builder.GetGlobalFlagName(datacenterCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), datacenterCmd, noPreRun, RunDataCenterList, "list", "List Data Centers",
		"Use this command to list all Data Centers on your account.", listDatacenterExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), datacenterCmd, PreRunDataCenterIdValidate, RunDataCenterGet, "get", "Get a Data Center",
		"Use this command to get information about a specified Data Center.\n\nRequired values to run command:\n\n* Data Center Id", getDatacenterExample, true)
	get.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id [Required flag]")
	get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), datacenterCmd, noPreRun, RunDataCenterCreate, "create", "Create a Data Center",
		`Use this command to create a Data Center. You can specify the name, description or location for the object.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.`, createDatacenterExample, true)
	create.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	create.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")
	create.AddStringFlag(config.ArgDataCenterRegion, "", "de/txl", "Location for the Data Center")
	create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterRegion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(context.TODO(), datacenterCmd, PreRunDataCenterIdValidate, RunDataCenterUpdate, "update", "Update a Data Center",
		`Use this command to change a Data Center's name, description.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, updateDatacenterExample, true)
	update.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id [Required flag]")
	update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgDataCenterName, "", "", "Name of the Data Center")
	update.AddStringFlag(config.ArgDataCenterDescription, "", "", "Description of the Data Center")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), datacenterCmd, PreRunDataCenterIdValidate, RunDataCenterDelete, "delete", "Delete a Data Center",
		`Use this command to delete a specified Data Center from your account. This is irreversible.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id`, deleteDatacenterExample, true)
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", "The unique Data Center Id [Required flag]")
	deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Data Center to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Data Center to be deleted [seconds]")

	return datacenterCmd
}

func PreRunDataCenterIdValidate(c *builder.PreCommandConfig) error {
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

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
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

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
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
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete data center")
	if err != nil {
		return err
	}
	resp, err := c.DataCenters().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)))
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
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
		dc = append(dc, resources.Datacenter{d})
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
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)

	datacenterSvc := resources.NewDataCenterService(clientSvc.Get(), context.TODO())
	datacenters, _, err := datacenterSvc.List()
	clierror.CheckError(err, outErr)

	dcIds := make([]string, 0)
	if datacenters.Datacenters.Items != nil {
		for _, d := range *datacenters.Datacenters.Items {
			dcIds = append(dcIds, *d.GetId())
		}
	} else {
		return nil
	}
	return dcIds
}
