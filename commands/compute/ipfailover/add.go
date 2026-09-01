package ipfailover

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IpFailoverAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add IP Failover group to a LAN",
		LongDesc: `Use this command to add a NIC to an IP failover group on a LAN, registering it with the group's floating public IP.

Setting up a working IP failover group takes three steps:

  1. Reserve a public IP block in the same region/location as the datacenter (` + "`" + `ionosctl ipblock create` + "`" + `) and assign one of its IPs to the NIC that will become the failover MASTER.
  2. Run this command with that IP and the master NIC's Id to create/enable the failover group.
  3. Assign the SAME reserved IP to the other NICs on the same LAN and run this command again for each of them. Those NICs join the group as MEMBERS (standby).

If the group does not exist yet on the LAN, the first ` + "`" + `add` + "`" + ` creates it; subsequent ` + "`" + `add` + "`" + ` calls with the same --ip extend it. The IP must belong to a reserved IP block, not an ad-hoc/DHCP address.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address`,
		Example: `# Register the master NIC with the group's floating IP
ionosctl compute ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID --lan-id LAN_ID --nic-id MASTER_NIC_ID --ip "203.0.113.10"

# Add a standby member NIC (on another server) to the same group, reusing the same IP
ionosctl compute ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID_2 --lan-id LAN_ID --nic-id MEMBER_NIC_ID --ip "203.0.113.10"`,
		PreCmdRun:  PreRunDcLanServerNicIdsIp,
		CmdRun:     RunIpFailoverAdd,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The unique ID of the Virtual Data Center that holds the LAN and servers", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", "The unique ID of the LAN the failover group lives on. All member NICs must be on this LAN", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", "The unique ID of the server that owns the NIC being added to the group", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", "The unique ID of the NIC to add to the failover group. The first NIC added becomes the master; later NICs (with the same IP) become standby members", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The floating public IP for the failover group. Must be an address from a reserved IP block in the same region/location as the datacenter. Reuse the exact same IP when adding member NICs", core.RequiredFlagOption())

	return cmd
}
