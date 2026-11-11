package mariadb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	return newRoot("mariadb", []string{"maria", "mar", "ma"})
}

// HiddenV2Root keeps the pre-promotion `mariadb-v2` tree reachable so existing scripts keep working.
func HiddenV2Root() *core.Command {
	cmd := newRoot("mariadb-v2", []string{"maria-v2", "mar-v2", "ma-v2"})
	cmd.Command.Hidden = true
	return cmd
}

func newRoot(use string, aliases []string) *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     use,
			Aliases: aliases,
			Short:   "DBaaS MariaDB Operations",
			Long: `Manage IONOS Database as a Service (DBaaS) MariaDB.

MariaDB is an open-source relational database, MySQL-compatible. This command tree lets you provision and operate fully managed MariaDB clusters - IONOS handles the underlying instances, replication, patching (during a weekly maintenance window) and the automatic backups used for point-in-time restore, so you only manage the cluster's shape and data.

MariaDB is a regional service, so commands operate against a specific location; set it with --location (or the config) where required.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.ClusterCmd())
	cmd.AddCommand(backup.BackupCmd())
	cmd.AddCommand(version.VersionCmd())

	// Config key stays `mariadbv2`: the v1 `mariadb` key points at the deactivated v1 API.
	return core.WithRegionalConfigOverride(cmd, []string{constants.FileConfigMariaDBV2}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations)
}
