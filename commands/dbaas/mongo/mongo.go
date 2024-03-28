package mongo

import (
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/apiversion"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/user"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func DBaaSMongoCmd() *core.Command {
	deprecatedAliases := []string{"mongodb", "mg"}

	mongoCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "mongo",
			Aliases: append(deprecatedAliases, "mdb", "m"),
			Short:   "DBaaS Mongo Operations",
			Long: `DBaaS Mongo Operations. Wiki: https://docs.ionos.com/cloud/managed-services/database-as-a-service/mongodb
With IONOS Cloud Database as a Service (DBaaS) MongoDB, you can quickly set up and manage MongoDB database clusters. It is an open-source, NoSQL database solution that does not require a relational Database Management System (DBMS). The feature offers flexible data schemas, managed MongoDB solution with deployment and monitoring of your databases. To cater to your workload use cases, IONOS provides MongoDB editions such as Playground, Business, and Enterprise models.`,
			TraverseChildren: true,
		},
	}

	mongoCmd.Command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		aliasUsed := os.Args[2]
		if slices.Contains(deprecatedAliases, aliasUsed) {
			cmd.PrintErrf("WARNING: '%s' is deprecated and will be removed in a future release, "+
				"please use any of '%s' instead.\n",
				aliasUsed,
				strings.Join(append([]string{mongoCmd.Command.Use}, deprecatedAliases...), "', '"),
			)
		}
	}

	mongoCmd.AddCommand(cluster.ClusterCmd())
	mongoCmd.AddCommand(templates.TemplatesCmd())
	mongoCmd.AddCommand(user.UserCmd())
	mongoCmd.AddCommand(snapshot.SnapshotCmd())
	mongoCmd.AddCommand(logs.LogsCmd())
	mongoCmd.AddCommand(apiversion.ApiVersionCmd())
	return mongoCmd
}
