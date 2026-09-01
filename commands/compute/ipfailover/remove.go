package ipfailover

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IpFailoverRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove IP Failover group from a LAN",
		LongDesc: `Use this command to remove a NIC from an IP failover group on a LAN. Identify the entry to remove by its --ip and --nic-id. Removing the master NIC leaves the remaining member NICs in the group; removing the last NIC effectively dissolves the group.

Use --all to clear every IP failover entry from the LAN at once.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address`,
		Example: `# Remove a single NIC from a failover group
ionosctl compute ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "203.0.113.10"

# Remove all IP failover entries from a LAN
ionosctl compute ipfailover remove --datacenter-id DATACENTER_ID --lan-id LAN_ID --all`,
		PreCmdRun:  PreRunDcLanServerNicIpRemove,
		CmdRun:     RunIpFailoverRemove,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The unique ID of the Virtual Data Center that holds the LAN and servers", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", "The unique ID of the LAN the failover group lives on", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", "The unique ID of the server that owns the NIC being removed", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", "The unique ID of the NIC to remove from the failover group", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The floating IP of the failover entry to remove. Together with --nic-id this identifies the specific entry", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all IP failover entries from the LAN. Mutually exclusive with --nic-id/--ip")

	return cmd
}
