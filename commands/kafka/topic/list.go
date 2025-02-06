package topic

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
)

func listCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "List all kafka topics",
			Aliases:   []string{"ls"},
			Example: `ionosctl kafka topic list --location LOCATION
ionosctl kafka topic list --location LOCATION --cluster-id CLUSTER_ID`,
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(cmd.Command, cmd.NS, constants.FlagLocation)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				if !cmd.Command.Command.Flags().Changed(constants.FlagClusterId) {
					return listAll(cmd)
				}

				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)
				topics, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(
					context.Background(), clusterID,
				).Execute()
				if err != nil {
					return err
				}

				cols, _ := cmd.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"items", jsonpaths.KafkaTopic, topics, tabheaders.GetHeadersAllDefault(defaultCols, cols),
				)
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), out)
				return nil
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
	clusters, _, err := client.Must().Kafka.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	var allTopicsConverted []map[string]interface{}
	var allTopics []kafka.TopicReadList
	for _, cluster := range clusters.Items {
		topics, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(context.Background(), *cluster.Id).Execute()
		if err != nil {
			return err
		}

		topicsConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.KafkaTopic, topics)
		if err != nil {
			return fmt.Errorf("error converting topics to table format: %w", err)
		}

		allTopicsConverted = append(
			allTopicsConverted, functional.Map(
				topicsConverted, func(topic map[string]interface{}) map[string]interface{} {
					topic["ClusterId"] = cluster.Id
					return topic
				},
			)...,
		)
		allTopics = append(allTopics, topics)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice("cols")
	out, err := jsontabwriter.GenerateOutputPreconverted(
		allTopics, allTopicsConverted, tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
