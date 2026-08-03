package backupunit

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func BackupUnitUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a BackupUnit",
		LongDesc: `Use this command to update a BackupUnit. Only the --password and --email may be changed; the --name (backup login) is immutable and cannot be updated here (to rename, delete and recreate the unit).

Changing --password rotates the backup agent login secret. Like on create, the password is never returned by the API, so record any new value you set.

Required values to run command:

* BackupUnit Id`,
		Example: `# Change the report e-mail
ionosctl compute backupunit update --backupunit-id BACKUPUNIT_ID --email newops@example.com

# Rotate the login password
ionosctl compute backupunit update --backupunit-id BACKUPUNIT_ID --password 'N3wS3cret!'`,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "New login secret for the backup agent. Write-only: never returned by the API, so record it")
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "New e-mail address for backup service reports")
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
