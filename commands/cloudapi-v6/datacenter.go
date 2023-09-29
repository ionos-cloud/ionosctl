package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allDatacenterJSONPaths = map[string]string{
		"DatacenterId":      "id",
		"Name":              "properties.name",
		"Location":          "properties.location",
		"Description":       "properties.description",
		"Version":           "properties.version",
		"State":             "metadata.state",
		"Features":          "properties.features",
		"CpuFamily":         "properties.cpuArchitecture.*.cpuFamily",
		"SecAuthProtection": "properties.secAuthProtection",
	}

	defaultDatacenterCols = []string{"DatacenterId", "Name", "Location", "CpuFamily", "State"}
	allDatacenterCols     = []string{"DatacenterId", "Name", "Location", "State", "Description", "Version", "Features", "CpuFamily", "SecAuthProtection"}
)

func DatacenterCmd() *core.Command {
	ctx := context.TODO()
	datacenterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"d", "dc", "vdc"},
			Args:             cobra.ExactValidArgs(1),
			Short:            "Data Center Operations",
			Long:             "The sub-commands of `ionosctl datacenter` allow you to create, list, get, update and delete Data Centers.",
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultDatacenterCols, tabheaders.ColsMessage(allDatacenterCols))
	_ = viper.BindPFlag(core.GetFlagName(datacenterCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = datacenterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Data Center", "Name of the Data Center")
	create.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Data Center")
	create.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "de/txl", "Location for the Data Center")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Data Center creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Data Center creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Data Center")
	update.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Data Center")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Data Center update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Data Center update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Data Center deletion")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Datacenters.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Data Center deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return datacenterCmd
}

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
}

func PreRunDataCenterDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	)
}

func PreRunDataCenterList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.DataCentersFilters(), completer.DataCentersFiltersUsage())
	}
	return nil
}

func RunDataCenterList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	datacenters, resp, err := c.CloudApiV6Services.DataCenters().List(listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", allDatacenterJSONPaths, datacenters.Datacenters,
		tabheaders.GetHeaders(allDatacenterCols, defaultDatacenterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	dc, resp, err := c.CloudApiV6Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allDatacenterJSONPaths, dc.Datacenter,
		tabheaders.GetHeaders(allDatacenterCols, defaultDatacenterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
	loc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the datacenter: Name: %v, Description: %v, Location: %v", name, description, loc))

	dc, resp, err := c.CloudApiV6Services.DataCenters().Create(name, description, loc, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allDatacenterJSONPaths, dc,
		tabheaders.GetHeaders(allDatacenterCols, defaultDatacenterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	input := resources.DatacenterProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
		input.SetDescription(description)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Description set: %v", description))
	}

	dc, resp, err := c.CloudApiV6Services.DataCenters().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		input,
		queryParams,
	)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allDatacenterJSONPaths, dc,
		tabheaders.GetHeaders(allDatacenterCols, defaultDatacenterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllDatacenters(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete data center", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Datacenter with ID: %v...", dcId))

	resp, err := c.CloudApiV6Services.DataCenters().Delete(dcId, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Datacenter successfully deleted"))

	return nil

}

func DeleteAllDatacenters(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Datacenters..."))

	datacenters, resp, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	datacentersItems, ok := datacenters.GetItemsOk()
	if !ok || datacentersItems == nil {
		return fmt.Errorf("could not get items of Datacenters")
	}

	if len(*datacentersItems) <= 0 {
		return fmt.Errorf("no Datacenters found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Datacenters to be deleted:"))

	for _, dc := range *datacentersItems {
		delIdAndName := ""
		if id, ok := dc.GetIdOk(); ok && id != nil {
			delIdAndName += "Datacenter Id: " + *id
		}

		if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " Datacenter Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Datacenters", viper.IsSet(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Datacenters..."))

	var multiErr error
	for _, dc := range *datacentersItems {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Datacenter with id: %v...", *id))

		resp, err = c.CloudApiV6Services.DataCenters().Delete(*id, resources.QueryParams{})
		if resp != nil && utils.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Datacenters successfully deleted"))
	return nil
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
