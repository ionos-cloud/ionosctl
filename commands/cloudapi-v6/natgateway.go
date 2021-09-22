package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultNatGatewayCols, printer.ColsMessage(defaultNatGatewayCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(natgatewayCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = natgatewayCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		LongDesc:   "Use this command to list NAT Gateways from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listNatGatewayExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunNatGatewayList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified NAT Gateway to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for waiting for NAT Gateway to be in AVAILABLE state [seconds]")

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
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "NAT Gateway", "Name of the NAT Gateway")
	create.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of public reserved IP addresses of the NAT Gateway", core.RequiredFlagOption())
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway creation [seconds]")

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
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the NAT Gateway")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of public reserved IP addresses of the NAT Gateway. This will overwrite the current values")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway update [seconds]")

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
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway deletion [seconds]")

	natgatewayCmd.AddCommand(NatgatewayRuleCmd())
	natgatewayCmd.AddCommand(NatgatewayLanCmd())
	natgatewayCmd.AddCommand(NatgatewayFlowLogCmd())

	return natgatewayCmd
}

func PreRunDcIdsNatGatewayIps(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIps)
}

func PreRunDcNatGatewayIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId)
}

func RunNatGatewayList(c *core.CommandConfig) error {
	natgateways, resp, err := c.CloudApiV6Services.NatGateways().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(nil, c, getNatGateways(natgateways)))
}

func RunNatGatewayGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, waiter.NatGatewayStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))); err != nil {
		return err
	}
	c.Printer.Verbose("NatGateway with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)))
	ng, resp, err := c.CloudApiV6Services.NatGateways().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(nil, c, []resources.NatGateway{*ng}))
}

func RunNatGatewayCreate(c *core.CommandConfig) error {
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
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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
	input := getNewNatGatewayInfo(c)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		*input,
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
	return c.Printer.Print(getNatGatewayPrint(resp, c, []resources.NatGateway{*ng}))
}

func RunNatGatewayDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nat gateway"); err != nil {
		return err
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	c.Printer.Verbose("NatGateway with id: %v is deleting...", natGatewayId)
	resp, err := c.CloudApiV6Services.NatGateways().Delete(dcId, natGatewayId)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayPrint(resp, c, nil))
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

// Output Printing

var defaultNatGatewayCols = []string{"NatGatewayId", "Name", "PublicIps", "State"}

type NatGatewayPrint struct {
	NatGatewayId string   `json:"NatGatewayId,omitempty"`
	Name         string   `json:"Name,omitempty"`
	PublicIps    []string `json:"PublicIps,omitempty"`
	State        string   `json:"State,omitempty"`
}

func getNatGatewayPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NatGateway) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getNatGatewaysKVMaps(ss)
			r.Columns = getNatGatewaysCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getNatGatewaysCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultNatGatewayCols
	}

	columnsMap := map[string]string{
		"NatGatewayId": "NatGatewayId",
		"Name":         "Name",
		"PublicIps":    "PublicIps",
		"State":        "State",
	}
	var natgatewayCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			natgatewayCols = append(natgatewayCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return natgatewayCols
}

func getNatGateways(natgateways resources.NatGateways) []resources.NatGateway {
	ss := make([]resources.NatGateway, 0)
	for _, s := range *natgateways.Items {
		ss = append(ss, resources.NatGateway{NatGateway: s})
	}
	return ss
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
		o := structs.Map(natgatewayPrint)
		out = append(out, o)
	}
	return out
}
