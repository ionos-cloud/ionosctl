package topic

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func listCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "List all kafka topics",
			Aliases:   []string{"ls"},
			Example: `ionosctl kafka topic list --cluster-id CLUSTER_ID
ionosctl kafka topic list --location LOCATION --cluster-id CLUSTER_ID`,
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return nil
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				if !cmd.Command.Command.Flags().Changed(constants.FlagClusterId) {
					return listAll(cmd)
				}

				if err := cmd.RequireExplicitLocation(); err != nil {
					return err
				}

				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)
				topics, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(
					context.Background(), clusterID,
				).Execute()
				if err != nil {
					return err
				}

				return cmd.Printer(allCols).Prefix("items").Print(topics)
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, "", "", "The ID of the cluster",
		core.RequiredFlagOption(), core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(read kafka.ClusterRead) string {
						return read.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	return cmd
}

func listAll(c *core.CommandConfig) error {
	// When --location is unset, ListAllLocations queries every location
	// concurrently, enumerating clusters and their topics per location and
	// merging the results with a Location column.
	return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
		kc := kafka.NewAPIClient(cfg)
		clusters, _, err := kc.ClustersApi.ClustersGet(context.Background()).Execute()
		if err != nil {
			return nil, err
		}

		var items []kafka.TopicRead
		for _, cluster := range clusters.Items {
			topics, _, err := kc.TopicsApi.ClustersTopicsGet(context.Background(), cluster.Id).Execute()
			if err != nil {
				return nil, err
			}
			items = append(items, topics.Items...)
		}

		return kafka.TopicReadList{Items: items}, nil
	})
}
