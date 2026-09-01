package backupunit

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func BackupUnitDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a BackupUnit",
		LongDesc: `Use this command to delete a BackupUnit. This is a DESTRUCTIVE and irreversible operation: a successful delete removes the backup plans inside the unit, ALL backups stored in it, the backup login user, and finally the BackupUnit itself.

Because the name (backup login) is immutable, deleting is also the only way to "rename" a unit: delete and recreate under a new name (note the recreated unit starts empty).

Required values to run command:

* BackupUnit Id`,
		Example: `# Delete one BackupUnit
ionosctl compute backupunit delete --backupunit-id BACKUPUNIT_ID

# Delete every BackupUnit under the contract
ionosctl compute backupunit delete --all`,
		PreCmdRun:  PreRunBackupUnitDelete,
		CmdRun:     RunBackupUnitDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all BackupUnits under the contract (each with its backups). Use instead of --backupunit-id")

	return cmd
}
