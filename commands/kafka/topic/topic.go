package topic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "ReplicationFactor", JSONPath: "properties.replicationFactor", Default: true},
	{Name: "NumberOfPartitions", JSONPath: "properties.numberOfPartitions", Default: true},
	{Name: "RetentionTime", JSONPath: "properties.logRetention.retentionTime", Default: true},
	{Name: "SegmentByes", JSONPath: "properties.logRetention.segmentBytes", Default: true},
	{Name: "ClusterId", JSONPath: "href"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "topic",
			Short: "Manage Kafka topics",
			Long: `Manage Kafka topics.

A topic is a named, append-only stream of messages inside a cluster. Each topic is split into partitions (the unit of parallelism and ordering) and each partition is copied to several brokers per the replication factor. Log retention (--retention-time / --segment-bytes) controls how long messages are kept and how the on-disk log is segmented.

All topic commands target a cluster, so pass --cluster-id (and --location).`,
			Aliases:          []string{"t"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(listCmd())

	return cmd
}
