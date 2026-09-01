package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayRuleCreateCmd() *core.Command {
	create := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway Rule",
		LongDesc: `Use this command to create a source-NAT (SNAT) rule on a NAT Gateway. The rule masquerades outbound packets whose source matches ` + "`" + `--source-subnet` + "`" + ` (and, optionally, whose destination matches ` + "`" + `--target-subnet` + "`" + ` / a target port range) behind the public IP given in ` + "`" + `--ip` + "`" + `.

The ` + "`" + `--ip` + "`" + ` value must be one of the public IPs already assigned to the parent NAT Gateway (` + "`" + `natgateway create/update --ips` + "`" + `); an address not on the gateway is rejected.

Protocol / port constraints: a target port range (` + "`" + `--port-range-start` + "`" + ` / ` + "`" + `--port-range-end` + "`" + `) is only meaningful for TCP and UDP. If ` + "`" + `--protocol` + "`" + ` is ICMP the target port range cannot be set. With the default ALL, leave the port range at its default.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state before returning.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Public IP
* Source Subnet`,
		Example: `# Masquerade all outbound traffic from the 10.0.1.0/24 LAN behind a gateway public IP
ionosctl compute natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name allow-lan --ip 203.0.113.10 --source-subnet 10.0.1.0/24

# TCP-only rule limited to a destination subnet and HTTPS port range
ionosctl compute natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name https-out --ip 203.0.113.10 --source-subnet 10.0.1.0/24 --target-subnet 198.51.100.0/24 --protocol TCP --port-range-start 443 --port-range-end 443`,
		PreCmdRun:  PreRunNatGatewayRuleCreate,
		CmdRun:     RunNatGatewayRuleCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Rule", "Human-friendly name for the rule")
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, string(ionoscloud.ALL), "Protocol the rule matches: TCP, UDP, ICMP or ALL (default ALL matches every protocol). A target port range is only valid for TCP/UDP; with ICMP the target port range must not be set", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIpFlag(cloudapiv6.ArgIp, "", nil, "Public IP used to masquerade the source address of matched outbound packets. Must be one of the public IPs already assigned to the parent NAT Gateway", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgSourceSubnet, "", "", "Source subnet (CIDR) the rule applies to, matched against each packet's source IP; typically the CIDR of the private LAN whose servers should get outbound access, e.g. 10.0.1.0/24", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgTargetSubnet, "", "", "Destination subnet (CIDR) the rule applies to, matched against each packet's destination IP. Leave unset to match any destination")
	create.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "First destination port (inclusive) the rule matches. Only applies to TCP/UDP; ignored for ICMP/ALL")
	create.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Last destination port (inclusive) the rule matches. Only applies to TCP/UDP; ignored for ICMP/ALL")

	return create
}
