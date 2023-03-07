package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerCmd() *core.Command {
	ctx := context.TODO()
	applicationloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "applicationloadbalancer",
			Aliases:          []string{"alb"},
			Short:            "Application Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl applicationloadbalancer` allow you to create, list, get, update, delete Application Load Balancers.",
			TraverseChildren: true,
		},
	}
	globalFlags := applicationloadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultApplicationLoadBalancerCols, printer.ColsMessage(defaultApplicationLoadBalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(applicationloadbalancerCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = applicationloadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultApplicationLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "applicationloadbalancer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Application Load Balancers",
		LongDesc:   "Use this command to list Application Load Balancers from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listApplicationLoadBalancerExample,
		PreCmdRun:  PreRunApplicationLoadBalancerList,
		CmdRun:     RunApplicationLoadBalancerList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "applicationloadbalancer",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an Application Load Balancer",
		LongDesc:   "Use this command to get information about a specified Application Load Balancer from a Virtual Data Center. You can also wait for Application Load Balancer to get in AVAILABLE state using `--wait-for-state` option.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id",
		Example:    getApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for specified Application Load Balancer to be in AVAILABLE state")
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an Application Load Balancer",
		LongDesc: `Use this command to create an Application Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunApplicationLoadBalancerCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Application Load Balancer", "The name of the Application Load Balancer.")
	create.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "ID of the listening (inbound) LAN.")
	create.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "ID of the balanced private target LAN (outbound).")
	create.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.")
	create.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.AlbTimeoutSeconds, "Timeout option for Request for Application Load Balancer creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an Application Load Balancer",
		LongDesc: `Use this command to update a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id`,
		Example:    updateApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Application Load Balancer", "The name of the Application Load Balancer.")
	update.AddIntFlag(cloudapiv6.ArgListenerLan, "", 0, "ID of the listening (inbound) LAN.")
	update.AddIntFlag(cloudapiv6.ArgTargetLan, "", 0, "ID of the balanced private target LAN (outbound).")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.")
	update.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete an Application Load Balancer",
		LongDesc: `Use this command to delete a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id`,
		Example:    deleteApplicationLoadBalancerExample,
		PreCmdRun:  PreRunApplicationLoadBalancerDelete,
		CmdRun:     RunApplicationLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Application Load Balancers")
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer deletion to be executed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerRuleCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerFlowLogCmd())

	return applicationloadbalancerCmd
}

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}

func PreRunApplicationLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func PreRunApplicationLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ApplicationLoadBalancersFilters(), completer.ApplicationLoadBalancersFiltersUsage())
	}
	return nil
}

func RunApplicationLoadBalancerListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	allDcs := getDataCenters(datacenters)
	var allApplicationLoadBalancers []resources.ApplicationLoadBalancer
	totalTime := time.Duration(0)
	for _, dc := range allDcs {
		ApplicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(*dc.GetId(), listQueryParams)
		if err != nil {
			return err
		}
		allApplicationLoadBalancers = append(allApplicationLoadBalancers, getApplicationLoadBalancers(ApplicationLoadBalancers)...)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Printer.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer.Print(getApplicationLoadBalancerPrint(nil, c, allApplicationLoadBalancers))
}

func RunApplicationLoadBalancerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunApplicationLoadBalancerListAll(c)
	}
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Getting ApplicationLoadBalancers from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationloadbalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(nil, c, getApplicationLoadBalancers(applicationloadbalancers)))
}

func RunApplicationLoadBalancerGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("Getting ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	if err := utils.WaitForState(c, waiter.ApplicationLoadBalancerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId))); err != nil {
		return err
	}
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(nil, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	proper := getNewApplicationLoadBalancerInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Printer.Verbose("Property TargetLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Printer.Verbose("Property ListenerLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	c.Printer.Verbose("Creating ApplicationLoadBalancer")
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		resources.ApplicationLoadBalancer{
			ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
				Properties: &proper.ApplicationLoadBalancerProperties,
			},
		},
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	input := getNewApplicationLoadBalancerInfo(c)
	c.Printer.Verbose("Updating ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		*input,
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	var (
		resp *resources.Response
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		err = DeleteAllApplicationLoadBalancer(c)
		if err != nil {
			return err
		}
	} else {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting ApplicationLoadBalancer")
		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), queryParams)
		if resp != nil {
			c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, nil))
}

func DeleteAllApplicationLoadBalancer(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	_ = c.Printer.Warn("Getting Application Load Balancers...")
	applicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if albItems, ok := applicationLoadBalancers.GetItemsOk(); ok && albItems != nil {
		if len(*albItems) > 0 {
			for _, alb := range *albItems {
				delIdAndName := ""
				if id, ok := alb.GetIdOk(); ok && id != nil {
					delIdAndName += "Application Load Balancer Id: " + *id
				}
				if propertiesOk, ok := alb.GetPropertiesOk(); ok && propertiesOk != nil {
					if name, ok := propertiesOk.GetNameOk(); ok && name != nil {
						delIdAndName += "Application Load Balancer Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Application Load Balancers"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Application Load Balancers...")
			var multiErr error
			for _, alb := range *albItems {
				if id, ok := alb.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Application Load Balancer with id: %v...", *id)
					resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Application Load Balancers found")
		}
	} else {
		return errors.New("could not get items of Application Load Balancers")
	}
}

func getNewApplicationLoadBalancerInfo(c *core.CommandConfig) *resources.ApplicationLoadBalancerProperties {
	input := ionoscloud.ApplicationLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		input.SetIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
		c.Printer.Verbose("Property IPs set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)) {
		input.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Printer.Verbose("Property ListenerLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)) {
		input.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Printer.Verbose("Property TargetLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)) {
		input.SetLbPrivateIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
		c.Printer.Verbose("Property LbPrivateIps set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
	}
	return &resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: input,
	}
}

// Output Printing

var (
	defaultApplicationLoadBalancerCols = []string{"ApplicationLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
	allApplicationLoadBalancerCols     = []string{"ApplicationLoadBalancerId", "DatacenterId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
)

type ApplicationLoadBalancerPrint struct {
	ApplicationLoadBalancerId string   `json:"ApplicationLoadBalancerId,omitempty"`
	Name                      string   `json:"Name,omitempty"`
	ListenerLan               int32    `json:"ListenerLan,omitempty"`
	Ips                       []string `json:"Ips,omitempty"`
	TargetLan                 int32    `json:"TargetLan,omitempty"`
	PrivateIps                []string `json:"PrivateIps,omitempty"`
	State                     string   `json:"State,omitempty"`
	DatacenterId              string   `json:"DatacenterId,omitempty"`
}

func getApplicationLoadBalancerPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.ApplicationLoadBalancer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getApplicationLoadBalancersKVMaps(ss)
			r.Columns = printer.GetHeadersListAll(allApplicationLoadBalancerCols, defaultApplicationLoadBalancerCols, "DatacenterId", viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)), viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)))
		}
	}
	return r
}

func getApplicationLoadBalancers(applicationloadbalancers resources.ApplicationLoadBalancers) []resources.ApplicationLoadBalancer {
	ss := make([]resources.ApplicationLoadBalancer, 0)
	for _, s := range *applicationloadbalancers.Items {
		ss = append(ss, resources.ApplicationLoadBalancer{ApplicationLoadBalancer: s})
	}
	return ss
}

func getApplicationLoadBalancersKVMaps(ss []resources.ApplicationLoadBalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var applicationloadbalancerPrint ApplicationLoadBalancerPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			applicationloadbalancerPrint.ApplicationLoadBalancerId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				applicationloadbalancerPrint.Name = *name
			}
			if listenerLan, ok := properties.GetListenerLanOk(); ok && listenerLan != nil {
				applicationloadbalancerPrint.ListenerLan = *listenerLan
			}
			if ips, ok := properties.GetIpsOk(); ok && ips != nil {
				applicationloadbalancerPrint.Ips = *ips
			}
			if targetLan, ok := properties.GetTargetLanOk(); ok && targetLan != nil {
				applicationloadbalancerPrint.TargetLan = *targetLan
			}
			if lbPrivateIps, ok := properties.GetLbPrivateIpsOk(); ok && lbPrivateIps != nil {
				applicationloadbalancerPrint.PrivateIps = *lbPrivateIps
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				applicationloadbalancerPrint.State = *state
			}
		}
		if hrefOk, ok := s.GetHrefOk(); ok && hrefOk != nil {
			// Get parent resource ID using HREF: `.../datacenter/[PARENT_ID_WE_WANT]/alb/[ALB_ID]`
			applicationloadbalancerPrint.DatacenterId = strings.Split(strings.Split(*hrefOk, "datacenter")[1], "/")[1]
		}
		o := structs.Map(applicationloadbalancerPrint)
		out = append(out, o)
	}
	return out
}
