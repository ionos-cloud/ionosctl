package postgres

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
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
	return pgsqlCmd
}
