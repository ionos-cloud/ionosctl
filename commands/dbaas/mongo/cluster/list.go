package cluster

import (
	"context"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ClusterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
			var limitPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
				limit := viper.GetInt32(f)
				limitPtr = &limit
			}
			var offsetPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
				offset := viper.GetInt32(f)
				offsetPtr = &offset
			}
			clusters, _, err := c.DbaasMongoServices.Clusters().List(name, limitPtr, offsetPtr)
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, clusters.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the Mongo Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
