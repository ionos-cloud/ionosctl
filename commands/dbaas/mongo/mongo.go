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
			Use:              "mongo",
			Aliases:          []string{"mongodb", "mdb", "m"},
			Short:            "DBaaS Mongo Operations",
			Long:             "The sub-commands of `ionosctl dbaas mongo` allow you to perform operations on DBaaS Mongo resources.",
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

const (
	flagClusterId = "cluster-id"
)
