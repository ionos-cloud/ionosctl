package backup

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			rows, err := resource2table.ConvertDbaasMariadbBackupsToTable(backups)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(backups, rows,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
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
