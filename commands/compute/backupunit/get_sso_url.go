package backupunit

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func BackupUnitGetSsoUrlCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "get-sso-url",
		ShortDesc:  "Get BackupUnit SSO URL",
		LongDesc:   "Use this command to obtain a Single Sign-On (SSO) URL for a BackupUnit. Copy the entire returned URL and paste it into a browser: you will be logged straight into the backup web console (https://backup.ionos.com) for that unit, without entering the login/password manually.\n\nThe URL is short-lived and grants access to manage backup plans, schedules and restores for this unit.\n\nRequired values to run command:\n\n* BackupUnit Id",
		Example:    `ionosctl compute backupunit get-sso-url --backupunit-id BACKUPUNIT_ID`,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitGetSsoUrl,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
