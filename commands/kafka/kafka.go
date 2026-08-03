package kafka

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/topic"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/user"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "kafka",
			Short: "Manage IONOS CLOUD Kafka",
			Long: `Manage IONOS CLOUD Kafka: managed Apache Kafka clusters for streaming data.

A cluster is 3 broker nodes running inside your own private LAN; inside it you create topics (message streams, split into partitions and replicated across brokers) and users (mTLS client identities). Clients connect over TLS on port 9093.

Kafka is regional — every command targets a --location (e.g. de/txl); most 'list' commands fan out across all regions when --location is omitted.

Docs: https://docs.ionos.com/cloud/data-analytics/kafka`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.Command())
	cmd.AddCommand(topic.Command())
	cmd.AddCommand(user.Command())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Kafka}, constants.KafkaApiRegionalURL, constants.KafkaLocations)
}
