package backupunit

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func BackupUnitGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a BackupUnit",
		LongDesc:   "Use this command to retrieve details about a specific BackupUnit.\n\nRequired values to run command:\n\n* BackupUnit Id",
		Example:    `ionosctl compute backupunit get --backupunit-id BACKUPUNIT_ID`,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, cloudapiv6.ArgIdShort, "", cloudapiv6.BackupUnitId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
