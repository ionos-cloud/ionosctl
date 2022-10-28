package mongo

import (
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func DBaaSMongoCmd() *core.Command {
	mongoCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "mongo",
			Aliases:          []string{"mongodb", "mdb", "m"},
			Short:            "DBaaS Mongo Operations",
			Long:             "The sub-commands of `ionosctl dbaas mongo` allow you to perform operations on DBaaS Mongo resources.",
			TraverseChildren: true,
		},
	}
	mongoCmd.AddCommand(cluster.ClusterCmd())
	mongoCmd.AddCommand(templates.TemplatesCmd())
	//mongoCmd.AddCommand(BackupCmd())
	//mongoCmd.AddCommand(PgsqlVersionCmd())
	//mongoCmd.AddCommand(APIVersionCmd())
	return mongoCmd
}