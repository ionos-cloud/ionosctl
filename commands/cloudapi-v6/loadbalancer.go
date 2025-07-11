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
	defaultLoadbalancerCols = []string{"LoadBalancerId", "Name", "Dhcp", "State"}
	allLoadbalancerCols     = []string{"LoadBalancerId", "Name", "Dhcp", "State", "Ip", "DatacenterId"}
)

func LoadBalancerCmd() *core.Command {
	ctx := context.TODO()
	loadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "loadbalancer",
			Aliases:          []string{"lb"},
			Short:            "Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl loadbalancer` manage your Load Balancers on your account. With the `ionosctl loadbalancer` command, you can list, create, delete Load Balancers and manage their configuration details.",
			TraverseChildren: true,
		},
	}
	globalFlags := loadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.FlagCols, "", defaultLoadbalancerCols, tabheaders.ColsMessage(allLoadbalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(loadbalancerCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = loadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLoadbalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "loadbalancer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Load Balancers",
		LongDesc:   "Use this command to retrieve a list of Load Balancers within a Virtual Data Center on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LoadbalancersFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLoadbalancerExample,
		PreCmdRun:  PreRunLoadBalancerList,
		CmdRun:     RunLoadBalancerList,
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
		return completer.LoadBalancersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadBalancersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, cloudapiv6.FlagListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "loadbalancer",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Load Balancer",
		LongDesc:   "Use this command to retrieve information about a Load Balancer instance.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		Example:    getLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.LoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Load Balancer",
		LongDesc: `Use this command to create a new Load Balancer within the Virtual Data Center. Load balancers can be used for public or private IP traffic. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createLoadbalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLoadBalancerCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Load Balancer", "Name of the Load Balancer")
	create.AddBoolFlag(cloudapiv6.FlagDhcp, "", cloudapiv6.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Load Balancer creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer creation [seconds]")
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Load Balancer",
		LongDesc: `Use this command to update the configuration of a specified Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`,
		Example:    updateLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.LoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "Name of the Load Balancer")
	update.AddIpFlag(cloudapiv6.FlagIp, "", nil, "The IP of the Load Balancer")
	update.AddBoolFlag(cloudapiv6.FlagDhcp, "", cloudapiv6.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Load Balancer update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer update [seconds]")
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Load Balancer",
		LongDesc: `Use this command to delete the specified Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`,
		Example:    deleteLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerDelete,
		CmdRun:     RunLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagLoadBalancerId, cloudapiv6.FlagIdShort, "", cloudapiv6.LoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Load Balancer deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all the LoadBlancers from a virtual Datacenter.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	loadbalancerCmd.AddCommand(LoadBalancerNicCmd())

	return core.WithConfigOverride(loadbalancerCmd, "compute", "")
}

func PreRunLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId},
		[]string{cloudapiv6.FlagAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.LoadBalancersFilters(), completer.LoadbalancersFiltersUsage())
	}
	return nil
}

func PreRunDcLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagLoadBalancerId)
}

func PreRunDcLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagLoadBalancerId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagAll},
	)
}

func RunLoadBalancerListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	allDcs := getDataCenters(datacenters)

	var allLoadbalancers []ionoscloud.Loadbalancers
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		loadBalancers, resp, err := c.CloudApiV6Services.Loadbalancers().List(*id, listQueryParams)
		if err != nil {
			return err
		}

		allLoadbalancers = append(allLoadbalancers, loadBalancers.Loadbalancers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.LoadBalancer, allLoadbalancers,
		tabheaders.GetHeaders(allLoadbalancerCols, defaultLoadbalancerCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLoadBalancerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		return RunLoadBalancerListAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting LoadBalancers from Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	lbs, resp, err := c.CloudApiV6Services.Loadbalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.LoadBalancer, lbs.Loadbalancers,
		tabheaders.GetHeaders(allLoadbalancerCols, defaultLoadbalancerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLoadBalancerGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Load balancer with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLoadBalancerId))))

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLoadBalancerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.LoadBalancer, lb.Loadbalancer,
		tabheaders.GetHeaders(allLoadbalancerCols, defaultLoadbalancerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLoadBalancerCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDhcp))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the load balancer: Name: %v, Dhcp: %v", name, dhcp))

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Create(dcId, name, dhcp, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.LoadBalancer, lb.Loadbalancer,
		tabheaders.GetHeaders(allLoadbalancerCols, defaultLoadbalancerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLoadBalancerUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	input := resources.LoadbalancerProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIp)) {
		ip := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIp))
		input.SetIp(ip)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Ip set: %v", ip))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDhcp))
		input.SetDhcp(dhcp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Dhcp set: %v", dhcp))
	}

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLoadBalancerId)),
		input,
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

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.LoadBalancer, lb.Loadbalancer,
		tabheaders.GetHeaders(allLoadbalancerCols, defaultLoadbalancerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunLoadBalancerDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLoadBalancerId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllLoadBalancers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete loadbalancer", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Load balancer with id: %v is deleting...", loadBalancerId))

	resp, err := c.CloudApiV6Services.Loadbalancers().Delete(dcid, loadBalancerId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Load Balancer successfully deleted"))
	return nil
}

func DeleteAllLoadBalancers(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcid))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting LoadBalancers..."))

	loadBalancers, resp, err := c.CloudApiV6Services.Loadbalancers().List(dcid, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	loadBalancersItems, ok := loadBalancers.GetItemsOk()
	if !ok || loadBalancersItems == nil {
		return fmt.Errorf("could not get items of Load Balancers")
	}

	if len(*loadBalancersItems) <= 0 {
		return fmt.Errorf("no Load Balancers found")
	}

	var multiErr error
	for _, lb := range *loadBalancersItems {
		name := lb.GetProperties().Name
		id := lb.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the LoadBalancer with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Loadbalancers().Delete(dcid, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
