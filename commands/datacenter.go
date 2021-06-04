package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func datacenter() *core.Command {
	ctx := context.TODO()
	datacenterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"d", "dc"},
			Short:            "Data Center Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl datacenter` + "`" + ` allow you to create, list, get, update and delete Data Centers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultDatacenterCols,
		fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", allDatacenterCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(datacenterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = datacenterCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allDatacenterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace:  "datacenter",
		Resource:   "datacenter",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Data Centers",
		LongDesc:   "Use this command to retrieve a complete list of Virtual Data Centers provisioned under your account.",
		Example:    listDatacenterExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunDataCenterList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace:  "datacenter",
		Resource:   "datacenter",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Data Center",
		LongDesc:   "Use this command to retrieve details about a Virtual Data Center by using its ID.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    getDatacenterExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgDataCenterId, config.ArgIdShort, "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Data Center",
		LongDesc: `Use this command to create a Virtual Data Center. You can specify the name, description or location for the object.

Virtual Data Centers (VDCs) are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.`,
		Example:    createDatacenterExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunDataCenterCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Data Center")
	create.AddStringFlag(config.ArgDescription, config.ArgDescriptionShort, "", "Description of the Data Center")
	create.AddStringFlag(config.ArgLocation, config.ArgLocationShort, "de/txl", "Location for the Data Center")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Data Center creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Data Center creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Data Center",
		LongDesc: `Use this command to change a Virtual Data Center's name, description.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    updateDatacenterExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgDataCenterId, config.ArgIdShort, "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Data Center")
	update.AddStringFlag(config.ArgDescription, config.ArgDescriptionShort, "", "Description of the Data Center")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Data Center update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Data Center update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Data Center",
		LongDesc: `Use this command to delete a specified Virtual Data Center (VDC) from your account. This will remove all objects within the VDC and remove the VDC object itself. 

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id`,
		Example:    deleteDatacenterExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, config.ArgIdShort, "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Data Center deletion")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Data Center deletion [seconds]")

	return datacenterCmd
}

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId)
}

func RunDataCenterList(c *core.CommandConfig) error {
	datacenters, _, err := c.DataCenters().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(nil, c, getDataCenters(datacenters)))
}

func RunDataCenterGet(c *core.CommandConfig) error {
	dc, _, err := c.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(nil, c, []resources.Datacenter{*dc}))
}

func RunDataCenterCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, config.ArgDescription))
	region := viper.GetString(core.GetFlagName(c.NS, config.ArgLocation))
	dc, resp, err := c.DataCenters().Create(name, description, region)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(resp, c, []resources.Datacenter{*dc}))
}

func RunDataCenterUpdate(c *core.CommandConfig) error {
	input := resources.DatacenterProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDescription)) {
		input.SetDescription(viper.GetString(core.GetFlagName(c.NS, config.ArgDescription)))
	}
	dc, resp, err := c.DataCenters().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(resp, c, []resources.Datacenter{*dc}))
}

func RunDataCenterDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete data center"); err != nil {
		return err
	}
	resp, err := c.DataCenters().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(resp, c, nil))
}

// Output Printing

var (
	defaultDatacenterCols = []string{"DatacenterId", "Name", "Location", "CpuFamily", "State"}
	allDatacenterCols     = []string{"DatacenterId", "Name", "Location", "State", "Description", "Version", "Features", "CpuFamily", "SecAuthProtection"}
)

type DatacenterPrint struct {
	DatacenterId      string   `json:"DatacenterId,omitempty"`
	Name              string   `json:"Name,omitempty"`
	Location          string   `json:"Location,omitempty"`
	Description       string   `json:"Description,omitempty"`
	Version           int32    `json:"Version,omitempty"`
	State             string   `json:"State,omitempty"`
	Features          []string `json:"Features,omitempty"`
	CpuFamily         []string `json:"CpuFamily,omitempty"`
	SecAuthProtection bool     `json:"SecAuthProtection,omitempty"`
}

func getDataCenterPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.Datacenter) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getDataCentersKVMaps(dcs)
			r.Columns = getDataCenterCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getDataCenterCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultDatacenterCols
	}

	columnsMap := map[string]string{
		"DatacenterId":      "DatacenterId",
		"Name":              "Name",
		"Location":          "Location",
		"Version":           "Version",
		"Description":       "Description",
		"State":             "State",
		"Features":          "Features",
		"CpuFamily":         "CpuFamily",
		"SecAuthProtection": "SecAuthProtection",
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
		var dcPrint DatacenterPrint
		if dcid, ok := dc.GetIdOk(); ok && dcid != nil {
			dcPrint.DatacenterId = *dcid
		}
		if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				dcPrint.Name = *name
			}
			if loc, ok := properties.GetLocationOk(); ok && loc != nil {
				dcPrint.Location = *loc
			}
			if description, ok := properties.GetDescriptionOk(); ok && description != nil {
				dcPrint.Description = *description
			}
			if ver, ok := properties.GetVersionOk(); ok && ver != nil {
				dcPrint.Version = *ver
			}
			if feat, ok := properties.GetFeaturesOk(); ok && feat != nil {
				dcPrint.Features = *feat
			}
			if secAuth, ok := properties.GetSecAuthProtectionOk(); ok && secAuth != nil {
				dcPrint.SecAuthProtection = *secAuth
			}
			if cpuArhis, ok := properties.GetCpuArchitectureOk(); ok && cpuArhis != nil {
				cpufamilies := make([]string, 0)
				for _, cpuArhi := range *cpuArhis {
					if cpuName, ok := cpuArhi.GetCpuFamilyOk(); ok && cpuName != nil {
						cpufamilies = append(cpufamilies, *cpuName)
					}
				}
				dcPrint.CpuFamily = cpufamilies
			}
		}
		if metadata, ok := dc.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				dcPrint.State = *state
			}
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
