package kafka

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/topic"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "kafka",
			Short:            "The sub-commands of the 'kafka' resource help manage kafka clusters",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.Command())
	cmd.AddCommand(topic.Command())

	return core.WithRegionalFlags(cmd, constants.KafkaApiRegionalURL, constants.KafkaLocations)
}
