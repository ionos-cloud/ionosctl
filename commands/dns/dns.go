package dns

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/dnssec"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/quota"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/record"
	reverse_record "github.com/ionos-cloud/ionosctl/v6/commands/dns/reverse-record"
	secondary_zones "github.com/ionos-cloud/ionosctl/v6/commands/dns/secondary-zones"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/zone"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dns",
			Short:            "The sub-commands of the 'dns' resource help automate DNS Zone and Record management",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(zone.ZoneCommand())
	cmd.AddCommand(record.RecordCommand())
	cmd.AddCommand(reverse_record.Root())
	cmd.AddCommand(quota.Root())
	cmd.AddCommand(dnssec.Root())
	cmd.AddCommand(secondary_zones.Root())

	return core.WithRegionalConfigOverride(cmd, []string{"dns"}, constants.DNSApiRegionalURL, constants.DNSLocations)
}
