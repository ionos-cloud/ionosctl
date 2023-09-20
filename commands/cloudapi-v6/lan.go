package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allLanJSONPaths = map[string]string{
		"LanId":  "id",
		"Name":   "properties.name",
		"Public": "properties.public",
		"PccId":  "properties.pcc",
		"State":  "metadata.state",
	}

	defaultLanCols = []string{"LanId", "Name", "Public", "PccId", "State"}
	allLanCols     = []string{"LanId", "Name", "Public", "PccId", "State", "DatacenterId"}
)

func LanCmd() *core.Command {
	ctx := context.TODO()
	lanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Aliases:          []string{"l"},
			Short:            "LAN Operations",
			Long:             "The sub-commands of `ionosctl lan` allow you to create, list, get, update, delete LANs.",
			TraverseChildren: true,
		},
	}
	globalFlags := lanCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLanCols, tabheaders.ColsMessage(allLanCols))
	_ = viper.BindPFlag(core.GetFlagName(lanCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = lanCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace:  "lan",
		Resource:   "lan",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List LANs",
		LongDesc:   "Use this command to retrieve a list of LANs provisioned in a specific Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LANsFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLanExample,
		PreCmdRun:  PreRunLansList,
		CmdRun:     RunLanList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LANsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LANsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace:  "lan",
		Resource:   "lan",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a LAN",
		LongDesc:   "Use this command to retrieve information of a given LAN.\n\nRequired values to run command:\n\n* Data Center Id\n* LAN Id",
		Example:    getLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a LAN",
		LongDesc: `Use this command to create a new LAN within a Virtual Data Center on your account. The name, the public option and the Private Cross-Connect Id can be set.

NOTE: IP Failover is configured after LAN creation using an update command.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createLanExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLanCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed LAN", "The name of the LAN")
	create.AddBoolFlag(cloudapiv6.ArgPublic, cloudapiv6.ArgPublicShort, cloudapiv6.DefaultPublic, "Indicates if the LAN faces the public Internet (true) or not (false). E.g.: --public=true, --public=false")
	create.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for LAN creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for LAN creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a LAN",
		LongDesc: `Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Private Cross-Connect.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example:    updateLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the LAN")
	update.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.ArgPublic, "", cloudapiv6.DefaultPublic, "Public option for LAN. E.g.: --public=true, --public=false")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for LAN update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for LAN update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a LAN",
		LongDesc: `Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example:    deleteLanExample,
		PreCmdRun:  PreRunLanDelete,
		CmdRun:     RunLanDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for LAN deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Lans from a Virtual Data Center.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for LAN deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return lanCmd
}

func PreRunLansList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.LANsFilters(), completer.LANsFiltersUsage())
	}
	return nil
}

func PreRunLanDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunLanListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	allDcs := getDataCenters(datacenters)

	var allLans []ionoscloud.Lans
	var allLansConverted []map[string]interface{}
	totalTime := time.Duration(0)
	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("failed to retrieve Datacenter ID")
		}

		lans, resp, err := c.CloudApiV6Services.Lans().List(*dc.GetId(), listQueryParams)
		if err != nil {
			return err
		}

		items, ok := lans.GetItemsOk()
		if !ok || items == nil {
			continue
		}

		for _, item := range *items {
			temp, err := json2table.ConvertJSONToTable("", allLanJSONPaths, item)
			if err != nil {
				return fmt.Errorf("failed to convert from JSON to Table format: %w", err)
			}

			temp[0]["DatacenterId"] = *id
			allLansConverted = append(allLansConverted, temp[0])
		}

		allLans = append(allLans, lans.Lans)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(allLans, allLansConverted,
		tabheaders.GetHeaders(allLanCols, defaultLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLanList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunLanListAll(c)
	}

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	lans, resp, err := c.CloudApiV6Services.Lans().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", allLanJSONPaths, lans.Lans,
		tabheaders.GetHeadersAllDefault(defaultLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLanGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Lan with id: %v from Datacenter with id: %v is getting...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	l, resp, err := c.CloudApiV6Services.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allLanJSONPaths, l.Lan,
		tabheaders.GetHeadersAllDefault(defaultLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLanCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
	properties := ionoscloud.LanPropertiesPost{
		Name:   &name,
		Public: &public,
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Lan: Name: %v, Public: %v", name, public))

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		properties.SetPcc(pcc)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Pcc set: %v", pcc))
	}

	input := resources.LanPost{
		LanPost: ionoscloud.LanPost{
			Properties: &properties,
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Creating LAN in Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	l, resp, err := c.CloudApiV6Services.Lans().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), input, queryParams)
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

	out, err := jsontabwriter.GenerateOutput("", allLanJSONPaths, l.LanPost,
		tabheaders.GetHeadersAllDefault(defaultLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLanUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	input := resources.LanProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublic)) {
		public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
		input.SetPublic(public)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Public set: %v", public))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		input.SetPcc(pcc)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Pcc set: %v", pcc))
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating LAN with ID: %v from Datacenter with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
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

	out, err := jsontabwriter.GenerateOutput("", allLanJSONPaths, lanUpdated.Lan,
		tabheaders.GetHeadersAllDefault(defaultLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLanDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllLans(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete lan", viper.GetBool(constants.ArgForce)) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId))

	resp, err := c.CloudApiV6Services.Lans().Delete(dcId, lanId, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Lan successfully deleted"))
	return nil
}

func DeleteAllLans(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Datacenter ID: %v", dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Lans..."))

	lans, resp, err := c.CloudApiV6Services.Lans().List(dcId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	lansItems, ok := lans.GetItemsOk()
	if !ok || lansItems == nil {
		return fmt.Errorf("could not get items of Lans")
	}

	if len(*lansItems) <= 0 {
		return fmt.Errorf("no Lans found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Lans to be deleted:"))

	for _, lan := range *lansItems {
		delIdAndName := ""
		if id, ok := lan.GetIdOk(); ok && id != nil {
			delIdAndName += "Lan Id: " + *id
		}

		if properties, ok := lan.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " Lan Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Lans", viper.GetBool(constants.ArgForce)) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Lans..."))

	var multiErr error
	for _, lan := range *lansItems {
		id, ok := lan.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Lan with id: %v...", *id))

		resp, err = c.CloudApiV6Services.Lans().Delete(dcId, *id, queryParams)
		if resp != nil && utils.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Lans successfully deleted"))
	return nil
}
