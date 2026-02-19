package volume

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VolumeGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "volume",
		Resource:   "volume",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Volume",
		LongDesc:   "Use this command to retrieve information about a Volume using its ID.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id",
		Example:    `ionosctl compute volume get --datacenter-id DATACENTER_ID --volume-id VOLUME_ID`,
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunVolumeGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
