package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

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

var (
	defaultApplicationLoadBalancerCols = []string{"ApplicationLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
	allApplicationLoadBalancerCols     = []string{"ApplicationLoadBalancerId", "DatacenterId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
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
	globalFlags.StringSliceP(constants.FlagCols, "", defaultApplicationLoadBalancerCols, tabheaders.ColsMessage(defaultApplicationLoadBalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(applicationloadbalancerCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = applicationloadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allApplicationLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
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
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
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
	list.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, cloudapiv6.FlagListAllDescription)

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
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.FlagWaitForState, constants.FlagWaitForStateShort, constants.DefaultWait, "Wait for specified Application Load Balancer to be in AVAILABLE state")
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)
	get.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds]")

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
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed Application Load Balancer", "The name of the Application Load Balancer.")
	create.AddIntFlag(cloudapiv6.FlagListenerLan, "", 2, "ID of the listening (inbound) LAN.")
	create.AddIntFlag(cloudapiv6.FlagTargetLan, "", 1, "ID of the balanced private target LAN (outbound).")
	create.AddStringSliceFlag(cloudapiv6.FlagIps, "", nil, "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.")
	create.AddStringSliceFlag(cloudapiv6.FlagPrivateIps, "", nil, "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.")
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.AlbTimeoutSeconds, "Timeout option for Request for Application Load Balancer creation [seconds]")
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

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
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Application Load Balancer", "The name of the Application Load Balancer.")
	update.AddIntFlag(cloudapiv6.FlagListenerLan, "", 0, "ID of the listening (inbound) LAN.")
	update.AddIntFlag(cloudapiv6.FlagTargetLan, "", 0, "ID of the balanced private target LAN (outbound).")
	update.AddStringSliceFlag(cloudapiv6.FlagIps, "", nil, "Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.")
	update.AddStringSliceFlag(cloudapiv6.FlagPrivateIps, "", nil, "Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer update [seconds]")
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all Application Load Balancers")
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Application Load Balancer deletion to be executed")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerRuleCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerFlowLogCmd())

	return core.WithConfigOverride(applicationloadbalancerCmd, "compute", "")
}

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId)
}

func PreRunApplicationLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagAll},
	)
}

func PreRunApplicationLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId},
		[]string{cloudapiv6.FlagAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
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
	var allApplicationLoadBalancers = make([]ionoscloud.ApplicationLoadBalancers, 0)
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		dcId, ok := dc.GetIdOk()
		if !ok || dcId == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		ApplicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(*dcId, listQueryParams)
		if err != nil {
			return err
		}

		allApplicationLoadBalancers = append(allApplicationLoadBalancers, ApplicationLoadBalancers.ApplicationLoadBalancers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.ApplicationLoadBalancer, allApplicationLoadBalancers,
		tabheaders.GetHeaders(allApplicationLoadBalancerCols, defaultApplicationLoadBalancerCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		return RunApplicationLoadBalancerListAll(c)
	}

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting ApplicationLoadBalancers from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))

	applicationloadbalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(dcId, listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}

	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.ApplicationLoadBalancer, applicationloadbalancers,
		tabheaders.GetHeadersAllDefault(defaultApplicationLoadBalancerCols, cols))
	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	queryParams := listQueryParams.QueryParams
	fmt.Fprintf(c.Command.Command.ErrOrStderr(),
		jsontabwriter.GenerateVerboseOutput(
			constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

	if err := waitfor.WaitForState(c, waiter.ApplicationLoadBalancerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))); err != nil {
		return err
	}

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}

	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApplicationLoadBalancer, ng.ApplicationLoadBalancer,
		tabheaders.GetHeadersAllDefault(defaultApplicationLoadBalancerCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	queryParams := listQueryParams.QueryParams
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	proper := getNewApplicationLoadBalancerInfo(c)

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))))
	}

	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagTargetLan)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property TargetLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagTargetLan))))
	}

	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagListenerLan)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property ListenerLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagListenerLan))))
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Creating ApplicationLoadBalancer"))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		resources.ApplicationLoadBalancer{
			ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
				Properties: &proper.ApplicationLoadBalancerProperties,
			},
		},
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Request href: %v ", resp.Header.Get("location")))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApplicationLoadBalancer, ng.ApplicationLoadBalancer,
		tabheaders.GetHeadersAllDefault(defaultApplicationLoadBalancerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	queryParams := listQueryParams.QueryParams
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))

	input := getNewApplicationLoadBalancerInfo(c)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Updating ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		*input,
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApplicationLoadBalancer, ng.ApplicationLoadBalancer,
		tabheaders.GetHeadersAllDefault(defaultApplicationLoadBalancerCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunApplicationLoadBalancerDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))

		err = DeleteAllApplicationLoadBalancer(c)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete application load balancer", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting ApplicationLoadBalancer"))

	resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)), queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Application Load Balancer successfully deleted"))

	return nil
}

func DeleteAllApplicationLoadBalancer(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Getting Application Load Balancers..."))

	applicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)), cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	albItems, ok := applicationLoadBalancers.GetItemsOk()
	if !ok || albItems == nil {
		return errors.New("could not get items of Application Load Balancers")
	}

	if len(*albItems) <= 0 {
		return errors.New("no Application Load Balancers found")
	}

	var multiErr error

	for _, alb := range *albItems {
		id := alb.GetId()
		name := alb.Properties.Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Application Load Balancer Id: %s , Name: %s ", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)), *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
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

func getNewApplicationLoadBalancerInfo(c *core.CommandConfig) *resources.ApplicationLoadBalancerProperties {
	input := ionoscloud.ApplicationLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIps)) {
		input.SetIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.FlagIps)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IPs set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.FlagIps))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagListenerLan)) {
		input.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagListenerLan)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property ListenerLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagListenerLan))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagTargetLan)) {
		input.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagTargetLan)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property TargetLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagTargetLan))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagPrivateIps)) {
		input.SetLbPrivateIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.FlagPrivateIps)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property LbPrivateIps set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.FlagPrivateIps))))
	}
	return &resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: input,
	}
}
