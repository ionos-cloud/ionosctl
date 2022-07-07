package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"go.uber.org/multierr"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
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
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultListDepth, "Controls the detail depth of the response objects. Max depth is 10.")
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgFlowLogId, cloudapiv6.ArgIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	get.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultGetDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed ALB Flow Log", "The name of the Application Load Balancer FlowLog.")
	create.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "ALL", "Specifies the traffic action pattern.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "INGRESS", "Specifies the traffic direction pattern.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 bucket name of an existing IONOS Cloud S3 bucket.", core.RequiredFlagOption())
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultCreateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgFlowLogId, cloudapiv6.ArgIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the Application Load Balancer FlowLog.")
	update.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "", "Specifies the traffic action pattern.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "", "Specifies the traffic direction pattern.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 bucket name of an existing IONOS Cloud S3 bucket.")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultUpdateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgFlowLogId, cloudapiv6.ArgIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Application Load Balancer FlowLogs")
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog deletion [seconds]")
	deleteCmd.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultDeleteDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	return applicationloadbalancerFlowLogCmd
}

func PreRunApplicationLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgS3Bucket)
}

func PreRunApplicationLoadBalancerFlowLogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgFlowLogId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func PreRunDcApplicationLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgFlowLogId)
}

func RunApplicationLoadBalancerFlowLogList(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("List Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
		if !structs.IsZero(listQueryParams.QueryParams) {
			c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams.QueryParams))
		}
	}
	c.Printer.Verbose("Getting FlowLogs")
	applicationloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(applicationloadbalancerFlowLogs)))
}

func RunApplicationLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("Getting FlowLog with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, []resources.FlowLog{*ng}))
}

func RunApplicationLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	c.Printer.Verbose("Datacenter ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	proper := getFlowLogPropertiesSet(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if !proper.HasDirection() {
		proper.SetDirection(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))
		c.Printer.Verbose("Property Direction set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))
	}
	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction)))
		c.Printer.Verbose("Property Action set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction)))
	}
	c.Printer.Verbose("Creating FlowLog")
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, []resources.FlowLog{*ng}))
}

func RunApplicationLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))
	input := getFlowLogPropertiesUpdate(c)
	c.Printer.Verbose("Updating FlowLog")
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		&input,
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, []resources.FlowLog{*ng}))
}

func RunApplicationLoadBalancerFlowLogDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	var (
		resp *resources.Response
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		err = DeleteAllApplicationLoadBalancerFlowLog(c)
		if err != nil {
			return err
		}
	} else {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		c.Printer.Verbose("FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer flowlog"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting FlowLog")
		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
			queryParams,
		)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, nil))
}

func DeleteAllApplicationLoadBalancerFlowLog(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	_ = c.Printer.Print("Getting Application Load Balancer FlowLogs...")
	applicationLoadBalancerFlowlogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.ListQueryParams{},
	)
	if err != nil {
		return err
	}
	if albFlowLogItems, ok := applicationLoadBalancerFlowlogs.GetItemsOk(); ok && albFlowLogItems != nil {
		if len(*albFlowLogItems) > 0 {
			for _, fl := range *albFlowLogItems {
				toPrint := ""
				if id, ok := fl.GetIdOk(); ok && id != nil {
					toPrint += "Application Load Balancer FlowLog Id: " + *id
				}
				if properties, ok := fl.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += "Application Load Balancer FlowLog Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Application Load Balancer FlowLogs"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Application Load Balancer FlowLogs...")
			var multiErr error
			for _, fl := range *albFlowLogItems {
				if id, ok := fl.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Application Load Balancer FlowLog with id: %v...", *id)
					resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
						viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
						viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), *id,
						queryParams,
					)
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
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Application Load Balancer Flow Logs found")
		}
	} else {
		return errors.New("could not get items of Application Load Balancer Flow Logs")
	}
}
