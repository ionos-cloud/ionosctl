package cluster

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId", Default: true},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId", Default: true},
	{Name: "BrokerAddresses", JSONPath: "properties.connections.0.brokerAddresses", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "StateMessage", JSONPath: "metadata.message", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Short:            "The sub-commands of 'ionosctl kafka cluster' allow you to manage kafka clusters",
			Aliases:          []string{"cl"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(FindByID())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Create())
	return cmd
}
