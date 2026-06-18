package topic

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
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
				if err := core.CheckRequiredFlagsSets(
					cmd.Command, cmd.NS,
					[]string{constants.FlagLocation, constants.FlagClusterId, constants.FlagKafkaTopicId},
					[]string{constants.FlagLocation, constants.FlagClusterId, constants.ArgAll},
				); err != nil {
					return err
				}

				if err := cmd.RequireExplicitLocation(); err != nil {
					return err
				}

				return nil
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				if cmd.Command.Command.Flags().Changed(constants.ArgAll) {
					err := deleteAll(cmd)
					if err != nil {
						return err
					}

					fmt.Fprintf(cmd.Command.Command.OutOrStdout(), "All topics deleted\n")
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

				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), "Topic %v deleted\n", topicID)
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
	cmd.AddStringFlag(
		constants.FlagKafkaTopicId, "", "", "The ID of the topic", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.Topics(cmd.Command.Flag(constants.FlagClusterId).Value.String())
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all topics")

	return cmd
}

func deleteAll(cmd *core.CommandConfig) error {
	clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)

	return core.DeleteAll(cmd, core.DeleteAllOptions[kafka.TopicRead]{
		Resource: "topic",
		List: func() ([]kafka.TopicRead, error) {
			topics, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(
				context.Background(), clusterID,
			).Execute()
			if err != nil {
				return nil, err
			}
			return topics.Items, nil
		},
		Summary: func(topic kafka.TopicRead) string {
			return fmt.Sprintf("name: %s, id: %s", topic.Properties.Name, topic.Id)
		},
		ID: func(topic kafka.TopicRead) string {
			return topic.Id
		},
		Delete: func(topic kafka.TopicRead) error {
			_, err := client.Must().Kafka.TopicsApi.ClustersTopicsDelete(
				context.Background(), clusterID, topic.Id,
			).Execute()
			return err
		},
	})
}
