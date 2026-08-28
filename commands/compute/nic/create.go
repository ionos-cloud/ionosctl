package nic

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func NicCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NIC",
		LongDesc: `Use this command to add a new NIC (Network Interface Card) to a server. The NIC attaches the server (--server-id) to a LAN (--lan-id) inside the given Data Center (--datacenter-id). If the target LAN does not exist yet, it is created implicitly when the NIC is created.

Addressing options:
* DHCP (default): with --dhcp=true the NIC reserves an IP automatically from the LAN's DHCP server. This is the usual choice.
* Static/reserved IPs: pass --ips to assign specific addresses. Explicitly assigned public IPs must come from reserved IP blocks; private-range addresses (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be assigned on private LANs. Leaving --ips empty lets IONOS pick an address automatically.

Firewall: --firewall-active toggles the per-NIC firewall (off by default). When active, all incoming traffic is blocked except what explicit firewall rules allow; when inactive, rules are ignored and traffic reaches the NIC directly. --firewall-type selects which traffic direction those rules govern (INGRESS/EGRESS/BIDIRECTIONAL).

IPv6: --ipv6-cidr-block, --ipv6-ips and --dhcpv6 only apply when the target LAN has IPv6 enabled; --ipv6-ips must fall within the NIC's IPv6 CIDR block.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the NIC to reach AVAILABLE state before the command returns.

Required values to run a command:

* Data Center Id
* Server Id`,
		Example: `# Create a DHCP NIC on LAN 1 (defaults: --dhcp=true, --lan-id=1, firewall off)
ionosctl compute nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name mynic

# Create a NIC on a specific LAN with static reserved IPs and an active ingress firewall
ionosctl compute nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name web-nic --lan-id 2 --dhcp=false --ips 203.0.113.10 --firewall-active=true --firewall-type INGRESS --wait`,
		PreCmdRun:  PreRunNicCreate,
		CmdRun:     RunNicCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		dcId, _ := c.Flags().GetString(cloudapiv6.ArgDataCenterId)
		return completer.ServersIds(dcId), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Internet Access", "A human-friendly name for the NIC (shown in the DCD and listings). Does not affect networking")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "One or more IPs to assign to the NIC. Explicitly assigned public IPs must come from a reserved IP block; private-range IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be used on private LANs. Leave empty to let IONOS assign an address automatically (see --dhcp)")
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", cloudapiv6.DefaultDhcp, "Whether the NIC reserves an IP automatically from the LAN's DHCP server. Set --dhcp=false to disable DHCP (typically when assigning static --ips). Default: true")
	cmd.AddBoolFlag(cloudapiv6.ArgFirewallActive, "", cloudapiv6.DefaultFirewallActive, "Enable the per-NIC firewall. When enabled, an empty ruleset blocks all incoming traffic and only explicitly configured firewall rules are allowed through; when disabled, all traffic reaches the NIC and rules are ignored. Default: false")
	cmd.AddStringFlag(cloudapiv6.ArgFirewallType, "", "INGRESS", "Direction of traffic the NIC's firewall rules apply to. INGRESS = inbound only, EGRESS = outbound only, BIDIRECTIONAL = both. Only meaningful when --firewall-active=true. Default: INGRESS")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgLanId, "", cloudapiv6.DefaultNicLanId, "The ID of the LAN this NIC attaches to, determining which network the server reaches through this NIC. If the LAN ID does not exist in the Data Center, it is created implicitly. Default: 1")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		dcId, _ := c.Flags().GetString(cloudapiv6.ArgDataCenterId)
		return completer.LansIds(dcId), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(cloudapiv6.FlagIPv6CidrBlock, "", "disable", cloudapiv6.FlagIPv6CidrBlockDescriptionForNIC)
	cmd.AddBoolFlag(cloudapiv6.FlagDHCPv6, "", true, "Whether the NIC reserves an IPv6 address automatically via DHCPv6. Only applies when the target LAN has IPv6 enabled. Set --dhcpv6=false to disable. Default: true")
	cmd.AddStringSliceFlag(cloudapiv6.FlagIPv6IPs, "", nil, "One or more IPv6 IPs to assign to the NIC. Each must fall within the NIC's IPv6 CIDR block (--ipv6-cidr-block), and the target LAN must have IPv6 enabled")

	return cmd
}
