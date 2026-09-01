package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func LanListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List LANs",
		LongDesc:  "Use this command to list the LANs provisioned in a Virtual Data Center. The output shows each LAN's ID, name, whether it is public, the Cross-Connect (PCC) it is attached to, its IPv6 CIDR block, and its provisioning state.\n\nUse `--all` to list LANs across every datacenter in the account instead of a single one (in which case `--datacenter-id` is not required).\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LANsFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example: `# List the LANs in one datacenter
ionosctl compute lan list --datacenter-id DATACENTER_ID

# List LANs across all datacenters in the account
ionosctl compute lan list --all`,
		PreCmdRun:  PreRunLansList,
		CmdRun:     RunLanList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	return cmd
}
