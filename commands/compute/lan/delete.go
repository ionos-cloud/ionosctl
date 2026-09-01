package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LanDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a LAN",
		LongDesc: `Use this command to delete a LAN from a Virtual Data Center. Detach any NICs from the LAN first; deleting a LAN that is attached to a Cross-Connect removes it from that Cross-Connect as well.

Use ` + "`--all`" + ` to delete every LAN in the datacenter in one call. Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the deletion request to complete, and ` + "`" + `--force` + "`" + ` to skip the confirmation prompt.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example: `ionosctl compute lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID

ionosctl compute lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID --wait`,
		PreCmdRun:  PreRunLanDelete,
		CmdRun:     RunLanDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every LAN in the given Virtual Data Center (--datacenter-id). --lan-id is not needed when this is set")

	return cmd
}
