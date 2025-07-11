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

func ApplicationLoadBalancerFlowLogCmd() *core.Command {
	ctx := context.TODO()
	applicationloadbalancerFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "Application Load Balancer FlowLog Operations",
			Long:             "The sub-commands of `ionosctl applicationloadbalancer flowlog` allow you to create, list, get, update, delete Application Load Balancer FlowLogs.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "flowlog",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Application Load Balancer FlowLogs",
		LongDesc:   "Use this command to list Application Load Balancer FlowLogs from a specified Application Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id",
		Example:    listApplicationLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerFlowLogList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "flowlog",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an Application Load Balancer FlowLog",
		LongDesc:   "Use this command to get information about a specified Application Load Balancer FlowLog from an Application Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id\n* Application Load Balancer FlowLog Id",
		Example:    getApplicationLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerFlowLogIds,
		CmdRun:     RunApplicationLoadBalancerFlowLogGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an Application Load Balancer FlowLog",
		LongDesc: `Use this command to create an Application Load Balancer FlowLog in a specified Application Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Bucket Name`,
		Example:    createApplicationLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunApplicationLoadBalancerFlowLogCreate,
		CmdRun:     RunApplicationLoadBalancerFlowLogCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed ALB Flow Log", "The name of the Application Load Balancer FlowLog.")
	create.AddStringFlag(cloudapiv6.FlagAction, cloudapiv6.FlagActionShort, "ALL", "Specifies the traffic action pattern.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "INGRESS", "Specifies the traffic direction pattern.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 bucket name of an existing IONOS Cloud S3 bucket.", core.RequiredFlagOption())
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog creation [seconds]")
	create.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an Application Load Balancer FlowLog",
		LongDesc: `Use this command to update a specified Application Load Balancer FlowLog from an Application Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id`,
		Example:    updateApplicationLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerFlowLogIds,
		CmdRun:     RunApplicationLoadBalancerFlowLogUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "The name of the Application Load Balancer FlowLog.")
	update.AddStringFlag(cloudapiv6.FlagAction, cloudapiv6.FlagActionShort, "", "Specifies the traffic action pattern.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "", "Specifies the traffic direction pattern.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 bucket name of an existing IONOS Cloud S3 bucket.")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog update [seconds]")
	update.AddStringSliceFlag(constants.FlagCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete an Application Load Balancer FlowLog",
		LongDesc: `Use this command to delete a specified Application Load Balancer FlowLog from an Application Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id`,
		Example:    deleteApplicationLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunApplicationLoadBalancerFlowLogDelete,
		CmdRun:     RunApplicationLoadBalancerFlowLogDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.FlagFlowLogId, cloudapiv6.FlagIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all Application Load Balancer FlowLogs")
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	return core.WithConfigOverride(applicationloadbalancerFlowLogCmd, "compute", "")
}

func PreRunApplicationLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.ArgS3Bucket)
}

func PreRunApplicationLoadBalancerFlowLogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagFlowLogId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagAll},
	)
}

func PreRunDcApplicationLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagFlowLogId)
}

func RunApplicationLoadBalancerFlowLogList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting FlowLogs"))

	applicationloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Flowlog, applicationloadbalancerFlowLogs,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting FlowLog with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Datacenter ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"ApplicationLoadBalancer ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

	proper := getFlowLogPropertiesSet(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))))
	}

	if !proper.HasDirection() {
		proper.SetDirection(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDirection)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Direction set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDirection))))
	}

	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagAction)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Action set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagAction))))
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Creating FlowLog"))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog, tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))))

	input := getFlowLogPropertiesUpdate(c)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating FlowLog"))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId)),
		&input,
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Flowlog, ng.FlowLog,
		tabheaders.GetHeadersAllDefault(defaultFlowLogCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerFlowLogDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	var resp *resources.Response
	queryParams := listQueryParams.QueryParams

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

		err = DeleteAllApplicationLoadBalancerFlowLog(c)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateVerboseOutput(
		"FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete application load balancer flowlog with id: %s", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId))), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFlowLogId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Application Load Balancers Flowlog successfully deleted"))

	return nil
}

func DeleteAllApplicationLoadBalancerFlowLog(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Getting Application Load Balancer FlowLogs..."))

	applicationLoadBalancerFlowlogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		cloudapiv6.ParentResourceListQueryParams,
	)
	if err != nil {
		return err
	}

	albFlowLogItems, ok := applicationLoadBalancerFlowlogs.GetItemsOk()
	if !ok || albFlowLogItems == nil {
		return errors.New("could not get items of Application Load Balancer Flow Logs")
	}

	if len(*albFlowLogItems) <= 0 {
		return errors.New("no Application Load Balancer Flow Logs found")
	}

	var multiErr error
	for _, fl := range *albFlowLogItems {
		id := fl.GetId()
		name := fl.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Application Load Balancer FlowLog Id: %s , Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)), *id,
			queryParams,
		)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
