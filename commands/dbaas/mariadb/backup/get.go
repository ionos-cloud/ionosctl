package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List MariaDB Backups",
		LongDesc:  "List all MariaDB Backups, or optionally provide a Cluster ID to list those of a certain cluster",
		Example:   "ionosctl dbaas mariadb backup get --backup-id BACKUP_ID",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			backup, _, err := client.Must().MariaClient.BackupsApi.BackupsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))).Execute()

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			rows, err := resource2table.ConvertDbaasMariadbBackupToTable(backup)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(backup, rows,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagBackupId, "", "", "The ID of the Backup to be retrieved")

	return cmd
}
