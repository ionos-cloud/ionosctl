package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List MariaDB Backups",
		LongDesc:  "List all MariaDB Backups, or optionally provide a Cluster ID to list those of a certain cluster",
		Example:   "ionosctl dbaas mariadb backup list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			var backups mariadb.BackupList
			var err error

			if clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)); clusterId != "" {
				backups, _, err = client.Must().MariaClient.BackupsApi.ClusterBackupsGet(context.Background(), clusterId).Execute()
			} else {
				backups, err = Backups()
			}

			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(backups)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "Optionally limit shown backups to those of a certain cluster",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return cluster.ClustersProperty(func(c mariadb.ClusterResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id
				})
			}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)

	return cmd
}
