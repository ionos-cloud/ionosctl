package mariadb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "mariadb",
			Aliases: []string{"maria", "mar", "ma"},
			Short:   "DBaaS MariaDB Operations",
			Long: `Manage IONOS Database as a Service (DBaaS) MariaDB.

MariaDB is an open-source relational database, MySQL-compatible. This command tree lets you provision and operate fully managed MariaDB clusters - IONOS handles the underlying instances, replication, patching (during a weekly maintenance window) and continuous backups, so you only manage the cluster's shape and data.

Sub-commands:
  cluster - create, list, inspect, resize, upgrade and delete MariaDB clusters.
  backup  - list and inspect the automatic backups used for point-in-time restore.

MariaDB is a regional service, so commands operate against a specific location; set it with --location (or the config) where required.`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(cluster.Root())
	cmd.AddCommand(backup.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Mariadb}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations)
}
