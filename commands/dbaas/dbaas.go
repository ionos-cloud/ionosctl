package dbaas

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func DataBaseServiceCmd() *core.Command {
	dbaasCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dbaas",
			Short:            "Database as a Service Operations",
			Aliases:          []string{"db"},
			Long:             "The sub-commands of `ionosctl dbaas` allow you to perform operations on DBaaS resources.",
			TraverseChildren: true,
		},
	}
	dbaasCmd.AddCommand(postgres.DBaaSPostgresCmd())
	dbaasCmd.AddCommand(mongo.DBaaSMongoCmd())
	dbaasCmd.AddCommand(mariadb.Root())
	dbaasCmd.AddCommand(inmemorydb.Root())
	return dbaasCmd
}
