package commands

import (
	"context"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkloadbalancerFlowLogCmd() *core.Command {
	ctx := context.TODO()
	networkloadbalancerFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "Network Load Balancer FlowLog Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer flowlog` allow you to create, list, get, update, delete Network Load Balancer FlowLogs.",
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
		LongDesc:   "Use this command to list Network Load Balancer FlowLogs from a specified Network Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id",
		Example:    listNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerFlowLogList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgNetworkLoadBalancerId, "", "", cloudapi_v6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
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
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgNetworkLoadBalancerId, "", "", cloudapi_v6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgFlowLogId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
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
* Bucket Name`,
		Example:    createNetworkLoadBalancerFlowLogExample,
		PreCmdRun:  PreRunNetworkLoadBalancerFlowLogCreate,
		CmdRun:     RunNetworkLoadBalancerFlowLogCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgNetworkLoadBalancerId, "", "", cloudapi_v6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "Unnamed FlowLog", "The name for the FlowLog")
	create.AddStringFlag(cloudapi_v6.ArgAction, cloudapi_v6.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgDirection, cloudapi_v6.ArgDirectionShort, "BIDIRECTIONAL", "Specifies the traffic Direction pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgS3Bucket, cloudapi_v6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapi_v6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
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
	update.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgNetworkLoadBalancerId, "", "", cloudapi_v6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgFlowLogId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.FlowLogId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "", "Name of the Network Load Balancer FlowLog")
	update.AddStringFlag(cloudapi_v6.ArgAction, cloudapi_v6.ArgActionShort, "", "Specifies the traffic Action pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgDirection, cloudapi_v6.ArgDirectionShort, "", "Specifies the traffic Direction pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgS3Bucket, cloudapi_v6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapi_v6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
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
		LongDesc: `Use this command to delete a specified Network Load Balancer FlowLog from a Network Load Balancer.

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
	deleteCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgNetworkLoadBalancerId, "", "", cloudapi_v6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgFlowLogId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancerFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapi_v6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog deletion [seconds]")

	return networkloadbalancerFlowLogCmd
}

func PreRunNetworkLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgNetworkLoadBalancerId, cloudapi_v6.ArgS3Bucket)
}

func PreRunDcNetworkLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgNetworkLoadBalancerId, cloudapi_v6.ArgFlowLogId)
}

func RunNetworkLoadBalancerFlowLogList(c *core.CommandConfig) error {
	networkloadbalancerFlowLogs, _, err := c.CloudApiV6Services.NetworkLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNetworkLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(networkloadbalancerFlowLogs)))
}

func RunNetworkLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	c.Printer.Verbose("NetworkLoadBalancerFlowLog with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId)))
	ng, _, err := c.CloudApiV6Services.NetworkLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, []resources.FlowLog{*ng}))
}

func RunNetworkLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	proper := getFlowLogPropertiesSet(c)
	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgAction)))
	}
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNetworkLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, []resources.FlowLog{*ng}))
}

func RunNetworkLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	input := getFlowLogPropertiesUpdate(c)
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId)),
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
	c.Printer.Verbose("NetworkLoadBalancerFlowLog with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId)))
	resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId)),
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
		config.GetServerUrl(),
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
