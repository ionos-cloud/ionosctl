package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "user",
			Short: "Manage Kafka users",
			Long: `Manage Kafka users.

A Kafka user is a client identity that authenticates to the cluster with mutual TLS (mTLS): instead of a password, each user has a client certificate + private key signed by the cluster's CA. 'user get-access' downloads those credentials so a producer/consumer can connect on port 9093.

All user commands target a cluster, so pass --cluster-id (and --location).`,
			Aliases:          []string{"u"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(GetAccess())
	return cmd
}
