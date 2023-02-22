package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayCmd() *core.Command {
	ctx := context.TODO()
	natgatewayCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "natgateway",
			Aliases:          []string{"nat", "ng"},
			Short:            "NAT Gateway Operations",
			Long:             "The sub-commands of `ionosctl natgateway` allow you to create, list, get, update, delete NAT Gateways.",
			TraverseChildren: true,
		},
	}
	globalFlags := natgatewayCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNatGatewayCols, printer.ColsMessage(defaultNatGatewayCols))
	_ = viper.BindPFlag(core.GetFlagName(natgatewayCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = natgatewayCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, natgatewayCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "natgateway",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List NAT Gateways",
		LongDesc:   "Use this command to list NAT Gateways from a specified Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.NATGatewaysFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listNatGatewayExample,
		PreCmdRun:  PreRunNATGatewayList,
		CmdRun:     RunNatGatewayList,
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
		return completer.NATGatewaysFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NATGatewaysFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, natgatewayCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "natgateway",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a NAT Gateway",
		LongDesc:   "Use this command to get information about a specified NAT Gateway from a Virtual Data Center. You can also wait for NAT Gateway to get in AVAILABLE state using `--wait-for-state` option.\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id",
		Example:    getNatGatewayExample,
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for specified NAT Gateway to be in AVAILABLE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for waiting for NAT Gateway to be in AVAILABLE state [seconds]")
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, natgatewayCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "natgateway",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway",
		LongDesc: `Use this command to create a NAT Gateway in a specified Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* IPs`,
		Example:    createNatGatewayExample,
		PreCmdRun:  PreRunDcIdsNatGatewayIps,
		CmdRun:     RunNatGatewayCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "NAT Gateway", "Name of the NAT Gateway")
	create.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of public reserved IP addresses of the NAT Gateway", core.RequiredFlagOption())
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, natgatewayCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "natgateway",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a NAT Gateway",
		LongDesc: `Use this command to update a specified NAT Gateway from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id`,
		Example:    updateNatGatewayExample,
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the NAT Gateway")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of public reserved IP addresses of the NAT Gateway. This will overwrite the current values")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, natgatewayCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "natgateway",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a NAT Gateway",
		LongDesc: `Use this command to delete a specified NAT Gateway from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id`,
		Example:    deleteNatGatewayExample,
		PreCmdRun:  PreRunNatGatewayDelete,
		CmdRun:     RunNatGatewayDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Natgateways.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	natgatewayCmd.AddCommand(NatgatewayRuleCmd())
	natgatewayCmd.AddCommand(NatgatewayLanCmd())
	natgatewayCmd.AddCommand(NatgatewayFlowLogCmd())

	return natgatewayCmd
}

func PreRunNATGatewayList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.NATGatewaysFilters(), completer.NATGatewaysFiltersUsage())
	}
	return nil
}

func PreRunDcIdsNatGatewayIps(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIps)
}

func PreRunDcNatGatewayIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId)
}

func PreRunNatGatewayDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	allDcs := getDataCenters(datacenters)
	var allNatGateways []resources.NatGateway
	totalTime := time.Duration(0)
	for _, dc := range allDcs {
		natGateways, resp, err := c.CloudApiV6Services.NatGateways().List(*dc.GetId(), listQueryParams)
		if err != nil {
			return err
		}
		allNatGateways = append(allNatGateways, getNatGateways(natGateways)...)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Printer.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer.Print(getNatGatewayPrint(nil, c, allNatGateways))
}

func RunNatGatewayList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunNatGatewayListAll(c)
	}
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	natgateways, resp, err := c.CloudApiV6Services.NatGateways().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(nil, c, getNatGateways(natgateways)))
}

func RunNatGatewayGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if err := utils.WaitForState(c, waiter.NatGatewayStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))); err != nil {
		return err
	}
	c.Printer.Verbose("NatGateway with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)))
	ng, resp, err := c.CloudApiV6Services.NatGateways().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(nil, c, []resources.NatGateway{*ng}))
}

func RunNatGatewayCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	proper := getNewNatGatewayInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	ng, resp, err := c.CloudApiV6Services.NatGateways().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		resources.NatGateway{
			NatGateway: ionoscloud.NatGateway{
				Properties: &proper.NatGatewayProperties,
			},
		},
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(resp, c, []resources.NatGateway{*ng}))
}

func RunNatGatewayUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	input := getNewNatGatewayInfo(c)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		*input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(resp, c, []resources.NatGateway{*ng}))
}

func RunNatGatewayDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNatgateways(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nat gateway"); err != nil {
			return err
		}
		c.Printer.Verbose("Starring deleting NatGateway with id: %v...", natGatewayId)
		resp, err := c.CloudApiV6Services.NatGateways().Delete(dcId, natGatewayId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getNatGatewayPrint(resp, c, nil))
	}
}

func getNewNatGatewayInfo(c *core.CommandConfig) *resources.NatGatewayProperties {
	input := ionoscloud.NatGatewayProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
		input.SetPublicIps(publicIps)
		c.Printer.Verbose("Property PublicIps set: %v", publicIps)
	}
	return &resources.NatGatewayProperties{
		NatGatewayProperties: input,
	}
}

func DeleteAllNatgateways(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting NatGateways...")
	natGateways, resp, err := c.CloudApiV6Services.NatGateways().List(dcId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if natGatewayItems, ok := natGateways.GetItemsOk(); ok && natGatewayItems != nil {
		if len(*natGatewayItems) > 0 {
			_ = c.Printer.Warn("NatGateway to be deleted:")
			for _, natGateway := range *natGatewayItems {
				delIdAndName := ""
				if id, ok := natGateway.GetIdOk(); ok && id != nil {
					delIdAndName += "NatGateway Id: " + *id
				}
				if properties, ok := natGateway.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += " NatGateway Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the NatGateways"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the NatGateways...")
			var multiErr error
			for _, natGateway := range *natGatewayItems {
				if id, ok := natGateway.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting NatGateway with id: %v...", *id)
					resp, err = c.CloudApiV6Services.NatGateways().Delete(dcId, *id, queryParams)
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
			return errors.New("no NatGateways found")
		}
	} else {
		return errors.New("could not get items of NatGateway")
	}
}

// Output Printing

var defaultNatGatewayCols = []string{"NatGatewayId", "Name", "PublicIps", "State"}
var allNatGatewayCols = []string{"NatGatewayId", "Name", "PublicIps", "State", "DatacenterId"}

type NatGatewayPrint struct {
	NatGatewayId string   `json:"NatGatewayId,omitempty"`
	Name         string   `json:"Name,omitempty"`
	PublicIps    []string `json:"PublicIps,omitempty"`
	State        string   `json:"State,omitempty"`
	DatacenterId string   `json:"DatacenterId,omitempty"`
}

func getNatGatewayPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NatGateway) printer.Result {
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
			r.KeyValue = getNatGatewaysKVMaps(ss)
			r.Columns = printer.GetHeadersListAll(allNatGatewayCols, defaultNatGatewayCols, "DatacenterId", viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)), viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)))
		}
	}
	return r
}

func getNatGateways(natgateways resources.NatGateways) []resources.NatGateway {
	natGatewayObjs := make([]resources.NatGateway, 0)
	if items, ok := natgateways.GetItemsOk(); ok && items != nil {
		for _, natGateway := range *items {
			natGatewayObjs = append(natGatewayObjs, resources.NatGateway{NatGateway: natGateway})
		}
	}
	return natGatewayObjs
}

func getNatGatewaysKVMaps(ss []resources.NatGateway) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var natgatewayPrint NatGatewayPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			natgatewayPrint.NatGatewayId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				natgatewayPrint.Name = *name
			}
			if ips, ok := properties.GetPublicIpsOk(); ok && ips != nil {
				natgatewayPrint.PublicIps = *ips
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				natgatewayPrint.State = *state
			}
		}
		if hrefOk, ok := s.GetHrefOk(); ok && hrefOk != nil {
			// Get parent resource ID using HREF: `.../datacenter/[PARENT_ID_WE_WANT]/natgateways/[ID]`
			natgatewayPrint.DatacenterId = strings.Split(strings.Split(*hrefOk, "datacenter")[1], "/")[1]
		}
		o := structs.Map(natgatewayPrint)
		out = append(out, o)
	}
	return out
}
