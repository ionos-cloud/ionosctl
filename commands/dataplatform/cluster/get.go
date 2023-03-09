package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	sdkdataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a Dataplatform Cluster by ID",
		Example:   "ionosctl dataplatform cluster get --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Printer.Verbose("Getting Cluster by id: %s", clusterId)

			client, err := config.GetClient()
			if err != nil {
				return err
			}

			cluster, _, err := client.DataplatformClient.DataPlatformClusterApi.GetCluster(c.Context, clusterId).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, &[]sdkdataplatform.ClusterResponseData{cluster}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
