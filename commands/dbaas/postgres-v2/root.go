package postgres_v2

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	pgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "postgres-v2",
			Aliases:          []string{"pg-v2"},
			Short:            "DBaaS PostgreSQL V2 Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres-v2` allow you to perform operations on DBaaS PostgreSQL V2 resources.",
			TraverseChildren: true,
		},
	}
	pgsqlCmd.AddCommand(cluster.ClusterCmd())
	pgsqlCmd.AddCommand(version.VersionCmd())
	pgsqlCmd.AddCommand(backup.BackupCmd())

	// todo: decide config override 'name' key
	return core.WithRegionalConfigOverride(pgsqlCmd, []string{"todo"}, constants.PostgresApiRegionalURL, constants.PostgresLocations)
}
