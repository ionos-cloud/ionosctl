package topic

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
)

func createCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "Create a Kafka topic",
			LongDesc: `Create a topic inside a Kafka cluster.

--partitions splits the topic so it can be produced/consumed in parallel; ordering is guaranteed only within a partition, and partitions cannot be reduced later.

--replication-factor is how many brokers keep a copy of each partition, giving fault tolerance. It cannot exceed the number of brokers in the cluster (3), so valid values are 1-3; 3 is recommended for production.

Log retention decides how long data is kept: --retention-time is the age (in milliseconds) after which messages become eligible for deletion (default 604800000 = 7 days); --segment-bytes is the size each on-disk log segment file reaches before a new one is rolled (default 1073741824 = 1 GiB). Retained data counts against the cluster's shared storage.

The cluster must be AVAILABLE before topics can be created.`,
			Aliases: []string{"c", "post"},
			Example: `ionosctl kafka topic create --location LOCATION --cluster-id CLUSTER_ID --name my-topic --partitions 3 --replication-factor 3
ionosctl kafka topic create --location LOCATION --cluster-id CLUSTER_ID --name events --partitions 6 --replication-factor 3 --retention-time 86400000 --segment-bytes 536870912`,
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return cmd.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.FlagName)
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
						Name:               name,
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

				return cmd.Printer(allCols).Print(topicRes)
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, "", "", "ID of the cluster to create the topic in", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(read kafka.ClusterRead) string {
						return read.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the topic", core.RequiredFlagOption())
	cmd.Command.Flags().Int32(constants.FlagKafkaPartitions, 3, "Number of partitions the topic is split into (parallelism / ordering unit); cannot be reduced later")
	cmd.Command.Flags().Int32(constants.FlagKafkaReplicationFactor, 3, "Copies of each partition kept across brokers for fault tolerance; 1-3 (cannot exceed the 3 brokers)")
	cmd.Command.Flags().Int32(constants.FlagKafkaRetentionTime, 604800000, "Age in milliseconds after which messages may be deleted (default 604800000 = 7 days)")
	cmd.Command.Flags().Int32(constants.FlagKafkaSegmentBytes, 1073741824, "Size in bytes a log segment file reaches before a new one is rolled (default 1073741824 = 1 GiB)")

	return cmd
}
