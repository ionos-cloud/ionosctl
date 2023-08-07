package mongo

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/apiversion"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/user"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func DBaaSMongoCmd() *core.Command {
	mongoCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "mongo",
			Aliases: []string{"mongodb", "mdb", "m"},
			Short:   "DBaaS Mongo Operations",
			Long: `DBaaS Mongo Operations. Wiki: https://docs.ionos.com/cloud/managed-services/database-as-a-service/mongodb
With IONOS Cloud Database as a Service (DBaaS) MongoDB, you can quickly set up and manage MongoDB database clusters. It is an open-source, NoSQL database solution that does not require a relational Database Management System (DBMS). The feature offers flexible data schemas, managed MongoDB solution with deployment and monitoring of your databases. To cater to your workload use cases, IONOS provides MongoDB editions such as Playground, Business, and Enterprise models.`,
			TraverseChildren: true,
		},
	}
	mongoCmd.AddCommand(cluster.ClusterCmd())
	mongoCmd.AddCommand(templates.TemplatesCmd())
	mongoCmd.AddCommand(user.UserCmd())
	mongoCmd.AddCommand(snapshot.SnapshotCmd())
	mongoCmd.AddCommand(logs.LogsCmd())
	mongoCmd.AddCommand(apiversion.ApiVersionCmd())
	return mongoCmd
}
