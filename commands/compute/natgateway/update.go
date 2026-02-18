package natgateway

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayUpdateCmd() *core.Command {
	update := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
		Example:    "ionosctl natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME",
		PreCmdRun:  cloudapiv6cmds.PreRunDcNatGatewayIds,
		CmdRun:     cloudapiv6cmds.RunNatGatewayUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIdShort, "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the NAT Gateway")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of public reserved IP addresses of the NAT Gateway. This will overwrite the current values")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway update [seconds]")

	return update
}
