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
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "dns",
			Short: "Manage IONOS CLOUD DNS zones and records",
			Long: `Manage IONOS CLOUD DNS.

A 'zone' is the authoritative container for a domain (e.g. example.com); inside it you add 'record's (A, AAAA, CNAME, MX, TXT, …) that answer lookups. Related sub-commands:
  zone             primary zones you host and edit directly
  record           the individual DNS entries inside a zone
  reverse-record   PTR-style records mapping your IONOS IPs back to names
  secondary-zones  read-only zones transferred in from an external primary
  dnssec           sign a zone so resolvers can verify its answers
  quota            your account's DNS limits

Docs: https://docs.ionos.com/cloud/network-services/cloud-dns`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(zone.ZoneCommand())
	cmd.AddCommand(record.RecordCommand())
	cmd.AddCommand(reverse_record.Root())
	cmd.AddCommand(quota.Root())
	cmd.AddCommand(dnssec.Root())
	cmd.AddCommand(secondary_zones.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.DNS}, constants.DNSApiRegionalURL, constants.DNSLocations)
}
