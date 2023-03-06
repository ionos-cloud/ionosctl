package dbaas

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func DataBaseServiceCmd() *core.Command {
	dbaasCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dbaas",
			Short:            "Database as a Service Operations",
			Long:             "The sub-commands of `ionosctl dbaas` allow you to perform operations on DBaaS resources.",
			TraverseChildren: true,
		},
	}
	dbaasCmd.AddCommand(postgres.DBaaSPostgresCmd())
	return dbaasCmd
}
