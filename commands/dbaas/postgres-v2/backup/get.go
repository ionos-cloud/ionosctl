package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func BackupGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "backup",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Backup",
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Backup by using its ID.\n\nRequired values to run command:\n\n* Backup Id",
		Example:    "ionosctl dbaas postgres-v2 backup get --backup-id <backup-id>",
		PreCmdRun:  PreRunBackupId,
		CmdRun:     RunBackupGet,
		InitClient: true,
	})
	get.AddUUIDFlag(constants.FlagBackupId, constants.FlagIdShort, "", "The unique ID of the Backup", core.RequiredFlagOption(),
		core.WithCompletion(completer.BackupIds, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	get.AddStringSliceFlag(constants.ArgCols, "", defaultBackupCols, table.ColsMessage(backupCols))

	return get
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupId)
}

func RunBackupGet(c *core.CommandConfig) error {
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))

	c.Verbose(constants.FlagBackupId, backupId)

	backup, _, err := client.Must().PostgresClientV2.BackupsApi.BackupsFindById(context.Background(), backupId).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))
	return c.Out(table.Sprint(backupCols, backup, cols))
}
