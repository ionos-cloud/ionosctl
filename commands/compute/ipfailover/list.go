package ipfailover

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IpFailoverListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "ipfailover",
		Resource:   "ipfailover",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IP Failovers groups from a LAN",
		LongDesc:   "Use this command to list the IP failover entries configured on a LAN. Each entry shows a floating IP and the NIC that is registered for it (NicId column). NICs sharing the same IP belong to the same failover group.\n\nRequired values to run command:\n\n* Data Center Id\n* Lan Id",
		Example:    `ionosctl compute ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID`,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunIpFailoverList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The unique ID of the Virtual Data Center that holds the LAN", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", "The unique ID of the LAN whose IP failover entries to list", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
