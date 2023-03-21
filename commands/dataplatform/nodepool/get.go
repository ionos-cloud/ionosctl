package nodepool

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get Dataplatform Nodepool by cluster and nodepool id",
		Example:   "ionosctl dataplatform nodepool get",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagNodepoolId)
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			npId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
			c.Printer.Verbose("Getting Nodepool %s...", npId)

			client, err := client2.Get()
			if err != nil {
				return err
			}

			np, _, err := client.DataplatformClient.DataPlatformNodePoolApi.GetClusterNodepool(c.Context, clusterId, npId).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getNodepoolsPrint(c, &[]ionoscloud.NodePoolResponseData{np}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, "", "", "The unique ID of the cluster. Must conform to the UUID format")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagNodepoolId, constants.FlagIdShort, "", "The unique ID of the nodepool. Must conform to the UUID format")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformNodepoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveError
	})
	cmd.Command.SilenceUsage = true

	return cmd
}
