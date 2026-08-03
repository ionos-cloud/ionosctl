package firewallrule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FirewallRuleUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "firewallrule",
		Resource:  "firewallrule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a FirewallRule",
		LongDesc: `Update the matching criteria of an existing Firewall Rule on a NIC. Only the flags you pass are changed; the rest keep their current values.

You can retune the match: --direction, --source-mac, --source-ip, --destination-ip, plus --port-range-start/--port-range-end (for a TCP/UDP rule) or --icmp-type/--icmp-code (for an ICMP rule), and --name.

NOTE: the rule's protocol is fixed at creation and CANNOT be changed here (there is no --protocol flag). To change protocol, delete the rule and create a new one. Editing port flags on an ICMP rule (or ICMP flags on a TCP/UDP rule) has no effect.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id`,
		Example: `# Rename a rule
ionosctl compute firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --name "New name" --wait

# Widen an existing TCP rule to the HTTPS port range and restrict it to one source IP
ionosctl compute firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --port-range-start 443 --port-range-end 443 --source-ip 192.0.2.0/24 --wait`,
		PreCmdRun:  PreRunDcServerNicFRuleIds,
		CmdRun:     RunFirewallRuleUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A human-friendly label for the rule. Has no effect on matching; used only to identify the rule in listings")
	cmd.AddStringFlag(cloudapiv6.ArgSourceMac, "", "", "Match only traffic originating from this MAC address. Format: aa:bb:cc:dd:ee:ff. Leave unset to allow any source MAC")
	cmd.AddIpFlag(cloudapiv6.ArgSourceIp, "", nil, "Match only traffic originating from this IP address or CIDR block (must match --ip-version). Leave unset to allow any source IP")
	cmd.AddIpFlag(cloudapiv6.ArgDestinationIp, "", nil, "When the NIC has multiple IPs, match only traffic directed to this IP address or CIDR block of the NIC (must match --ip-version). Leave unset to allow any target IP. WARNING: the short-hand flag `-D` is deprecated")
	cmd.AddIntFlag(cloudapiv6.ArgIcmpType, "", 0, "Only for an ICMP rule: match this ICMP type (0-254), e.g. 8 = echo request (ping), 0 = echo reply. Leave unset to allow all types. Has no effect on a TCP/UDP/ANY rule")
	cmd.AddIntFlag(cloudapiv6.ArgIcmpCode, "", 0, "Only for an ICMP rule: match this ICMP code (0-254). Leave unset to allow all codes. Has no effect on a TCP/UDP/ANY rule")
	cmd.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Only for a TCP/UDP rule: first port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports. Has no effect on an ICMP/ANY rule")
	cmd.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Only for a TCP/UDP rule: last port of the allowed destination-port range (1-65535, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports. Has no effect on an ICMP/ANY rule")
	cmd.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "", "Direction of traffic the rule matches: INGRESS (entering the NIC) or EGRESS (leaving the NIC)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgFirewallRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.FlagIPVersion, "", "IPv4", []string{"IPv4", "IPv6"}, "The IP version this rule applies to. If --source-ip/--destination-ip are given it must match their version; if omitted it is deduced from those addresses")

	return cmd
}
