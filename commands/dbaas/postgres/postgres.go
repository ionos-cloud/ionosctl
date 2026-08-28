package postgres

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/database"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/user"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func DBaaSPostgresCmd() *core.Command {
	pgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "postgres",
			Aliases: []string{"pg", "pgsql", "postgresql", "psql"},
			Short:   "Manage DBaaS PostgreSQL clusters, databases, users and backups",
			Long: `Manage IONOS DBaaS PostgreSQL, a fully managed PostgreSQL service.

A ` + "`cluster`" + ` is the top-level resource: one or more PostgreSQL instances (one master plus n-1 read-standbys) provisioned in a physical location, reachable over a private LAN in one of your virtual datacenters. Each cluster holds logical databases owned by users (roles), and supports automated point-in-time backups used for restore and cloning.`,
			TraverseChildren: true,
		},
	}
	pgsqlCmd.AddCommand(ClusterCmd())
	pgsqlCmd.AddCommand(LogsCmd())
	pgsqlCmd.AddCommand(BackupCmd())
	pgsqlCmd.AddCommand(PgsqlVersionCmd())
	pgsqlCmd.AddCommand(APIVersionCmd())
	pgsqlCmd.AddCommand(user.UserCmd())
	pgsqlCmd.AddCommand(database.DatabaseCmd())
	return core.WithConfigOverride(pgsqlCmd, []string{fileconfiguration.PSQL}, constants.DefaultApiURL+"/databases/postgresql")
}
