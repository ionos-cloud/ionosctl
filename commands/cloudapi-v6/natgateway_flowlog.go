package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayFlowLogCmd() *core.Command {
	ctx := context.TODO()
	natgatewayFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "NAT Gateway FlowLog Operations",
			Long:             "The sub-commands of `ionosctl natgateway flowlog` allow you to create, list, get, update, delete NAT Gateway FlowLogs.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, natgatewayFlowLogCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "flowlog",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List NAT Gateway FlowLogs",
		LongDesc:   "Use this command to list NAT Gateway FlowLogs from a specified NAT Gateway.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.FlowLogsFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id",
		Example:    listNatGatewayFlowLogExample,
		PreCmdRun:  PreRunNATGatewayFlowLogList,
		CmdRun:     RunNatGatewayFlowLogList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FlowLogsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FlowLogsFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, natgatewayFlowLogCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "flowlog",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a NAT Gateway FlowLog",
		LongDesc:   "Use this command to get information about a specified NAT Gateway FlowLog from a NAT Gateway.\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id\n* NAT Gateway FlowLog Id",
		Example:    getNatGatewayFlowLogExample,
		PreCmdRun:  PreRunDcNatGatewayFlowLogIds,
		CmdRun:     RunNatGatewayFlowLogGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayFlowLogsIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, natgatewayFlowLogCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway FlowLog",
		LongDesc: `Use this command to create a NAT Gateway FlowLog in a specified NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Bucket Name`,
		Example:    createNatGatewayFlowLogExample,
		PreCmdRun:  PreRunNatGatewayFlowLogCreate,
		CmdRun:     RunNatGatewayFlowLogCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed FlowLog", "The name for the FlowLog")
	create.AddStringFlag(cloudapiv6.FlagAction, cloudapiv6.FlagActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "BIDIRECTIONAL", "Specifies the traffic Direction pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway FlowLog creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway FlowLog creation [seconds]")
	create.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, natgatewayFlowLogCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "flowlog",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a NAT Gateway FlowLog",
		LongDesc: `Use this command to update a specified NAT Gateway FlowLog from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway FlowLog Id`,
		Example:    updateNatGatewayFlowLogExample,
		PreCmdRun:  PreRunDcNatGatewayFlowLogIds,
		CmdRun:     RunNatGatewayFlowLogUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayFlowLogsIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "Name of the NAT Gateway FlowLog")
	update.AddStringFlag(cloudapiv6.FlagAction, cloudapiv6.FlagActionShort, "", "Specifies the traffic Action pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "", "Specifies the traffic Direction pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway FlowLog update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway FlowLog update [seconds]")
	update.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, natgatewayFlowLogCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "flowlog",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a NAT Gateway FlowLog",
		LongDesc: `Use this command to delete a specified NAT Gateway FlowLog from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway FlowLog Id`,
		Example:    deleteNatGatewayFlowLogExample,
		PreCmdRun:  PreRunNatGatewayFlowlogDelete,
		CmdRun:     RunNatGatewayFlowLogDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayFlowLogsIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway FlowLog deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all Natgateway flowlogs.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway FlowLog deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	return core.WithConfigOverride(natgatewayFlowLogCmd, "compute", "")
}

func PreRunNATGatewayFlowLogList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.FlowLogsFilters(), completer.FlowLogsFiltersUsage())
	}
	return nil
}

func PreRunNatGatewayFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.ArgS3Bucket)
}

func PreRunNatGatewayFlowlogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagFlowLogId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagAll},
	)
}

func PreRunDcNatGatewayFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagFlowLogId)
}

func RunNatGatewayFlowLogList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	natgatewayFlowLogs, resp, err := c.CloudApiV6Services.NatGateways().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Flowlog, natgatewayFlowLogs.FlowLogs,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayFlowLogGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"NatGatewayFlowLogGet with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))))

	ng, resp, err := c.CloudApiV6Services.NatGateways().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayFlowLogCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	proper := getFlowLogPropertiesSet(c)

	ng, resp, err := c.CloudApiV6Services.NatGateways().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayFlowLogUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	input := getFlowLogPropertiesUpdate(c)

	ng, resp, err := c.CloudApiV6Services.NatGateways().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId)),
		&input,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayFlowLogDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	natgatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId))
	flowlogId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllNatGatewayFlowLogs(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nat gateway flowlog", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting NatGatewayFlowLog with id: %v...", flowlogId))

	resp, err := c.CloudApiV6Services.NatGateways().DeleteFlowLog(dcId, natgatewayId, flowlogId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("NAT Gateway Flowlog successfully deleted"))
	return nil
}

func DeleteAllNatGatewayFlowLogs(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	natgatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("NatGateway ID: %v", natgatewayId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting NatGatewayFlowLogs..."))

	flowlogs, resp, err := c.CloudApiV6Services.NatGateways().ListFlowLogs(dcId, natgatewayId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	natgatewaysItems, ok := flowlogs.GetItemsOk()
	if !ok || natgatewaysItems == nil {
		return fmt.Errorf("could not get items of NAT Gateway FlowLogs")
	}

	if len(*natgatewaysItems) <= 0 {
		return fmt.Errorf("no Nat Gateway FlowLogs found")
	}

	var multiErr error
	for _, natgateway := range *natgatewaysItems {
		name := natgateway.GetProperties().Name
		id := natgateway.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the NatGatewayFlowLog with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NatGateways().DeleteFlowLog(dcId, natgatewayId, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
