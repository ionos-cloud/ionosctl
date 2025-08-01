package postgres

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/database"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/user"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func DBaaSPostgresCmd() *core.Command {
	pgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "postgres",
			Aliases:          []string{"pg"},
			Short:            "DBaaS PostgreSQL Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres` allow you to perform operations on DBaaS PostgreSQL resources.",
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
	return core.WithConfigOverride(pgsqlCmd, []string{"psql"}, constants.DefaultApiURL+"/databases/postgresql")
}
