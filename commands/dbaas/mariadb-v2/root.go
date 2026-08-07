package mariadb_v2

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "mariadb-v2",
			Aliases:          []string{"maria-v2", "mar-v2", "ma-v2"},
			Short:            "DBaaS MariaDB V2 Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb-v2` allow you to perform operations on MariaDB V2 resources.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.ClusterCmd())
	cmd.AddCommand(backup.BackupCmd())
	cmd.AddCommand(version.VersionCmd())

	// Use a v2-specific config key so a `cfg login` config can hold both the v1
	// (`mariadb`) and v2 (`mariadbv2`) endpoint overrides without one clobbering
	// the other. Same convention as postgres-v2 (psqlv2).
	return core.WithRegionalConfigOverride(cmd, []string{constants.FileConfigMariaDBV2}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations)
}
