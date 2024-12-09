package topic

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
)

func getCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "Get a kafka topic",
			Aliases:   []string{"g"},
			Example:   "ionosctl kafka topic get --location LOCATION --cluster-id CLUSTER_ID --topic-id TOPIC_ID",
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(
					cmd.Command, cmd.NS, constants.FlagLocation, constants.FlagClusterId, constants.FlagKafkaTopicId,
				)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)
				topicID, _ := cmd.Command.Command.Flags().GetString(constants.FlagKafkaTopicId)

				topic, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsFindById(
					context.Background(), clusterID, topicID,
				).Execute()
				if err != nil {
					return err
				}

				cols, _ := cmd.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"", jsonpaths.KafkaTopic, topic, tabheaders.GetHeadersAllDefault(defaultCols, cols),
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
						return *read.Id
					},
				)
			}, constants.KafkaApiRegionalURL,
		),
	)
	cmd.AddStringFlag(
		constants.FlagKafkaTopicId, "", "", "The ID of the topic", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.Topics(cmd.Command.Flag(constants.FlagClusterId).Value.String())
			}, constants.KafkaApiRegionalURL,
		),
	)

	return cmd
}
