package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a Mongo Cluster by ID",
		Example:   "ionosctl dbaas mongo cluster get --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Cluster by id: %s", clusterId))

			cluster, _, err := client.Must().MongoClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			clusterConverted, err := jsonpaths.ConvertClusterToTable(cluster)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(cluster, clusterConverted,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
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
