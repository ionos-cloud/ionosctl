package commands

import (
	"context"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
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

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "flowlog",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Application Load Balancer FlowLog",
		LongDesc:   "Use this command to get information about a specified Application Load Balancer FlowLog from a Application Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id\n* Application Load Balancer FlowLog Id",
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

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Application Load Balancer FlowLog",
		LongDesc: `Use this command to create a Application Load Balancer FlowLog in a specified Application Load Balancer.

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
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed ALB Flow Log", "The name for the FlowLog")
	create.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "INGRESS", "Specifies the traffic Direction pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Application Load Balancer FlowLog",
		LongDesc: `Use this command to update a specified Application Load Balancer FlowLog from a Application Load Balancer.

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
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Application Load Balancer FlowLog")
	update.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "", "Specifies the traffic Action pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "", "Specifies the traffic Direction pattern")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer FlowLog update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer FlowLog update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, applicationloadbalancerFlowLogCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Application Load Balancer FlowLog",
		LongDesc: `Use this command to delete a specified Application Load Balancer FlowLog from a Application Load Balancer.

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
	c.Printer.Verbose("Getting FlowLogs for ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(applicationloadbalancerFlowLogs)))
}

func RunApplicationLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting FlowLog with ID: %v for ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, []resources.FlowLog{*ng}))
}

func RunApplicationLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
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
	c.Printer.Verbose("Creating FlowLog for ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	input := getFlowLogPropertiesSet(c)
	c.Printer.Verbose("Updating FlowLog with ID: %v from ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		&input,
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	var (
		resp *resources.Response
		err  error
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		resp, err = DeleteAllApplicationLoadBalancerFlowLog(c)
		if err != nil {
			return err
		}
	} else {
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer flowlog"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting FlowLog with ID: %v from ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		)
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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

func DeleteAllApplicationLoadBalancerFlowLog(c *core.CommandConfig) (*resources.Response, error) {
	_ = c.Printer.Print("Application Load Balancer FlowLogs to be deleted:")
	applicationLoadBalancerFlowlogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if err != nil {
		return nil, err
	}
	if itemsOk, ok := applicationLoadBalancerFlowlogs.GetItemsOk(); ok && itemsOk != nil {
		for _, alb := range *itemsOk {
			if idOk, ok := alb.GetIdOk(); ok && idOk != nil {
				_ = c.Printer.Print("Application Load Balancer FlowLog Id: " + *idOk)
			}
			if propertiesOk, ok := alb.GetPropertiesOk(); ok && propertiesOk != nil {
				if name, ok := propertiesOk.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print("Application Load Balancer FlowLog Name: " + *name)
				}
			}
		}
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Application Load Balancer FlowLogs"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Application Load Balancer FlowLogs...")
		for _, itemOk := range *itemsOk {
			if idOk, ok := itemOk.GetIdOk(); ok && idOk != nil {
				c.Printer.Verbose("Starting deleting Application Load Balancer FlowLog with id: %v...", *idOk)
				resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
					viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
					viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), *idOk)
				if resp != nil {
					c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
					c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
				}
				if err != nil {
					return nil, err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, err
}
