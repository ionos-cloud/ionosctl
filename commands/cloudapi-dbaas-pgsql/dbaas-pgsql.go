package cloudapi_dbaas_pgsql

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/spf13/cobra"
)

func DBaaSPgsqlCmd() *core.Command {
	dbaasPgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dbaas-pgsql",
			Short:            "DBaaS Pgsql Operations",
			Long:             "The sub-commands of `ionosctl dbaas-pgsql` allow you to perform operations on DBaaS Pgsql.",
			TraverseChildren: true,
		},
	}
	dbaasPgsqlCmd.AddCommand(ClusterCmd())
	return dbaasPgsqlCmd
}
