package datacenter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func DatacenterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Virtual Data Center and everything inside it",
		LongDesc: `Delete a Virtual Data Center. This removes the VDC object AND every resource it contains - servers, volumes (and their data), LANs, NICs, firewall rules - in a single cascading operation.

NOTE: This is a highly destructive, irreversible operation. Deleted volumes cannot be recovered unless you have snapshots or backups stored elsewhere. Use with extreme caution.

You must identify the VDC either with ` + "`--datacenter-id`" + ` (delete one) or with ` + "`--all`" + ` (delete every VDC on the account) - the two are mutually exclusive. Combine with ` + "`--force`" + ` (` + "`-f`" + `) to skip the interactive confirmation prompt (useful in scripts) and ` + "`--wait`" + ` (` + "`-w`" + `) to block until deletion completes.

Required values to run command:

* Data Center Id (or --all)`,
		Example: `# Delete a single VDC (prompts for confirmation)
ionosctl compute datacenter delete --datacenter-id DATACENTER_ID

# Delete a VDC non-interactively and wait for it to finish
ionosctl compute datacenter delete --datacenter-id DATACENTER_ID --force --wait`,
		PreCmdRun:  PreRunDataCenterDelete,
		CmdRun:     RunDataCenterDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every Virtual Data Center on the account (and all resources inside them). Mutually exclusive with --datacenter-id")

	return cmd
}
