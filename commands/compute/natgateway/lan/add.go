package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayLanAddCmd() *core.Command {
	add := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "lan",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a NAT Gateway Lan",
		LongDesc: `Use this command to attach a private LAN to a NAT Gateway so servers on that LAN can route their outbound traffic through it. The gateway becomes reachable on the LAN via the gateway IPs given in ` + "`" + `--ips` + "`" + ` (the next-hop address servers use to reach the internet).

If ` + "`" + `--ips` + "`" + ` is not set, a gateway IP is generated automatically (with a /24 subnet). Gateway IPs must include a valid subnet mask and should belong to the same subnet as the LAN.

Attaching the LAN does not by itself translate traffic; add SNAT rules (` + "`" + `natgateway rule create` + "`" + `) whose ` + "`" + `--source-subnet` + "`" + ` covers the servers on this LAN.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Lan Id`,
		Example: `# Attach LAN 1 and let the gateway IP be auto-assigned
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id 1

# Attach LAN 1 with explicit gateway IPs (include the subnet mask)
ionosctl compute natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id 1 --ips 10.0.1.1/24`,
		PreCmdRun:  PreRunDcNatGatewayLanIds,
		CmdRun:     RunNatGatewayLanAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 1, cloudapiv6.LanId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Comma-separated gateway IPs (with subnet mask, e.g. 10.0.1.1/24) that the gateway uses on this LAN as the servers' next hop. Should belong to the LAN's subnet. If omitted, an IP is auto-generated with a /24 subnet")
	add.AddColsFlag(allCols)

	return add
}
