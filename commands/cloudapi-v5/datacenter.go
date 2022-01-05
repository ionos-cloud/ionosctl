package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DatacenterCmd() *core.Command {
	ctx := context.TODO()
	datacenterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"d", "dc", "vdc"},
			Short:            "Data Center Operations",
			Long:             "The sub-commands of `ionosctl datacenter` allow you to create, list, get, update and delete Data Centers.",
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultDatacenterCols, printer.ColsMessage(allDatacenterCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(datacenterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = datacenterCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allDatacenterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, datacenterCmd, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Data Centers",
		LongDesc: `Use this command to retrieve a complete list of Virtual Data Centers provisioned under your account. You can setup multiple query parameters.

You can filter the results using ` + "`" + `--filters` + "`" + ` option. Use the following format to set filters: ` + "`" + `--filters KEY1=VALUE1,KEY2=VALUE2` + "`" + `
` + completer.DataCentersFiltersUsage(),
		Example:    listDatacenterExample,
		PreCmdRun:  PreRunDataCenterList,
		CmdRun:     RunDataCenterList,
		InitClient: true,
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersFilters(), cobra.ShellCompDirectiveNoFileComp
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
		Example:    getDatacenterExample,
		LongDesc:   "Use this command to retrieve details about a Virtual Data Center by using its ID.\n\nRequired values to run command:\n\n* Data Center Id",
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv5.ArgDataCenterId, cloudapiv5.ArgIdShort, "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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

Virtual Data Centers are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.`,
		Example:    createDatacenterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunDataCenterCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "Unnamed Data Center", "Name of the Data Center")
	create.AddStringFlag(cloudapiv5.ArgDescription, cloudapiv5.ArgDescriptionShort, "", "Description of the Data Center")
	create.AddStringFlag(cloudapiv5.ArgLocation, cloudapiv5.ArgLocationShort, "de/txl", "Location for the Data Center")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	update.AddStringFlag(cloudapiv5.ArgDataCenterId, cloudapiv5.ArgIdShort, "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "Name of the Data Center")
	update.AddStringFlag(cloudapiv5.ArgDescription, cloudapiv5.ArgDescriptionShort, "", "Description of the Data Center")
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
		LongDesc: `Use this command to delete a specified Virtual Data Center from your account. This will remove all objects within the VDC and remove the VDC object itself.

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Data Center Id`,
		Example:    deleteDatacenterExample,
		PreCmdRun:  PreRunDataCenterDelete,
		CmdRun:     RunDataCenterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgDataCenterId, cloudapiv5.ArgIdShort, "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Data Center deletion")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Datacenters.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Data Center deletion [seconds]")

	return datacenterCmd
}

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId)
}

func PreRunDataCenterList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.DataCentersFilters(), completer.DataCentersFiltersUsage())
	}
	return nil
}

func PreRunDataCenterDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgDataCenterId},
		[]string{cloudapiv5.ArgAll},
	)
}

func RunDataCenterList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	datacenters, resp, err := c.CloudApiV5Services.DataCenters().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(nil, c, getDataCenters(datacenters)))
}

func RunDataCenterGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	dc, resp, err := c.CloudApiV5Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(nil, c, []resources.Datacenter{*dc}))
}

func RunDataCenterCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDescription))
	location := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLocation))
	c.Printer.Verbose("Properties set for creating the Datacenter: Name: %v, Description: %v, Location: %v", name, description, location)
	dc, resp, err := c.CloudApiV5Services.DataCenters().Create(name, description, location)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(resp, c, []resources.Datacenter{*dc}))
}

func RunDataCenterUpdate(c *core.CommandConfig) error {
	input := resources.DatacenterProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDescription))
		input.SetDescription(description)
		c.Printer.Verbose("Property Description set: %v", description)
	}
	dc, resp, err := c.CloudApiV5Services.DataCenters().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)), input)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getDataCenterPrint(resp, c, []resources.Datacenter{*dc}))
}

func RunDataCenterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll)) {
		if err := DeleteAllDatacenters(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete data center"); err != nil {
			return err
		}
		dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
		c.Printer.Verbose("Starting deleting Datacenter with id: %v...", dcId)
		resp, err := c.CloudApiV5Services.DataCenters().Delete(dcId)
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getDataCenterPrint(resp, c, nil))
	}
}

func DeleteAllDatacenters(c *core.CommandConfig) error {
	_ = c.Printer.Print("Datacenters to be deleted:")
	datacenters, _, err := c.CloudApiV5Services.DataCenters().List(resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if datacentersItems, ok := datacenters.GetItemsOk(); ok && datacentersItems != nil {
		for _, dc := range *datacentersItems {
			toPrint := ""
			if id, ok := dc.GetIdOk(); ok && id != nil {
				toPrint += "Datacenter Id: " + *id
			}
			if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					toPrint += " Datacenter Name: " + *name
				}
			}
			_ = c.Printer.Print(toPrint)
		}
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Datacenters"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the Datacenters...")
		var multiErr error
		for _, dc := range *datacentersItems {
			if id, ok := dc.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting Datacenter with id: %v...", *id)
				resp, err := c.CloudApiV5Services.DataCenters().Delete(*id)
				if resp != nil && printer.GetId(resp) != "" {
					c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
				}
				if err != nil {
					multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
					continue
				} else {
					_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return err
				}
			}
			_ = c.Printer.Print("\n")
		}
		if multiErr != nil {
			return multiErr
		}
		return nil
	} else {
		return errors.New("could not get items of Datacenters")
	}
}

// Output Printing

var (
	defaultDatacenterCols = []string{"DatacenterId", "Name", "Location", "Features", "State"}
	allDatacenterCols     = []string{"DatacenterId", "Name", "Location", "State", "Description", "Version", "Features", "SecAuthProtection"}
)

type DatacenterPrint struct {
	DatacenterId      string   `json:"DatacenterId,omitempty"`
	Name              string   `json:"Name,omitempty"`
	Location          string   `json:"Location,omitempty"`
	Description       string   `json:"Description,omitempty"`
	Version           int32    `json:"Version,omitempty"`
	State             string   `json:"State,omitempty"`
	Features          []string `json:"Features,omitempty"`
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
	if items, ok := datacenters.GetItemsOk(); ok && items != nil {
		for _, datacenter := range *items {
			dc = append(dc, resources.Datacenter{Datacenter: datacenter})
		}
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
