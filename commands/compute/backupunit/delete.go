package backupunit

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
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
		LongDesc: `Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id`,
		Example:    `ionosctl backupunit delete --backupunit-id BACKUPUNIT_ID`,
		PreCmdRun:  cloudapiv6cmds.PreRunBackupUnitDelete,
		CmdRun:     cloudapiv6cmds.RunBackupUnitDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait,
		"Wait for the Request for BackupUnit deletion to be executed")
	cmd.Command.Flags().MarkHidden(constants.ArgWaitForRequest) // Backupunit resources are not tracked by /requests endpoint yet - but keep the flag for backward compatibility
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all BackupUnits.")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit deletion [seconds]")

	return cmd
}
