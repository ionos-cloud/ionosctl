package dataplatform

import (
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/postgres"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func DataplatformCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dataplatform",
			Short:            "Managed Stackable Data Platform by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable Platform.",
			Aliases:          []string{"mdp", "dp", "stackable", "managed-dataplatform"},
			Long:             "The sub-commands of `ionosctl dataplatform` allow you to perform operations on DBaaS resources.",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(postgres.DBaaSPostgresCmd())
	cmd.AddCommand(mongo.DBaaSMongoCmd())
	return cmd
}
