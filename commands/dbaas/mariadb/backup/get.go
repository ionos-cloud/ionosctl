package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a MariaDB Backup",
		LongDesc:  "Retrieve details of a single MariaDB backup by its ID: the cluster it belongs to, the earliest timestamp you can restore to, its total size, and the individual base backups. The recovery window runs from earliestRecoveryTargetTime up to now.",
		Example:   "ionosctl dbaas mariadb backup get --backup-id BACKUP_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagBackupId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			backup, _, err := client.Must().MariaClient.BackupsApi.BackupsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(backup)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagBackupId, "", "", "The unique ID of the backup to retrieve",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return BackupsProperty(func(c mariadb.BackupResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id + "\t" + fmt.Sprintf("(%d MiB)", *c.Properties.Size)
				})
			}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	return cmd
}
