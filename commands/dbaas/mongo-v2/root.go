package mongo_v2

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/template"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "mongo-v2",
			Aliases: []string{"mongodb-v2", "mg-v2"},
			Short:   "DBaaS MongoDB V2 Operations",
			Long: "The sub-commands of `ionosctl dbaas mongo-v2` allow you to perform operations on MongoDB V2 resources.\n\n" +
				"The MongoDB V2 API only accepts token authentication: use `IONOS_TOKEN` or a token from `ionosctl cfg login`, not a username and password.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.ClusterCmd())
	cmd.AddCommand(snapshot.SnapshotCmd())
	cmd.AddCommand(template.TemplateCmd())
	cmd.AddCommand(version.VersionCmd())

	return core.WithRegionalConfigOverride(cmd, []string{constants.FileConfigMongoDBV2}, constants.MongoDBApiRegionalURL, constants.MongoDBLocations)
}
