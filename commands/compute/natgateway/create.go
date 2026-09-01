package natgateway

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func NatgatewayCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "natgateway",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway",
		LongDesc: `Use this command to create a NAT Gateway in a specified Virtual Data Center. The gateway provides source-NAT (SNAT) outbound internet access for servers on private LANs, masquerading their traffic behind the gateway's public IPs.

Creating the gateway only allocates it and its public IPs. To make traffic actually flow you still need to attach it to a private LAN (` + "`" + `natgateway lan add` + "`" + `) and add at least one SNAT rule (` + "`" + `natgateway rule create` + "`" + `).

The addresses passed to ` + "`" + `--ips` + "`" + ` must be public IPs you have already reserved in the same location as the datacenter (see ` + "`" + `ionosctl compute ipblock` + "`" + `); arbitrary or in-use addresses are rejected.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* IPs`,
		Example: `# Create a NAT Gateway with a single reserved public IP
ionosctl compute natgateway create --datacenter-id DATACENTER_ID --name my-gateway --ips 203.0.113.10

# Create a NAT Gateway with two public IPs and wait until it is AVAILABLE
ionosctl compute natgateway create --datacenter-id DATACENTER_ID --name my-gateway --ips 203.0.113.10,203.0.113.11 --wait`,
		PreCmdRun:  PreRunDcIdsNatGatewayIps,
		CmdRun:     RunNatGatewayCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "NAT Gateway", "Human-friendly name for the NAT Gateway")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Comma-separated public IP addresses the gateway masquerades outbound traffic behind. Must be IPs you have already reserved in the same location as the datacenter (see `ionosctl compute ipblock`). SNAT rules on this gateway can only reference IPs listed here", core.RequiredFlagOption())

	return cmd
}
