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

func createCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "Create a kafka topic",
			Aliases:   []string{"c", "post"},
			Example:   "ionosctl kafka topic create --location LOCATION --name my-topic --cluster-id CLUSTER_ID --partitions 1 --replication-factor 1",
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(
					cmd.Command, cmd.NS, constants.FlagClusterId, constants.FlagName, constants.FlagLocation,
				)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				name, _ := cmd.Command.Command.Flags().GetString(constants.FlagName)
				partitions, _ := cmd.Command.Command.Flags().GetInt32(constants.FlagKafkaPartitions)
				replicationFactor, _ := cmd.Command.Command.Flags().GetInt32(constants.FlagKafkaReplicationFactor)
				retentionTime, _ := cmd.Command.Command.Flags().GetInt32(constants.FlagKafkaRetentionTime)
				segmentBytes, _ := cmd.Command.Command.Flags().GetInt32(constants.FlagKafkaSegmentBytes)
				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)

				topic := kafka.NewTopicCreate(
					kafka.Topic{
						Name:               &name,
						NumberOfPartitions: &partitions,
						ReplicationFactor:  &replicationFactor,
						LogRetention: &kafka.TopicLogRetention{
							RetentionTime: &retentionTime,
							SegmentBytes:  &segmentBytes,
						},
					},
				)

				topicRes, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsPost(
					context.Background(), clusterID,
				).TopicCreate(*topic).Execute()
				if err != nil {
					return err
				}

				cols, _ := cmd.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput(
					"", jsonpaths.KafkaTopic, topicRes, tabheaders.GetHeadersAllDefault(defaultCols, cols),
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
		constants.FlagClusterId, "", "", "The ID of the cluster", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(read kafka.ClusterRead) string {
						return *read.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.AddStringFlag(constants.FlagName, "", "", "The name of the topic", core.RequiredFlagOption())
	cmd.Command.Flags().Int32(constants.FlagKafkaPartitions, 3, "The number of partitions")
	cmd.Command.Flags().Int32(constants.FlagKafkaReplicationFactor, 3, "The replication factor")
	cmd.Command.Flags().Int32(constants.FlagKafkaRetentionTime, 604800000, "The retention time in milliseconds")
	cmd.Command.Flags().Int32(constants.FlagKafkaSegmentBytes, 1073741824, "The segment bytes")

	return cmd
}
