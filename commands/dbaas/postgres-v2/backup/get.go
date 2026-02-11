package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
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
		Example:    "ionosctl dbaas postgres backup get --backup-id <backup-id>",
		PreCmdRun:  PreRunBackupId,
		CmdRun:     RunBackupGet,
		InitClient: true,
	})
	get.AddUUIDFlag(constants.FlagBackupId, constants.FlagIdShort, "", "The unique ID of the Backup", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(constants.ArgCols, "", defaultBackupCols, tabheaders.ColsMessage(allBackupCols))

	return get
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupId)
}

func RunBackupGet(c *core.CommandConfig) error {
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.FlagBackupId, backupId))

	backup, _, err := client.Must().PostgresClientV2.BackupsApi.BackupsFindById(context.Background(), backupId).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresV2Backup, backup,
		tabheaders.GetHeaders(allBackupCols, defaultBackupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
