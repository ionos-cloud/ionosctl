package firewallrule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FirewallRuleCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "firewallrule",
		Resource:  "firewallrule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Firewall Rule",
		LongDesc: `Add a new Firewall Rule to the NIC identified by --datacenter-id / --server-id / --nic-id. Every Firewall Rule belongs to exactly one NIC.

A rule WHITELISTS a slice of traffic: while the NIC's firewall is active, traffic is only allowed if a rule matches it (default-deny). --direction (alias --type) selects INGRESS (traffic entering the NIC) or EGRESS (traffic leaving the NIC); it defaults to INGRESS.

--protocol determines which other match flags apply:
  * TCP / UDP  -> --port-range-start and --port-range-end restrict the destination port range. Leave both unset to allow all ports.
  * ICMP       -> --icmp-type and --icmp-code restrict the ICMP message. Leave unset to allow all types/codes.
  * ANY        -> matches every protocol; port and ICMP flags do not apply.
--source-mac, --source-ip and --destination-ip narrow the match further for any protocol; any of them left unset acts as a wildcard (allow all).

NOTE: --protocol is fixed at creation time. It cannot be changed later (the update command has no --protocol flag); to change protocol you must delete the rule and create a new one.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Protocol`,
		Example: `# Allow inbound SSH (TCP port 22) from any source
ionosctl compute firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol TCP --direction INGRESS --port-range-start 22 --port-range-end 22 --name "Allow SSH"

# Allow inbound ICMP echo-request (ping, type 8) from a single source IP only
ionosctl compute firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol ICMP --direction INGRESS --icmp-type 8 --source-ip 192.0.2.10 --name "Allow ping from admin host"`,
		PreCmdRun:  PreRunDcServerNicIdsFRuleProtocol,
		CmdRun:     RunFirewallRuleCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Rule", "A human-friendly label for the rule. Has no effect on matching; used only to identify the rule in listings")
	cmd.AddStringFlag(cloudapiv6.ArgProtocol, "", "", "The IP protocol this rule matches. TCP/UDP also honour --port-range-start/--port-range-end; ICMP also honours --icmp-type/--icmp-code; ANY matches every protocol. Fixed at creation - it cannot be changed by a later update", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"TCP", "UDP", "ICMP", "ANY"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgSourceMac, "", "", "Match only traffic originating from this MAC address. Format: aa:bb:cc:dd:ee:ff. Leave unset to allow any source MAC")
	cmd.AddIpFlag(cloudapiv6.ArgSourceIp, "", nil, "Match only traffic originating from this IP address or CIDR block (must match --ip-version). Leave unset to allow any source IP")
	cmd.AddIpFlag(cloudapiv6.ArgDestinationIp, "", nil, "When the NIC has multiple IPs, match only traffic directed to this IP address or CIDR block of the NIC (must match --ip-version). Leave unset to allow any target IP. WARNING: the short-hand flag `-D` is deprecated")
	cmd.AddIntFlag(cloudapiv6.ArgIcmpType, "", 0, "Only when --protocol ICMP: match this ICMP type (0-254), e.g. 8 = echo request (ping), 0 = echo reply. Leave unset to allow all types")
	cmd.AddIntFlag(cloudapiv6.ArgIcmpCode, "", 0, "Only when --protocol ICMP: match this ICMP code (0-254). Leave unset to allow all codes")
	cmd.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Only when --protocol TCP or UDP: first port of the allowed destination-port range (1-65534, inclusive). Set both --port-range-start and --port-range-end (use the same value for a single port); leave both unset to allow all ports")
	cmd.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Only when --protocol TCP or UDP: last port of the allowed destination-port range (1-65534, inclusive). Set both --port-range-start and --port-range-end; leave both unset to allow all ports")
	cmd.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "INGRESS", "Direction of traffic the rule matches: INGRESS (entering the NIC) or EGRESS (leaving the NIC). Defaults to INGRESS")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
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
	cmd.AddSetFlag(cloudapiv6.FlagIPVersion, "", "IPv4", []string{"IPv4", "IPv6"}, "The IP version this rule applies to. If --source-ip/--destination-ip are given it must match their version; if omitted it is deduced from those addresses. With no IPs given the rule only allows the selected version (defaults to IPv4)")

	return cmd
}
