package ipfailover

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

func IpFailoverAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add IP Failover group to a LAN",
		LongDesc: `Use this command to add an IP Failover group to a LAN.

Successfully setting up an IP Failover group requires three steps:

* Add a reserved IP address to a NIC that will become the IP Failover master.
* Use ` + "`" + `ionosctl ipfailover add` + "`" + ` command to enable IP Failover by providing the relevant IP and NIC Id values.
* Add the same reserved IP address to any other NICs that are a member of the same LAN. Those NICs will become IP Failover members.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address`,
		Example:    `ionosctl ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID --lan-id LAN_ID --nic-id NIC_ID --ip "x.x.x.x"`,
		PreCmdRun:  cloudapiv6cmds.PreRunDcLanServerNicIdsIp,
		CmdRun:     cloudapiv6cmds.RunIpFailoverAdd,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP address to be added to IP Failover Group", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IP Failover creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IP Failover creation [seconds]")

	return cmd
}
