package commands

import (
	"context"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func networkloadbalancerFlowLog() *core.Command {
	ctx := context.TODO()
	networkloadbalancerFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "Network Load Balancer FlowLog Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl networkloadbalancer flowlog` + "`" + ` allow you to create, list, get, update, delete Network Load Balancer FlowLogs.`,
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, networkloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "flowlog",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancer FlowLogs",
		LongDesc:   "Use this command to list Network Load Balancer FlowLogs from a specified NAT Gateway.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id",
		Example:    listNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerFlowLogList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, utils.ColsMessage(defaultFlowLogCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, networkloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "flowlog",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Network Load Balancer FlowLog",
		LongDesc:   "Use this command to get information about a specified Network Load Balancer FlowLog from a Network Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id\n* Network Load Balancer FlowLog Id",
		Example:    getNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerFlowLogIds,
		CmdRun:     RunNetworkLoadBalancerFlowLogGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.RequiredFlagFlowLogId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, utils.ColsMessage(defaultFlowLogCols))
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, networkloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer FlowLog",
		LongDesc: `Use this command to create a Network Load Balancer FlowLog in a specified Network Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Name
* Direction
* Bucket Name`,
		Example:    createNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunNetworkLoadBalancerFlowLogCreate,
		CmdRun:     RunNetworkLoadBalancerFlowLogCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the FlowLog "+config.RequiredFlag)
	create.AddStringFlag(config.ArgAction, config.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDirection, config.ArgDirectionShort, "", "Specifies the traffic Direction pattern "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgBucketName, config.ArgBucketNameShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, utils.ColsMessage(defaultFlowLogCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, networkloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "flowlog",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer FlowLog",
		LongDesc: `Use this command to update a specified Network Load Balancer FlowLog from a Network Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id`,
		Example:    updateNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerFlowLogIds,
		CmdRun:     RunNetworkLoadBalancerFlowLogUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.RequiredFlagFlowLogId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Network Load Balancer FlowLog")
	update.AddStringFlag(config.ArgAction, config.ArgActionShort, "", "Specifies the traffic Action pattern")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgDirection, config.ArgDirectionShort, "", "Specifies the traffic Direction pattern")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgBucketName, config.ArgBucketNameShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, utils.ColsMessage(defaultFlowLogCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, networkloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "flowlog",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Network Load Balancer FlowLog",
		LongDesc: `Use this command to delete a specified Network Load Balancer FlowLog from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Network Load Balancer FlowLog Id`,
		Example:    deleteNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerFlowLogIds,
		CmdRun:     RunNetworkLoadBalancerFlowLogDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.RequiredFlagFlowLogId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog deletion [seconds]")

	return networkloadbalancerFlowLogCmd
}

func PreRunNetworkLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNetworkLoadBalancerId, config.ArgName, config.ArgDirection, config.ArgBucketName)
}

func PreRunDcNetworkLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNetworkLoadBalancerId, config.ArgFlowLogId)
}

func RunNetworkLoadBalancerFlowLogList(c *core.CommandConfig) error {
	networkloadbalancerFlowLogs, _, err := c.NetworkLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(networkloadbalancerFlowLogs)))
}

func RunNetworkLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	ng, _, err := c.NetworkLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, []resources.FlowLog{*ng}))
}

func RunNetworkLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	proper := getFlowLogPropertiesSet(c)
	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, config.ArgAction)))
	}
	ng, resp, err := c.NetworkLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, []resources.FlowLog{*ng}))
}

func RunNetworkLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	input := getFlowLogPropertiesSet(c)
	ng, resp, err := c.NetworkLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId)),
		&input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, []resources.FlowLog{*ng}))
}

func RunNetworkLoadBalancerFlowLogDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete network load balancer flowlog"); err != nil {
		return err
	}
	resp, err := c.NetworkLoadBalancers().DeleteFlowLog(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, nil))
}

func getNetworkLoadBalancerFlowLogsIds(outErr io.Writer, datacenterId, networkloadbalancerId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(clientSvc.Get(), context.TODO())
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
