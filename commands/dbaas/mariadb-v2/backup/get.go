package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func BackupGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "backup",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a MariaDB Backup",
		Example:    "ionosctl dbaas mariadb-v2 backup get --backup-id <backup-id>",
		LongDesc:   "Use this command to retrieve details about a MariaDB Backup by using its ID.\n\nThe backup represents a recovery WINDOW: earliestRecoveryTargetTime bounds the start of the range you can restore to (the window extends to the present). Pass a timestamp inside that range as `--recovery-time` on `cluster restore` / `cluster create --backup-id` to zoom into a specific point; omit it to use the latest.\n\nRequired values to run command:\n\n* Backup Id",
		PreCmdRun:  PreRunBackupId,
		CmdRun:     RunBackupGet,
		InitClient: true,
	})

	get.AddUUIDFlag(constants.FlagBackupId, constants.FlagIdShort, "", "The unique ID of the Backup", core.RequiredFlagOption(),
		core.WithCompletion(completer.BackupIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	return get
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagBackupId)
}

func RunBackupGet(c *core.CommandConfig) error {
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
	c.Verbose("Getting Backup %v...", backupId)

	backup, _, err := client.Must().MariaClientV2.BackupsApi.BackupsFindById(
		context.Background(), backupId).Execute()
	if err != nil {
		return fmt.Errorf("could not get backup: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(backupCols, backup, cols))
}
