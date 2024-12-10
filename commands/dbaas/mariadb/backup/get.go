package backup

import (
	"context"
	"fmt"

	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mariadb"

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
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a MariaDB Backup",
		Example:   "ionosctl dbaas mariadb backup get --backup-id BACKUP_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupId)
		},
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

	cmd.AddStringFlag(constants.FlagBackupId, "", "", "The ID of the Backup to be retrieved",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return BackupsProperty(func(c ionoscloud.BackupResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id + "\t" + fmt.Sprintf("(%d MiB)", *c.Properties.Size)
				})
			}, constants.MariaDBApiRegionalURL),
	)
	return cmd
}
