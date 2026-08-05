package mariadb_v2

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Mariadb}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations)
}
