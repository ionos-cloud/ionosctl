package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func natgatewayLan() *core.Command {
	ctx := context.TODO()
	natgatewayLanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "NAT Gateway Lan Operations",
			Long:             "The sub-commands of `ionosctl natgateway lan` allow you to add, list, remove NAT Gateway Lans.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, natgatewayLanCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "lan",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List NAT Gateway Lans",
		LongDesc: `Use this command to list NAT Gateway Lans from a specified NAT Gateway.

Required values to run command:

* Data Center Id
* NAT Gateway Id`,
		Example:    listNatGatewayLanExample,
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayLanList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultNatGatewayLanCols, utils.ColsMessage(defaultNatGatewayLanCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, natgatewayLanCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "lan",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a NAT Gateway Lan",
		LongDesc: `Use this command to add a NAT Gateway Lan in a specified NAT Gateway.

If IPs are not set manually, using ` + "`" + `--ips` + "`" + ` option, an IP will be automatically assigned. IPs must contain valid subnet mask. If user will not provide any IP then system will generate an IP with /24 subnet.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id`,
		Example:    addNatGatewayLanExample,
		PreCmdRun:  PreRunDcNatGatewayLanIds,
		CmdRun:     RunNatGatewayLanAdd,
		InitClient: true,
	})
	add.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = add.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = add.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIntFlag(config.ArgLanId, config.ArgIdShort, 1, config.RequiredFlagLanId)
	_ = add.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringSliceFlag(config.ArgIps, "", []string{""}, "Collection of Gateway IPs. If not set, it will automatically reserve public IPs")
	add.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway Lan addition to be executed")
	add.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Lan addition [seconds]")
	add.AddStringSliceFlag(config.ArgCols, "", defaultNatGatewayLanCols, utils.ColsMessage(defaultNatGatewayLanCols))
	_ = add.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, natgatewayLanCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "lan",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a NAT Gateway Lan",
		LongDesc: `Use this command to remove a specified NAT Gateway Lan from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id`,
		Example:    removeNatGatewayLanExample,
		PreCmdRun:  PreRunDcNatGatewayLanIds,
		CmdRun:     RunNatGatewayLanRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddIntFlag(config.ArgLanId, config.ArgIdShort, 1, config.RequiredFlagLanId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway Lan deletion to be executed")
	removeCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Lan deletion [seconds]")
	removeCmd.AddStringSliceFlag(config.ArgCols, "", defaultNatGatewayLanCols, utils.ColsMessage(defaultNatGatewayLanCols))
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	return natgatewayLanCmd
}

func PreRunDcNatGatewayLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNatGatewayId, config.ArgLanId)
}

func RunNatGatewayLanList(c *core.CommandConfig) error {
	ng, resp, err := c.NatGateways().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayLanPrint(nil, c, getNatGatewayLans(ng)))
}

func RunNatGatewayLanAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId))
	ng, _, err := c.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Adding NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId)
	input := getNewNatGatewayLanInfo(c, ng)
	ng, resp, err := c.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayLanPrint(resp, c, getNatGatewayLans(ng)))
}

func RunNatGatewayLanRemove(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove nat gateway lan"); err != nil {
		return err
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId))
	ng, _, err := c.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Removing NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId)
	input := removeNatGatewayLanInfo(c, ng)
	ng, resp, err := c.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayLanPrint(resp, c, nil))
}

func getNewNatGatewayLanInfo(c *core.CommandConfig, oldNg *v6.NatGateway) *v6.NatGatewayProperties {
	var proper []ionoscloud.NatGatewayLanProperties
	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				proper = *lans
			}
		}
	}
	input := ionoscloud.NatGatewayLanProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgLanId)) {
		lanId := viper.GetInt32(core.GetFlagName(c.NS, config.ArgLanId))
		input.SetId(lanId)
		c.Printer.Verbose("Property Id set: %v", lanId)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIps)) {
		gatewayIps := viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgIps))
		input.SetGatewayIps(gatewayIps)
		c.Printer.Verbose("Property GatewayIps set: %v", gatewayIps)
	}
	proper = append(proper, input)
	return &v6.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

func removeNatGatewayLanInfo(c *core.CommandConfig, oldNg *v6.NatGateway) *v6.NatGatewayProperties {
	proper := make([]ionoscloud.NatGatewayLanProperties, 0)
	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					if id, ok := lanItem.GetIdOk(); ok && id != nil {
						if *id != viper.GetInt32(core.GetFlagName(c.NS, config.ArgLanId)) {
							proper = append(proper, lanItem)
						}
					}
				}
			}
		}
	}
	return &v6.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

// Output Printing

var defaultNatGatewayLanCols = []string{"NatGatewayLanId", "GatewayIps"}

type NatGatewayLanPrint struct {
	NatGatewayLanId int32    `json:"NatGatewayLanId,omitempty"`
	GatewayIps      []string `json:"GatewayIps,omitempty"`
}

func getNatGatewayLanPrint(resp *v6.Response, c *core.CommandConfig, ss []v6.NatGatewayLanProperties) printer.Result {
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
			r.KeyValue = getNatGatewayLansKVMaps(ss)
			r.Columns = getNatGatewayLansCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getNatGatewayLansCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultNatGatewayLanCols
	}

	columnsMap := map[string]string{
		"NatGatewayLanId": "NatGatewayLanId",
		"GatewayIps":      "GatewayIps",
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

func getNatGatewayLans(ng *v6.NatGateway) []v6.NatGatewayLanProperties {
	ss := make([]v6.NatGatewayLanProperties, 0)
	if ng != nil {
		if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					ss = append(ss, v6.NatGatewayLanProperties{
						NatGatewayLanProperties: lanItem,
					})
				}
			}
		}
	}
	return ss
}

func getNatGatewayLansKVMaps(ss []v6.NatGatewayLanProperties) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var natgatewayPrint NatGatewayLanPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			natgatewayPrint.NatGatewayLanId = *id
		}
		if ips, ok := s.GetGatewayIpsOk(); ok && ips != nil {
			natgatewayPrint.GatewayIps = *ips
		}
		o := structs.Map(natgatewayPrint)
		out = append(out, o)
	}
	return out
}
