package zone

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone/file"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.zoneName", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "NameServers", JSONPath: "metadata.nameServers", Default: true},
	{Name: "Enabled", JSONPath: "properties.enabled", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "zone",
			Aliases: []string{"z", "zones"},
			Short:   "Manage primary DNS zones",
			Long: `Manage primary DNS zones.

A zone is the authoritative container for one domain name (e.g. example.com): IONOS answers lookups for that domain and everything under it from the records you add here. Creating a zone does not by itself route traffic — delegate the domain to the IONOS name servers at your registrar, then add 'dns record's inside the zone.

A zone can be --enabled (served) or disabled (kept but not answered). To pull a zone in from an external primary instead of hosting it here, see 'dns secondary-zone'.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())
	cmd.AddCommand(file.Root())

	return cmd
}
