package topic

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/spf13/viper"
)

func deleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "delete",
			Namespace: "kafka",
			Resource:  "topic",
			ShortDesc: "Delete a kafka topic",
			Aliases:   []string{"d"},
			Example:   "ionosctl kafka topic delete --location LOCATION --topic-id TOPIC_ID",
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return core.CheckRequiredFlagsSets(
					cmd.Command, cmd.NS,
					[]string{constants.FlagLocation, constants.FlagClusterId, constants.FlagKafkaTopicId},
					[]string{constants.FlagLocation, constants.FlagClusterId, constants.ArgAll},
				)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				if cmd.Command.Command.Flags().Changed(constants.ArgAll) {
					err := deleteAll(cmd)
					if err != nil {
						return err
					}

					fmt.Fprintf(cmd.Command.Command.OutOrStdout(), "All topics deleted")
					return nil
				}

				clusterID := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))
				topicID := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagKafkaTopicId))

				if !confirm.FAsk(
					cmd.Command.Command.InOrStdin(), fmt.Sprintf("delete topic %v", topicID),
					viper.GetBool(constants.ArgForce),
				) {
					return fmt.Errorf(confirm.UserDenied)
				}

				_, err := client.Must().Kafka.TopicsApi.ClustersTopicsDelete(
					context.Background(), clusterID, topicID,
				).Execute()
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), "Topic %v deleted", topicID)
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
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all topics")

	return cmd
}

func deleteAll(cmd *core.CommandConfig) error {
	clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)

	topics, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(
		context.Background(), clusterID,
	).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyAndAggregateErrors(
		*topics.Items, func(topic kafka.TopicRead) error {
			if !confirm.FAsk(
				cmd.Command.Command.InOrStdin(), fmt.Sprintf("delete topic %v", *topic.Id),
				viper.GetBool(constants.ArgForce),
			) {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err := client.Must().Kafka.TopicsApi.ClustersTopicsDelete(
				context.Background(), clusterID, *topic.Id,
			).Execute()
			if err != nil {
				return fmt.Errorf("failed deleting topic %v: %w", *topic.Id, err)
			}

			return nil
		},
	)
}
