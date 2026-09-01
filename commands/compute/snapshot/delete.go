package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func SnapshotDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete a Snapshot",
		LongDesc:   "Use this command to permanently delete a Snapshot. This removes only the stored image; Volumes previously created or restored from it are unaffected, and deleting the source Volume does not delete its Snapshots. If the Snapshot was created with --sec-auth-protection, deletion requires the Contract Owner or a re-authenticated user.\n\nUse `--all` to delete every Snapshot in the contract. Use `--wait` (`-w`) to block until deletion completes.\n\nRequired values to run command:\n\n* Snapshot Id",
		Example:    "ionosctl compute snapshot delete --snapshot-id SNAPSHOT_ID --wait",
		PreCmdRun:  PreRunSnapshotDelete,
		CmdRun:     RunSnapshotDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every Snapshot in the contract instead of a single one. Mutually exclusive with --snapshot-id")

	return cmd
}
