package natgateway

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		LongDesc: `Use this command to update the name or public IPs of a specified NAT Gateway.

Note that ` + "`" + `--ips` + "`" + ` REPLACES the whole set of public IPs rather than appending to it, so you must pass every IP the gateway should keep. Removing an IP that a SNAT rule still references will break that rule, so update or delete the affected rules first.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id`,
		Example: `# Rename a NAT Gateway
ionosctl compute natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name renamed-gateway

# Replace the gateway's public IPs (this overwrites the existing set)
ionosctl compute natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --ips 203.0.113.10,203.0.113.12`,
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayUpdate,
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
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New human-friendly name for the NAT Gateway")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "New set of reserved public IP addresses (same location as the datacenter). This OVERWRITES the current set, so include every IP the gateway should keep; dropping an IP referenced by a SNAT rule breaks that rule")

	return update
}
