package pg

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/spf13/cobra"
)

func DBaaSPgCmd() *core.Command {
	pgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "pg",
			Aliases:          []string{"postgres", "pgsql"},
			Short:            "DBaaS PostgreSQL Operations",
			Long:             "The sub-commands of `ionosctl pg` allow you to perform operations on DBaaS PostgreSQL resources.",
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
