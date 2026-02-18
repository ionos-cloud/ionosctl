package nic

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

func NicCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NIC",
		LongDesc: `Use this command to create/add a new NIC to the target Server. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run a command:

* Data Center Id
* Server Id`,
		Example:    `ionosctl nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name NAME`,
		PreCmdRun:  cloudapiv6cmds.PreRunNicCreate,
		CmdRun:     cloudapiv6cmds.RunNicCreate,
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
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Internet Access", "The name of the NIC")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "IPs assigned to the NIC. This can be a collection")
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", cloudapiv6.DefaultDhcp, "Set to false if you wish to disable DHCP on the NIC. E.g.: --dhcp=true, --dhcp=false")
	cmd.AddBoolFlag(cloudapiv6.ArgFirewallActive, "", cloudapiv6.DefaultFirewallActive, "Activate or deactivate the Firewall. E.g.: --firewall-active=true, --firewall-active=false")
	cmd.AddStringFlag(cloudapiv6.ArgFirewallType, "", "INGRESS", "The type of Firewall Rules that will be allowed on the NIC")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgLanId, "", cloudapiv6.DefaultNicLanId, "The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NIC creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NIC creation [seconds]")

	cmd.AddStringFlag(cloudapiv6.FlagIPv6CidrBlock, "", "disable", cloudapiv6.FlagIPv6CidrBlockDescriptionForNIC)
	cmd.AddBoolFlag(cloudapiv6.FlagDHCPv6, "", true, "Set to false if you wish to disable DHCPv6 on the NIC. E.g.: --dhcpv6=true, --dhcpv6=false")
	cmd.AddStringSliceFlag(cloudapiv6.FlagIPv6IPs, "", nil, "IPv6 IPs assigned to the NIC. They need to be within the NIC IPv6 Cidr Block.")

	return cmd
}
