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
		Example:   "ionosctl dbaas mariadb backup get --backup-id BACKUP_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupId)
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

	cmd.AddStringFlag(constants.FlagBackupId, "", "", "The ID of the Backup to be retrieved",
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
