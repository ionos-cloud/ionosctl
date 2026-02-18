package lan

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayLanRemoveCmd() *core.Command {
	removeCmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
		Example:    "ionosctl natgateway lan remove --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunDcNatGatewayLanRemove,
		CmdRun:     cloudapiv6cmds.RunNatGatewayLanRemove,
		InitClient: true,
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 1, cloudapiv6.LanId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Lan deletion to be executed")
	removeCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Lan deletion [seconds]")
	removeCmd.AddStringSliceFlag(constants.ArgCols, "", defaultNatGatewayLanCols, tabheaders.ColsMessage(defaultNatGatewayLanCols))
	_ = removeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayLanCols, cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all NAT Gateway Lans.")

	return removeCmd
}
