package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/spf13/viper"
)

func FindByID() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "get",
			Aliases:   []string{"g"},
			ShortDesc: "Retrieve a cluster",
			Example:   "ionosctl kafka cl get --cluster-id ID",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlags(
					c.Command, c.NS, constants.FlagClusterId, constants.FlagLocation,
				); err != nil {
					return err
				}

				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				clusterID := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
				r, _, err := client.Must().Kafka.ClustersApi.ClustersFindById(
					context.Background(),
					clusterID,
				).Execute()
				if err != nil {
					return err
				}

				cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"", jsonpaths.KafkaCluster, r,
					tabheaders.GetHeadersAllDefault(defaultCols, cols),
				)
				if err != nil {
					return err
				}

				fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
				return nil
			},
			InitClient: true,
		},
	)
	cmd.AddStringFlag(
		constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the cluster you want to retrieve",
		core.RequiredFlagOption(), core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(k kafka.ClusterRead) string {
						return *k.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
