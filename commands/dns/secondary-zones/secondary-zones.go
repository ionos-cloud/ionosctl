package secondary_zones

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/secondary-zones/transfer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.zoneName", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "PrimaryIPs", JSONPath: "properties.primaryIps", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "secondary-zone",
			Aliases: []string{"secondary-zones", "sz"},
			Short:   "Manage secondary (transferred-in) DNS zones",
			Long: `Manage secondary DNS zones.

A secondary zone is a read-only copy of a zone whose master lives on an external primary name server. IONOS acts as a secondary: you point it at the primary's IPs (--primary-ips) and it pulls the records via zone transfer (AXFR/IXFR). You don't edit records here — you edit them on the primary and re-transfer.

Use the 'transfer' sub-commands to start a transfer and check its status.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(transfer.Root())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(updateCmd())

	return cmd
}
