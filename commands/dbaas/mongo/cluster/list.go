package cluster

import (
	"context"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
)

func ClusterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Mongo Clusters",
		LongDesc:  "Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:   "ionosctl dbaas mongo cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Clusters...")
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			clusters, _, err := c.DbaasMongoServices.Clusters().List(name)
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, clusters.GetItems()))
		},
		InitClient: true,
	})

	// TODO: Move ArgName to DBAAS level constants
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the Mongo Clusters that contain the specified name in the DisplayName field. The value is case insensitive")

	return cmd
}
