package nic

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NicUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a NIC",
		LongDesc: `Use this command to update the configuration of an existing NIC, identified within its Data Center (--datacenter-id) and Server (--server-id) by --nic-id. Only the flags you set are changed; everything else is left as-is.

Common updates:
* Move the NIC to a different network by changing --lan-id.
* Add reserved public IPs or assign private IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) with --ips.
* Toggle DHCP with --dhcp.
* Enable/disable the per-NIC firewall with --firewall-active and set its direction with --firewall-type (INGRESS/EGRESS/BIDIRECTIONAL). When enabled, incoming traffic is filtered by the NIC's firewall rules; when disabled, all traffic reaches the NIC directly and the rules are ignored.

Restriction: the primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the NIC to reach AVAILABLE state before the command returns.

Required values to run command:

* Data Center Id
* NIC Id`,
		Example: `# Move a NIC to a different LAN
ionosctl compute nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id 2 --wait

# Rename a NIC, add a reserved public IP, and enable a bidirectional firewall
ionosctl compute nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name web-nic --ips 203.0.113.10 --firewall-active=true --firewall-type BIDIRECTIONAL --wait`,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunNicUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, cloudapiv6.ArgIdShort, "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A human-friendly name for the NIC (shown in the DCD and listings). Does not affect networking")
	cmd.AddIntFlag(cloudapiv6.ArgLanId, "", cloudapiv6.DefaultNicLanId, "Move the NIC to this LAN ID, changing which network the server reaches through this NIC. If the LAN ID does not exist in the Data Center, it is created implicitly")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgFirewallActive, "", cloudapiv6.DefaultFirewallActive, "Enable the per-NIC firewall. When enabled, an empty ruleset blocks all incoming traffic and only explicitly configured firewall rules are allowed through; when disabled, all traffic reaches the NIC and rules are ignored")
	cmd.AddStringFlag(cloudapiv6.ArgFirewallType, "", "INGRESS", "Direction of traffic the NIC's firewall rules apply to. INGRESS = inbound only, EGRESS = outbound only, BIDIRECTIONAL = both. Only meaningful when --firewall-active=true")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", cloudapiv6.DefaultDhcp, "Whether the NIC reserves an IP automatically from the LAN's DHCP server (true) or not (false). Set --dhcp=false when managing addresses via --ips")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "IPs to assign to the NIC. Explicitly assigned public IPs must come from a reserved IP block; private-range IPs (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) may be used on private LANs")

	cmd.AddStringFlag(cloudapiv6.FlagIPv6CidrBlock, "", "disable", cloudapiv6.FlagIPv6CidrBlockDescriptionForNIC)
	cmd.AddBoolFlag(cloudapiv6.FlagDHCPv6, "", true, "Whether the NIC reserves an IPv6 address automatically via DHCPv6. Only applies when the target LAN has IPv6 enabled. Set --dhcpv6=false to disable")
	cmd.AddStringSliceFlag(cloudapiv6.FlagIPv6IPs, "", nil, "One or more IPv6 IPs to assign to the NIC. Each must fall within the NIC's IPv6 CIDR block, and the target LAN must have IPv6 enabled")

	return cmd
}
