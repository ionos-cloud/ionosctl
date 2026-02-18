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

func BackupUnitUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a BackupUnit",
		LongDesc: `Use this command to update details about a specific BackupUnit. The password and the email may be updated.

Required values to run command:

* BackupUnit Id`,
		Example:    `ionosctl backupunit update --backupunit-id BACKUPUNIT_ID --email EMAIL`,
		PreCmdRun:  cloudapiv6cmds.PreRunBackupUnitId,
		CmdRun:     cloudapiv6cmds.RunBackupUnitUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Alphanumeric password you want to update for the BackupUnit")
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The e-mail address you want to update for the BackupUnit")
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait,
		"Wait for the Request for BackupUnit update to be executed")
	cmd.Command.Flags().MarkHidden(constants.ArgWaitForRequest) // Backupunit resources are not tracked by /requests endpoint yet - but keep the flag for backward compatibility
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit update [seconds]")

	return cmd
}
