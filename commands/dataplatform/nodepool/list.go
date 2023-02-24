package nodepool

import (
	"context"

	"github.com/ionos-cloud/ionosctl/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Dataplatform Nodepools of a certain cluster",
		Example:   "ionosctl dataplatform nodepool list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Nodepools...")
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			client, err := config.GetClient()
			if err != nil {
				return err
			}

			np, _, err := client.DataplatformClient.DataPlatformNodePoolApi.GetClusterNodepools(c.Context, clusterId).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getNodepoolsPrint(c, np.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster. Must conform to the UUID format")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true

	return cmd
}
