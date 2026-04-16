package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "snapshot",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List the snapshots of your Mongo Cluster",
		Example:   "ionosctl dbaas mongo cluster snapshot ls --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Verbose("Getting snapshots of Cluster %s", clusterId)

			snapshots, _, err := client.Must().MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(snapshots)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true

	return cmd
}
